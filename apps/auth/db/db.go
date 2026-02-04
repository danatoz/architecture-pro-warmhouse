package db

import (
	"auth/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new DB instance
func New(connString string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// CreateUser добавляет нового пользователя в базу данных
func (db *DB) CreateUser(ctx context.Context, user *models.UserCreate) (*models.User, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, first_name, last_name, email, role, created_at, updated_at
	`

	var newUser models.User
	err := db.Pool.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Role, time.Now(), time.Now()).Scan(
		&newUser.ID,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Email,
		&newUser.Role,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

// UpdateUser обновляет информацию о пользователе
func (db *DB) UpdateUser(ctx context.Context, userID int, user *models.UserUpdate) (*models.User, error) {
	query := `
		UPDATE users
		SET first_name = COALESCE($1, first_name), 
			last_name = COALESCE($2, last_name),
			email = COALESCE($3, email),
			password = COALESCE($4, password),
			role = COALESCE($5, role),
			updated_at = $6
		WHERE id = $7
		RETURNING id, first_name, last_name, email, role, created_at, updated_at
	`

	var updatedUser models.User
	err := db.Pool.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Role, time.Now(), userID).Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.Role,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

// GetUserByID получает пользователя по ID
func (db *DB) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// DeleteUser удаляет пользователя по ID
func (db *DB) DeleteUser(ctx context.Context, userID int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	commandTag, err := db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}

// CreateHome добавляет новый дом
func (db *DB) CreateHome(ctx context.Context, h models.HomeCreate) (models.Home, error) {
	query := `
		INSERT INTO homes (user_id, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	var home models.Home
	now := time.Now()
	err := db.Pool.QueryRow(ctx, query, h.UserID, h.Address, now, now).Scan(
		&home.ID,
		&home.UserID,
		&home.Address,
		&home.CreatedAt,
		&home.UpdatedAt,
	)

	if err != nil {
		return models.Home{}, fmt.Errorf("failed to create home: %w", err)
	}

	return home, nil
}

// GetHomeByID получает дом по ID
func (db *DB) GetHomeByID(ctx context.Context, homeID int) (*models.Home, error) {
	query := `
		SELECT id, user_id, address, created_at, updated_at
		FROM homes
		WHERE id = $1
	`

	var home models.Home
	err := db.Pool.QueryRow(ctx, query, homeID).Scan(
		&home.ID,
		&home.UserID,
		&home.Address,
		&home.CreatedAt,
		&home.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("home with id %d not found", homeID)
		}
		return nil, fmt.Errorf("failed to get home: %w", err)
	}

	return &home, nil
}

// GetHomesByUserID получает все дома конкретного пользователя
func (db *DB) GetHomesByUserID(ctx context.Context, userID int) ([]*models.Home, error) {
	query := `
		SELECT id, user_id, address, created_at, updated_at
		FROM homes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get homes: %w", err)
	}
	defer rows.Close()

	var homes []*models.Home
	for rows.Next() {
		var home models.Home
		if err := rows.Scan(&home.ID, &home.UserID, &home.Address, &home.CreatedAt, &home.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan home: %w", err)
		}
		homes = append(homes, &home)
	}

	return homes, nil
}

// UpdateHome обновляет данные дома
func (db *DB) UpdateHome(ctx context.Context, homeID int, home *models.Home) (*models.Home, error) {
	query := `
		UPDATE homes
		SET address = COALESCE($1, address),
			updated_at = $2
		WHERE id = $3
		RETURNING updated_at
	`

	now := time.Now()
	err := db.Pool.QueryRow(ctx, query, home.Address, now, homeID).Scan(&home.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("home with id %d not found", homeID)
		}
		return nil, fmt.Errorf("failed to update home: %w", err)
	}

	home.ID = homeID
	return home, nil
}

// DeleteHome удаляет дом по ID
func (db *DB) DeleteHome(ctx context.Context, homeID int) error {
	query := `
		DELETE FROM homes
		WHERE id = $1
	`

	commandTag, err := db.Pool.Exec(ctx, query, homeID)
	if err != nil {
		return fmt.Errorf("failed to delete home: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("home with id %d not found", homeID)
	}

	return nil
}
