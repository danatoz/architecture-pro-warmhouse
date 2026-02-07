package handlers

import (
	"auth/db"
	"auth/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	DB *db.DB
}

// NewHomeHandler создаёт новый HomeHandler
func NewHomeHandler(db *db.DB) *HomeHandler {
	return &HomeHandler{
		DB: db,
	}
}

// RegisterRoutes регистрирует маршруты для домов
func (h *HomeHandler) RegisterRoutes(router *gin.RouterGroup) {
	homes := router.Group("/homes")
	{
		homes.GET("", h.GetAllHomes)
		homes.GET("/:id", h.GetHomeByID)
		homes.POST("", h.CreateHome)
		homes.PUT("/:id", h.UpdateHome)
		homes.DELETE("/:id", h.DeleteHome)
	}
}

// @Summary      Получить дома пользователя
// @Description  Возвращает список домов по user_id
// @Tags         homes
// @Accept       json
// @Produce      json
// @Param        user_id  query     int  true  "User ID"
// @Success      200      {array}   models.Home
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/homes [get]
func (h *HomeHandler) GetAllHomes(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	homes, err := h.DB.GetHomesByUserID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, homes)
}

// @Summary      Получить дом
// @Description  Возвращает дом по ID
// @Tags         homes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Home ID"
// @Success      200  {object}  models.Home
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/homes/{id} [get]
func (h *HomeHandler) GetHomeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid home id"})
		return
	}

	home, err := h.DB.GetHomeByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, home)
}

// @Summary      Создать дом
// @Description  Создаёт новый дом
// @Tags         homes
// @Accept       json
// @Produce      json
// @Param        home  body      models.HomeCreate  true  "Home data"
// @Success      201   {object}  models.Home
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/homes [post]
func (h *HomeHandler) CreateHome(c *gin.Context) {
	var input models.HomeCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	home, err := h.DB.CreateHome(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, home)
}

// @Summary      Обновить дом
// @Description  Обновляет информацию о доме
// @Tags         homes
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "Home ID"
// @Param        home  body      models.Home  true  "Updated home data"
// @Success      200   {object}  models.Home
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/homes/{id} [put]
func (h *HomeHandler) UpdateHome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid home id"})
		return
	}

	var input models.Home
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedHome, err := h.DB.UpdateHome(context.Background(), id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedHome)
}

// @Summary      Удалить дом
// @Description  Удаляет дом по ID
// @Tags         homes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Home ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/homes/{id} [delete]
func (h *HomeHandler) DeleteHome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid home id"})
		return
	}

	if err := h.DB.DeleteHome(context.Background(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "home deleted"})
}
