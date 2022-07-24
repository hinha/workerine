package workerine

import (
	"time"

	"github.com/hibiken/asynq"
)

type QueueStateSnapshot struct {
	// Name of the queue.
	Queue string `json:"queue"`
	// Total number of bytes the queue and its tasks require to be stored in redis.
	MemoryUsage int64 `json:"memory_usage_bytes"`
	// Total number of tasks in the queue.
	Size int `json:"size"`
	// Totoal number of groups in the queue.
	Groups int `json:"groups"`
	// Latency of the queue in milliseconds.
	LatencyMillisec int64 `json:"latency_msec"`
	// Latency duration string for display purpose.
	DisplayLatency string `json:"display_latency"`

	// Number of tasks in each state.
	Active      int `json:"active"`
	Pending     int `json:"pending"`
	Aggregating int `json:"aggregating"`
	Scheduled   int `json:"scheduled"`
	Retry       int `json:"retry"`
	Archived    int `json:"archived"`
	Completed   int `json:"completed"`

	// Total number of tasks processed during the given date.
	// The number includes both succeeded and failed tasks.
	Processed int `json:"processed"`
	// Breakdown of processed tasks.
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	// Paused indicates whether the queue is paused.
	// If true, tasks in the queue will not be processed.
	Paused bool `json:"paused"`
	// Time when this snapshot was taken.
	Timestamp time.Time `json:"timestamp"`
}

func SetQueueStateSnapshot(info *asynq.QueueInfo) *QueueStateSnapshot {
	return &QueueStateSnapshot{
		Queue:           info.Queue,
		MemoryUsage:     info.MemoryUsage,
		Size:            info.Size,
		Groups:          info.Groups,
		LatencyMillisec: info.Latency.Milliseconds(),
		DisplayLatency:  info.Latency.Round(10 * time.Millisecond).String(),
		Active:          info.Active,
		Pending:         info.Pending,
		Aggregating:     info.Aggregating,
		Scheduled:       info.Scheduled,
		Retry:           info.Retry,
		Archived:        info.Archived,
		Completed:       info.Completed,
		Processed:       info.Processed,
		Succeeded:       info.Processed - info.Failed,
		Failed:          info.Failed,
		Paused:          info.Paused,
		Timestamp:       info.Timestamp,
	}
}
