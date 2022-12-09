package usecase

import "github.com/gin-gonic/gin"

type UserWebAPIHandler interface {
	GetAllUsers(*gin.Context)
	GetUserByID(*gin.Context)
	CreateUser(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}
