package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AleksandrVishniakov/tgbots-observer/internal/app"
	"github.com/AleksandrVishniakov/tgbots-observer/internal/configs"
)

const name = "tgbots-observer"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println(name, "start")

	cfg := configs.MustConfigs()
	app := app.New(cfg, os.Stdout)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("%s app closed: %s\n", name, err.Error())
	}
}
