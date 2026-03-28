package services

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/auth"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
)

type AuthService struct {
	DB *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) AuthService {
	return AuthService{
		DB: db,
	}
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}

	query := `
	SELECT id, name, email, role, provider, provider_id, avatar_url, created_at, updated_at
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

func (s *AuthService) Register(userInput dto.UserRegister) (*models.User, error) {

	user := models.User{}

	hashedPassword, err := auth.HashPassword(userInput.Password)

	if err != nil {
		return nil, err
	}

	existingUser, err := s.GetUserByEmail(userInput.Email)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	query := `
	INSERT INTO users (name, email, password_hash, provider)
	VALUES ($1, $2, $3, 'local')
	RETURNING id, name, email, role, provider, provider_id, avatar_url, created_at, updated_at
	`
	err = s.DB.QueryRow(context.Background(),
		query,
		userInput.Name,
		userInput.Email,
		hashedPassword,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Provider,
		&user.ProviderID,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("Failed to add user: ", err)
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(userInput dto.UserLogin) (string, error) {
	existingUser, err := s.GetUserByEmail(userInput.Email)

	if err != nil {
		log.Println("Failed to get user: ", err)
		return "", errors.New("user not found")
	}

	if existingUser.PasswordHash == nil {
		return "", errors.New("invalid login method")
	}

	isPasswordSame := auth.CheckPassword(userInput.Password, *existingUser.PasswordHash)

	if !isPasswordSame {
		log.Println("Incorrect password")
		return "", errors.New("incorrect password or email")
	}

	token, err := auth.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Role)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *AuthService) GetOrCreateOAuthUser(email, name, avatarURL, provider, providerID string) (*models.User, error) {

	existingUser := models.User{}

	findUserByProviderAndProviderIDQuery := `
	SELECT id, name, email, role, provider, provider_id, avatar_url, created_at, updated_at
	FROM users 
	WHERE provider = $1 AND provider_id = $2
	`
	linkExistingUserQuery := `
	UPDATE users
	SET name = $1, avatar_url = $2, provider = $3, provider_id = $4
	WHERE email = $5
	RETURNING id, name, email, role, provider, provider_id, avatar_url, created_at, updated_at
	`

	createNewUserQuery := `
	INSERT INTO users (name, email, provider, provider_id, avatar_url)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, name, email, role, provider, provider_id, avatar_url, created_at, updated_at
	`

	err := s.DB.QueryRow(context.Background(), findUserByProviderAndProviderIDQuery, provider, providerID).Scan(
		&existingUser.ID,
		&existingUser.Name,
		&existingUser.Email,
		&existingUser.Role,
		&existingUser.Provider,
		&existingUser.ProviderID,
		&existingUser.AvatarUrl,
		&existingUser.CreatedAt,
		&existingUser.UpdatedAt,
	)

	if err == nil {
		return &existingUser, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	foundUserByEmail, err := s.GetUserByEmail(email)

	if err != nil {
		log.Println("User not found using email")
	}

	if foundUserByEmail != nil {
		err = s.DB.QueryRow(context.Background(), linkExistingUserQuery, name, avatarURL, provider, providerID, foundUserByEmail.Email).Scan(
			&existingUser.ID,
			&existingUser.Name,
			&existingUser.Email,
			&existingUser.Role,
			&existingUser.Provider,
			&existingUser.ProviderID,
			&existingUser.AvatarUrl,
			&existingUser.CreatedAt,
			&existingUser.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to link users details")
			return nil, err
		}
		return &existingUser, nil
	}

	err = s.DB.QueryRow(context.Background(), createNewUserQuery, name, email, provider, providerID, avatarURL).Scan(
		&existingUser.ID,
		&existingUser.Name,
		&existingUser.Email,
		&existingUser.Role,
		&existingUser.Provider,
		&existingUser.ProviderID,
		&existingUser.AvatarUrl,
		&existingUser.CreatedAt,
		&existingUser.UpdatedAt,
	)

	if err != nil {
		log.Println("Failed to create user")
		return nil, err
	}

	return &existingUser, nil
}
