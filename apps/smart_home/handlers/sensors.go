package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"smarthome/db"
	"smarthome/models"
	"smarthome/services"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// SensorHandler handles sensor-related requests
type SensorHandler struct {
	DB                 *db.DB
	TemperatureService *services.TemperatureService
	KafkaProducer      *services.KafkaProducer
}

// NewSensorHandler creates a new SensorHandler
func NewSensorHandler(db *db.DB, temperatureService *services.TemperatureService, kp *services.KafkaProducer) *SensorHandler {
	return &SensorHandler{
		DB:                 db,
		TemperatureService: temperatureService,
		KafkaProducer:      kp,
	}
}

// RegisterRoutes registers the sensor routes
func (h *SensorHandler) RegisterRoutes(router *gin.RouterGroup) {
	sensors := router.Group("/sensors")
	{
		sensors.GET("", h.GetSensors)
		sensors.GET("/:id", h.GetSensorByID)
		sensors.POST("", h.CreateSensor)
		sensors.PUT("/:id", h.UpdateSensor)
		sensors.DELETE("/:id", h.DeleteSensor)
		sensors.PATCH("/:id/value", h.UpdateSensorValue)
		sensors.GET("/temperature/:location", h.GetTemperatureByLocation)
	}
}

// @Summary      Получить датчики
// @Description  Возвращает все датчики
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Sensor
// @Router       /api/v1/sensors [get]
func (h *SensorHandler) GetSensors(c *gin.Context) {
	sensors, err := h.DB.GetSensors(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var wg sync.WaitGroup
	errChan := make(chan error, len(sensors))
	// Update temperature sensors with real-time data from the external API
	for i, sensor := range sensors {
		if sensor.Type == models.Temperature {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
				if err == nil {
					// Update sensor with real-time data
					sensors[i].Value = tempData.Value
					sensors[i].Status = tempData.Status
					sensors[i].LastUpdated = tempData.Timestamp
					log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
				} else {
					log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
				}
				errChan <- err
			}(i)
		}
	}
	wg.Wait()
	close(errChan)
	for i := 0; i < len(sensors); i++ {
		err := <-errChan
		if err != nil {
			log.Printf("Error updating sensor data: %v", err)
		}
	}

	c.JSON(http.StatusOK, sensors)
}

// @Summary      Получить датчик по идентификатору
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "ID датчика"
// @Success      200  {object}  map[string]string  "Сообщение об успешном обновлении"
// @Failure      400  {object}  map[string]string  "Ошибка запроса или данных"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /api/v1/sensors/{id} [get]
func (h *SensorHandler) GetSensorByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	sensor, err := h.DB.GetSensorByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
		return
	}

	// If this is a temperature sensor, fetch real-time data from the temperature API
	if sensor.Type == models.Temperature {
		tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
		if err == nil {
			// Update sensor with real-time data
			sensor.Value = tempData.Value
			sensor.Status = tempData.Status
			sensor.LastUpdated = tempData.Timestamp
			log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
		} else {
			log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
		}
	}

	c.JSON(http.StatusOK, sensor)
}

// @Summary      Получить температуру в локации
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Param        location    query      string     true  "Локация"
// @Success      200  {object}  map[string]string  "Сообщение об успешном обновлении"
// @Failure      400  {object}  map[string]string  "Ошибка запроса или данных"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /api/v1/sensors/temperature [get]
func (h *SensorHandler) GetTemperatureByLocation(c *gin.Context) {
	location := c.DefaultQuery("location", "")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	// Fetch temperature data from the external API
	tempData, err := h.TemperatureService.GetTemperature(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch temperature data: %v", err),
		})
		return
	}

	// Return the temperature data
	c.JSON(http.StatusOK, gin.H{
		"location":    tempData.Location,
		"value":       tempData.Value,
		"unit":        tempData.Unit,
		"status":      tempData.Status,
		"timestamp":   tempData.Timestamp,
		"description": tempData.Description,
	})
}

// @Summary      Создать датчик
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Param        body  body 	models.SensorCreate  true "Тело"
// @Success      200  {object}  map[string]string  "Сообщение об успешном обновлении"
// @Failure      400  {object}  map[string]string  "Ошибка запроса или данных"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /api/v1/sensors [post]
func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var sensorCreate models.SensorCreate
	if err := c.ShouldBindJSON(&sensorCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensor, err := h.DB.CreateSensor(context.Background(), sensorCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	AsyncWrapper(func() error {
		h.SendTelemetry(sensor.ID, &sensor.Value, time.Now(), sensor.Status, string(sensor.Type), sensor.Unit)
		return nil
	})

	c.JSON(http.StatusCreated, sensor)
}

// @Summary      Обновить значение датчика
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "ID датчика"
// @Param        body  body 	models.SensorUpdate  true "Тело"
// @Success      200  {object}  map[string]string  "Сообщение об успешном обновлении"
// @Failure      400  {object}  map[string]string  "Ошибка запроса или данных"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /api/v1/sensors/{id} [put]
func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	var sensorUpdate models.SensorUpdate
	if err := c.ShouldBindJSON(&sensorUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensor, err := h.DB.UpdateSensor(context.Background(), id, sensorUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	AsyncWrapper(func() error {
		h.SendTelemetry(id, sensorUpdate.Value, time.Now(), sensorUpdate.Status, string(sensor.Type), sensor.Unit)
		return nil
	})

	c.JSON(http.StatusOK, sensor)
}

// @Summary      Удалить датчик
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "ID датчика"
// @Success      200  {object}  map[string]string  "Сообщение об успешном обновлении"
// @Failure      400  {object}  map[string]string  "Ошибка запроса или данных"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /api/v1/sensors/{id} [delete]
func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	err = h.DB.DeleteSensor(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted successfully"})
}

// @Summary      Обновить значение датчика
// @Tags         sensors
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "ID датчика"
// @Param        body  body 	models.SensorPath  true "Тело"
// @Success      200  {object}  map[string]string  "Сообщение об успешном обновлении"
// @Failure      400  {object}  map[string]string  "Ошибка запроса или данных"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /api/v1/sensors/{id}/value [patch]
func (h *SensorHandler) UpdateSensorValue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	var input models.SensorPath
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.DB.UpdateSensorValue(context.Background(), id, input.Value, input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	AsyncWrapper(func() error {
		h.SendTelemetry(id, input.Value, time.Now(), input.Status, "", "")
		return nil
	})

	c.JSON(http.StatusOK, gin.H{"message": "Sensor value updated successfully"})
}

func (h *SensorHandler) SendTelemetry(id int, value *float64, timestamp time.Time, status string, t string, unit string) {
	sensorData := models.SensorData{
		SensorID:  id,
		Value:     value,
		Timestamp: time.Now(),
		Status:    status,
		Type:      t,
		Unit:      unit,
	}
	jsonData, err := json.Marshal(sensorData)
	if err != nil {
		log.Fatal(err)
	}
	message := string(jsonData)
	headers := []sarama.RecordHeader{
		{Key: []byte("message-type"), Value: []byte(fmt.Sprintf("%s", "sensors_data"))},
		{Key: []byte("timestamp"), Value: []byte(timestamp.Format(time.RFC3339))},
	}

	h.KafkaProducer.SendMessage("smart-home.sensors-datas.telemetry", message, headers)
}

func AsyncWrapper(fn func() error) {
	go func() {
		// Обработка ошибки в фоновом потоке
		if err := fn(); err != nil {
			log.Printf("Error executing function asynchronously: %v", err)
		}
	}()
}
