package origin

type Server struct {
    storage *Storage
}

func NewServer() *Server {
    return &Server{storage: NewStorage()}
}

func (s *Server) Fetch(key string) (string, bool) {
    return s.storage.Get(key)
}