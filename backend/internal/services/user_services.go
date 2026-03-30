package services

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
)

type UserService struct {
	DB *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) UserService {
	return UserService{
		DB: db,
	}
}

func (s *UserService) GetUser(email string) (*models.User, error) {
	user := models.User{}

	query := `
	SELECT id, name, email, password_hash, role, provider, provider_id, avatar_url, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	err := s.DB.QueryRow(context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Provider,
		&user.ProviderID,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		log.Println("Failed to get user: ", err)
		return nil, err
	}

	return &user, nil
}
