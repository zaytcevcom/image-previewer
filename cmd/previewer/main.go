package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zaytcevcom/image-previewer/internal/app"
	"github.com/zaytcevcom/image-previewer/internal/cacher"
	"github.com/zaytcevcom/image-previewer/internal/fetcher"
	"github.com/zaytcevcom/image-previewer/internal/logger"
	"github.com/zaytcevcom/image-previewer/internal/resizer"
	internalhttp "github.com/zaytcevcom/image-previewer/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/previewer/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config: ", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(config.Logger.Level, nil)

	previewer := app.New(
		logg,
		fetcher.New(),
		cacher.New(config.Cache.Capacity),
		resizer.New(),
	)

	port := 80
	server := internalhttp.New(logg, previewer, "", port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info(fmt.Sprintf("Previewer listening on port: %d", port))

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
