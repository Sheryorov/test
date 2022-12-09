package repository

import (
	"time"

	"github.com/sheryorov/test/internal/entity"
)

type UserRepository interface {
	GetUserByID(int) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
	CreateUser(string, string, string, time.Time) error
	UpdateUser(int, string, string, string, time.Time) error
	DeleteUser(int) error
}
