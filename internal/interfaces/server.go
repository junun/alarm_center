package interfaces

import (
	"alarm_center/internal/application/service"
	"alarm_center/internal/config"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server server
type Server struct {
	config      *config.AppConfig
	router      *gin.Engine
	container   *dig.Container
}

func NewServer(router *gin.Engine, conf *config.AppConfig, di *dig.Container) *Server {
	return &Server{
		router: 	router,
		config:     conf,
		container: 	di,
	}
}

func (s *Server) CronJob() {
	var cronSvc *service.CronService
	s.container.Invoke(func(cs *service.CronService) {
		cronSvc   = cs
	})

	go func(cronSvc *service.CronService) {
		cronSvc.StartJobOnBoot()
	}(cronSvc)

	cronSvc.JobRun()
}

// Handler mux http handler
//func (s *Server) Handler() http.Handler {
//	mux := http.NewServeMux()
//	mux.HandleFunc("/", s.index)
//	mux.HandleFunc("/users", s.Users)
//
//	return mux
//}

func (s *Server) Run() {
	address := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	log.Printf("server run on: %s\n", address)

	// create http services
	server := &http.Server{
		// Handler: http.TimeoutHandler(router, time.Second*6, `{code:503,"message":"services timeout"}`),
		Handler:      s.router,
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//go func() {
	//	log.Println("StartJobOnBoot will run...")
	//	var cronSvc *service.CronService
	//
	//	err := s.container.Invoke(func(cs *service.CronService) {
	//		cronSvc   = cs
	//	})
	//	fmt.Println(err)
	//	service.StartJob(cronSvc)
	//}()

	// run http services in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Println("services listen error:", err)
				return
			}

			log.Println("services will exit...")
		}
	}()

	// graceful exit
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recv signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), s.config.GraceWait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx)
	<-ctx.Done()

	log.Println("services shutdown success")
}

// Users user http handler
//func (s *Server) index(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("hello inject demo"))
//}
//
//func (s *Server) Users(w http.ResponseWriter, r *http.Request) {
//	users, err := s.userService.FindUsers()
//	if err != nil {
//		w.Write([]byte("request error: " + err.Error()))
//		return
//	}
//
//	var b []byte
//	b, err = json.Marshal(users)
//	if err != nil {
//		w.Write([]byte("json encode error: " + err.Error()))
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	w.Write(b)
//}
