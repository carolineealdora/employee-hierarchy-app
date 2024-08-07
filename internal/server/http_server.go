package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func RunServer(r http.Handler) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("failed on reading .env file", err.Error())
	}

	srv := http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Println("running server on port :", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("error while server listen and serve: ", err)
			}
		}
		log.Println("server is not receiving new requests...")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	graceDuration := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), graceDuration)
	defer cancel()

	log.Println("attempt to shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("error shutting down server: ", err)
	}

	log.Println("http server is shutting down gracefully")
}
