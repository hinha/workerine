package server

import (
	"github.com/hinha/workerine/server/server/routers"
	"github.com/hinha/workerine/server/tasks"
)

type Server interface {
	Router() []*routers.Router
}

type server struct {
	handler tasks.TaskHandler
}

func NewServer(handler tasks.TaskHandler) Server {
	return &server{handler: handler}
}

func (s *server) Router() []*routers.Router {
	return []*routers.Router{
		{
			Path:    s.handler.PingPong().GetTaskName(),
			Handler: s.handler.PingPong(),
		},
	}
}
