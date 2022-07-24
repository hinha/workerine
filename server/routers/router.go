package routers

import (
	"log"

	"github.com/hibiken/asynq"

	"github.com/hinha/workerine/server/config"
)

type Router struct {
	Path    string
	Handler asynq.Handler
}

type routing struct {
	routers []*Router
}

func New(routers []*Router) Routers {
	return &routing{routers: routers}
}

// Routers contains the functions of http handler to clean payloads and pass it the service
type Routers interface {
	Serve(cfg *config.Config)
}

func (r *routing) Serve(cfg *config.Config) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Addr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": int(cfg.Task.Queues.Critical),
				"default":  int(cfg.Task.Queues.Default),
				"low":      int(cfg.Task.Queues.Low),
			},
		},
	)

	mux := asynq.NewServeMux()
	for _, router := range r.routers {
		mux.Handle(router.Path, router.Handler)
	}

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
