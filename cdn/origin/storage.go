package origin

type Storage struct {
    data map[string]string
}

func NewStorage() *Storage {
    s := &Storage{data: make(map[string]string)}
    s.data["img1"] = "image_data"
    return s
}

func (s *Storage) Get(key string) (string, bool) {
    v, ok := s.data[key]
    return v, ok
}