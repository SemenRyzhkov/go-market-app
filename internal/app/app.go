package app

import (
	"log"
	"net/http"
	"time"

	"github.com/SemenRyzhkov/go-market-app/internal/config"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/userhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories"
	"github.com/SemenRyzhkov/go-market-app/internal/router"
	"github.com/SemenRyzhkov/go-market-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/go-market-app/internal/service/userservice"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	urlRepository, err := repositories.New(cfg.DataBaseAddress)
	if err != nil {
		return nil, err
	}
	urlService := userservice.New(urlRepository)
	cookieService, err := cookieservice.New(cfg.Key)
	if err != nil {
		return nil, err
	}
	urlHandler := userhandlers.NewHandler(urlService, cookieService)
	urlRouter := router.NewRouter(urlHandler)

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      urlRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//defer closeHTTPServerAndStopWorkerPool(server, urlRepository)
	return &App{
		HTTPServer: server,
	}, nil
}

//func closeHTTPServerAndStopWorkerPool(server *http.Server, repository repositories.URLRepository) {
//	sigs := make(chan os.Signal, 1)
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//	go func() {
//		<-sigs
//		server.Close()
//		repository.StopWorkerPool()
//	}()
//
//}

func (app *App) Run() error {
	log.Println("run server")
	return app.HTTPServer.ListenAndServe()
}
