package repository

import (
	"database/sql"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user domain.User) error {
	_, err := r.db.Exec(`
        INSERT INTO users (full_name, email, password_hash, role)
        VALUES ($1, $2, $3, $4)
    `, user.FullName, user.Email, user.PasswordHash, user.Role)

	return err
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.QueryRow(`
        SELECT id, full_name, email, password_hash, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `, email).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(id string) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.QueryRow(`
        SELECT id, full_name, email, password_hash, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `, id).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
