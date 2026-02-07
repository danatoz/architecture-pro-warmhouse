package models

import (
	"time"
)

type HomeType string

// Структура для дома
type Home struct {
	ID        int       `json:"id"`         // Уникальный идентификатор дома
	UserID    int       `json:"user_id"`    // Идентификатор пользователя, который владеет домом (связь с пользователем)
	Address   string    `json:"address"`    // Адрес дома
	CreatedAt time.Time `json:"created_at"` // Дата создания записи
	UpdatedAt time.Time `json:"updated_at"` // Дата последнего обновления
}

// Структура для создания нового дома
type HomeCreate struct {
	UserID  int    `json:"user_id"`
	Address string `json:"address" validate:"required"`
}

// Структура для обновления информации о доме
type HomeUpdate struct {
	Address *string `json:"address,omitempty"` // Поле может быть пустым, если не обновляется
}
