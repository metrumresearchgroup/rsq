package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/metrumresearchgroup/rsq/server"
	"github.com/sirupsen/logrus"
)

// NewHTTPServer provides a new http server
func NewHTTPServer(js server.JobService, version string, port string, n int, lg *logrus.Logger) {
	httpClient := NewJobHandler(js, n, lg)
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(NewStructuredLogger(lg))
	r.Use(middleware.Recoverer)

	// When a client closes their connection midway through a request, the
	// http.CloseNotifier will cancel the request context (ctx).
	// r.Use(middleware.CloseNotify)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	// r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(version))
	})

	r.Route("/job", func(r chi.Router) {
		r.Post("/", httpClient.HandleSubmitJob)
		r.Route("/{jobID}", func(r chi.Router) {
			r.Use(httpClient.JobCtx)
			r.Get("/", httpClient.HandleGetJobByID) // GET /jobs/123
		})
		r.Route("/cancel", func(r chi.Router) {
			r.Route("/{jobID}", func(r chi.Router) {
				r.Use(httpClient.JobCtx)
				r.Put("/", httpClient.HandleCancelJob) // GET /jobs/123
			})
		})
	})
	r.Route("/jobs", func(r chi.Router) {
		r.Get("/", httpClient.HandleGetJobsByStatus)
	})

	lg.Info(fmt.Sprintf("rsq listening on %v", port))
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
