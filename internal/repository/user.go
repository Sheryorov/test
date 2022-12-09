package repository

import (
	"time"

	"github.com/sheryorov/test/internal/entity"
	"gorm.io/gorm"
)

type userRepo struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{db}
}

func (u *userRepo) GetUserByID(userID int) (entity.User, error) {
	user := entity.User{}
	res := u.First(&user, userID)
	return user, res.Error
}

func (u *userRepo) GetAllUsers() ([]entity.User, error) {
	users := []entity.User{}
	res := u.Find(&users)
	return users, res.Error
}

func (u *userRepo) CreateUser(firstName string, lastName string, password string, birthDay time.Time) error {
	user := entity.User{
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  &birthDay,
	}
	res := u.Create(&user)
	return res.Error
}
func (u *userRepo) UpdateUser(id int, firstName string, lastName string, password string, birthDay time.Time) error {
	user := entity.User{
		ID:        uint(id),
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  &birthDay,
		Password:  password,
	}
	res := u.Save(&user)
	return res.Error
}
func (u *userRepo) DeleteUser(userID int) error {
	return u.Delete(&entity.User{}, userID).Error
}
