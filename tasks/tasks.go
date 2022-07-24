package tasks

type taskHandler struct {
}

type TaskHandler interface {
	PingPong() *pingpong
}

func NewTask() TaskHandler {
	return &taskHandler{}
}
