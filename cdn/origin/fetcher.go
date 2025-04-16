package origin

type Fetcher struct {
    server *Server
}

func NewFetcher() *Fetcher {
    return &Fetcher{server: NewServer()}
}

func (f *Fetcher) Fetch(key string) (string, bool) {
    return f.server.Fetch(key)
}