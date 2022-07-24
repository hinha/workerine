package tasks

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

func (c *taskHandler) PingPong() *pingpong {
	return newPingpong(c)
}

type pingpong struct {
	*taskHandler
}

func newPingpong(h *taskHandler) *pingpong {
	return &pingpong{h}
}

func (c *pingpong) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println("Pong")
	return nil
}

func (c *pingpong) GetTaskName() string {
	return "pingpong"
}
