package interfaces

import (
	"alarm_center/internal/application/service"
	"github.com/gin-gonic/gin"
)

func (s *Server) InitRouter() {
	r := s.router
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.GET("/check", HealthCheck)

	apiV1 := s.router.Group("/api/v1")
	s.userRoutes(apiV1)
	s.dingtalkRoutes(apiV1)
	s.emailRoutes(apiV1)
}

func (s *Server) userRoutes(api *gin.RouterGroup) {
	userRoutes := api.Group("/users")
	{
		var userSvc *service.UserService
		s.container.Invoke(func(us *service.UserService) {
			userSvc   = us
		})

		userRoutes.GET("/", userSvc.FindUsers)
		userRoutes.DELETE("/:id", userSvc.Delete)
	}
}

func (s *Server) dingtalkRoutes(api *gin.RouterGroup) {
	dingtalkRoutes := api.Group("/dt")
	{
		var dtSvc *service.DingtalkService
		s.container.Invoke(func(dts *service.DingtalkService) {
			dtSvc   = dts
		})

		dingtalkRoutes.POST("/", dtSvc.SendDingTalk)
	}
}

func (s *Server) emailRoutes(api *gin.RouterGroup) {
	dingtalkRoutes := api.Group("/email")
	{
		var emSvc *service.EmailService
		s.container.Invoke(func(em *service.EmailService) {
			emSvc   = em
		})

		dingtalkRoutes.POST("/", emSvc.SendEmail)
	}
}

// HealthCheck 监控检测
func HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, map[string]interface{}{
		"code":  0,
		"alive": true,
	})
}