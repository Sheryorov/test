package usecase

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sheryorov/test/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userWebAPIHandler struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

type reqBody struct {
	FirstName string
	LastName  string
	BirthDay  time.Time
	Password  string
}

func NewUserIteractor(db *gorm.DB, l *zap.Logger) *userWebAPIHandler {
	r := repository.NewUserRepository(db)
	return &userWebAPIHandler{
		userRepo: r,
		logger:   l,
	}
}

// GetAllUsers godoc
// @Summary      Get users LIST
// @Description  GET all users lists
// @Tags         users
// @Produce      json
// @Success      200  {object}  entity.User
// @Router       /users [get]
func (u *userWebAPIHandler) GetAllUsers(c *gin.Context) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": users})
}

func (u *userWebAPIHandler) GetUserByID(c *gin.Context) {
	userID := c.Params.ByName("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
	user, err := u.userRepo.GetUserByID(userIDInt)
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (u *userWebAPIHandler) CreateUser(c *gin.Context) {
	userReq := reqBody{}
	if err := c.ShouldBindJSON(&userReq); err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if err := u.userRepo.CreateUser(userReq.FirstName, userReq.LastName, userReq.Password, userReq.BirthDay); err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (u *userWebAPIHandler) UpdateUser(c *gin.Context) {
	userID := c.Params.ByName("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
	userReq := reqBody{}
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if err := u.userRepo.UpdateUser(userIDInt, userReq.FirstName, userReq.LastName, userReq.Password, userReq.BirthDay); err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (u *userWebAPIHandler) DeleteUser(c *gin.Context) {
	userID := c.Params.ByName("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}

	if err := u.userRepo.DeleteUser(userIDInt); err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
