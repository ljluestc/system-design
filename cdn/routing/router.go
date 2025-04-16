package routing

import (
    "cdn/edge"
)

type Router struct {
    edges []*edge.Server
}

func NewRouter(edges []*edge.Server) *Router {
    return &Router{edges: edges}
}

func (r *Router) Route(key, location string) *edge.Server {
    return r.edges[GetEdgeID(key, len(r.edges))]
}