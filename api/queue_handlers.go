package api

import (
	"encoding/json"
	"net/http"

	"github.com/hibiken/asynq"

	workerine "github.com/hinha/workerine/server"
)

func newListQueuesHandlerFunc(inspector *asynq.Inspector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qnames, err := inspector.Queues()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		snapshots := make([]*workerine.QueueStateSnapshot, len(qnames))
		for i, qname := range qnames {
			qinfo, err := inspector.GetQueueInfo(qname)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			snapshots[i] = workerine.SetQueueStateSnapshot(qinfo)
		}
		payload := map[string]interface{}{"queues": snapshots}
		json.NewEncoder(w).Encode(payload)
	}
}
