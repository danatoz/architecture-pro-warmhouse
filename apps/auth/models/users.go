package models

import (
	"time"
)

// Тип пользователя (например, администратор, обычный пользователь и т.д.)
type UserType string

// Структура для пользователя
type User struct {
	ID        int       `json:"id"`         // Уникальный идентификатор пользователя
	FirstName string    `json:"first_name"` // Имя пользователя
	LastName  string    `json:"last_name"`  // Фамилия пользователя
	Email     string    `json:"email"`      // Email пользователя
	Password  string    `json:"password"`   // Хэш пароля
	Role      UserType  `json:"role"`       // Роль пользователя (например, "admin" или "user")
	CreatedAt time.Time `json:"created_at"` // Дата создания записи
	UpdatedAt time.Time `json:"updated_at"` // Дата последнего обновления
}

// Структура для создания нового пользователя
type UserCreate struct {
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=6"`
	Role      UserType `json:"role" validate:"required,oneof=admin user"` // Обратите внимание на валидацию роли
}

// Структура для обновления информации о пользователе
type UserUpdate struct {
	FirstName *string   `json:"first_name,omitempty"` // Поле может быть пустым, если не обновляется
	LastName  *string   `json:"last_name,omitempty"`  // Поле может быть пустым, если не обновляется
	Email     *string   `json:"email,omitempty"`      // Поле может быть пустым, если не обновляется
	Password  *string   `json:"password,omitempty"`   // Поле может быть пустым, если не обновляется
	Role      *UserType `json:"role,omitempty"`       // Поле может быть пустым, если не обновляется
}
