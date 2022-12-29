package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/SemenRyzhkov/go-market-app/internal/app"
	"github.com/SemenRyzhkov/go-market-app/internal/common/utils"
	"github.com/SemenRyzhkov/go-market-app/internal/config"
)

func main() {
	utils.LoadEnvironments(".env")

	utils.HandleFlag()
	flag.Parse()

	serverAddress := utils.GetServerAddress()
	dbAddress := utils.GetDBAddress()
	key := utils.GetKey()

	cfg := config.New(serverAddress, key, dbAddress)
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
