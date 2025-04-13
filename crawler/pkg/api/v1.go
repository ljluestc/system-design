package api

import (
    "github.com/gorilla/mux"
    "crawler/internal/scheduler"
    "crawler/pkg/models"
)

type Router struct {
    scheduler *scheduler.Scheduler
}

func NewRouter(s *scheduler.Scheduler) *mux.Router {
    r := &Router{scheduler: s}
    router := mux.NewRouter()
    router.HandleFunc("/urls", r.AddURL).Methods("POST")
    return router
}

func (r *Router) AddURL(w http.ResponseWriter, req *http.Request) {
    var url models.URL
    // Handle JSON decoding and add to scheduler
}