package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/hinha/workerine/server/config"
)

func Server(cfg *config.Config) (*http.Server, *HTTPHandler) {
	redisConnOpt, err := redisConnOpt(cfg)
	if err != nil {
		log.Fatal(err)
	}

	h := New(Options{
		RedisConnOpt: redisConnOpt,
	})

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})
	mux.Handle("/", c.Handler(h))

	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return srv, h
}
