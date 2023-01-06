package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SemenRyzhkov/go-market-app/internal/config"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/orderhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/userhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/withdrawhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/orderrepository"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/userrepository"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/withdrawrepository"
	"github.com/SemenRyzhkov/go-market-app/internal/router"
	"github.com/SemenRyzhkov/go-market-app/internal/security"
	"github.com/SemenRyzhkov/go-market-app/internal/service/orderservice"
	"github.com/SemenRyzhkov/go-market-app/internal/service/userservice"
	"github.com/SemenRyzhkov/go-market-app/internal/service/withdrawservice"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	db, err := repositories.InitDB(cfg.DataBaseAddress)
	if err != nil {
		return nil, err
	}
	userRepository := userrepository.New(db)
	orderRepository, err := orderrepository.New(db, cfg.AccrualServiceAddress, cfg.ClientDuration)
	if err != nil {
		return nil, err
	}
	withdrawRepository := withdrawrepository.New(db)
	userService := userservice.New(userRepository)
	orderService := orderservice.New(orderRepository)
	withdrawService := withdrawservice.New(withdrawRepository, orderRepository)
	//cookieService, err := cookieservice.New(cfg.Key)
	//if err != nil {
	//	return nil, err
	//}
	jwtHelper, err := security.New(cfg.Key)
	if err != nil {
		return nil, err
	}
	urlHandler := userhandlers.NewHandler(userService, jwtHelper)
	orderHandler := orderhandlers.NewHandler(orderService, jwtHelper)
	withdrawHandler := withdrawhandlers.NewHandler(withdrawService, jwtHelper)
	urlRouter := router.NewRouter(urlHandler, orderHandler, withdrawHandler)

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      urlRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	defer closeHTTPServerAndStopWorkerPoolAndScheduler(server, orderRepository)
	return &App{
		HTTPServer: server,
	}, nil
}

func closeHTTPServerAndStopWorkerPoolAndScheduler(server *http.Server, repository orderrepository.OrderRepository) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		server.Close()
		repository.StopSchedulerAndWorkerPool()
	}()

}

func (app *App) Run() error {
	log.Println("run server")
	return app.HTTPServer.ListenAndServe()
}
