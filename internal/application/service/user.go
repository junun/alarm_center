package service

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

// UserService user service
type UserService struct {
	config  *config.AppConfig
	repo 	repo.UserRepository
	rdb 	*redis.Client
}

// NewUserService return user service
func NewUserService(config *config.AppConfig, userRepo repo.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		config:     config,
		repo: 		userRepo,
		rdb: rdb,
	}
}

// FindUsers return users
func (s *UserService) FindUsers(c *gin.Context)  {

	users, err := s.repo.FindAll()
	if err != nil {
		log.Println("get users error: ", err)
	}

	fmt.Println(users)

	fmt.Println(s.rdb.Ping().Result())

	c.JSON(200, map[string]interface{}{
		"code":  0,
		"alive": true,
	})
}

func (s *UserService)  Delete(c *gin.Context)  {
	s.repo.Delete(c.PostForm("id"))
}

