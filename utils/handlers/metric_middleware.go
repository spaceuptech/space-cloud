package handlers

import (
	"github.com/gorilla/mux"
	"github.com/miolini/datacounter"
	"github.com/spaceuptech/space-cloud/utils/metrics"
	"net/http"
)

func HandleMetricMiddleWare(next http.Handler, metrics *metrics.Module) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID, ok := vars["project"]
		if ok {
			readerCounter := datacounter.NewReaderCounter(r.Body)
			writerCounter := datacounter.NewResponseWriterCounter(w)
			next.ServeHTTP(writerCounter, r)
			metrics.AddIngress(projectID, readerCounter.Count()+200)
			metrics.AddEgress(projectID, writerCounter.Count()+200)
			return
		}
		next.ServeHTTP(w, r)

	})
}
