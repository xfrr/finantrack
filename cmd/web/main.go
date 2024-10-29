package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/xfrr/finantrack/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app.RouteWithRegexp(web.MainPath, func() app.Composer {
		return &web.App{}
	})

	app.RunWhenOnBrowser()

	// Server routing:
	handler := http.NewServeMux()
	handler.Handle(web.MainPath, web.NewHandler())

	server := &http.Server{
		Addr:         ":8000",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Log("server error", err)
		}
	}()

	<-ctx.Done()
}
