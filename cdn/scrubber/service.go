package scrubber

type Service struct {
    filter *Filter
}

func NewService() *Service {
    return &Service{filter: NewFilter()}
}

func (s *Service) IsSafe(request string) bool {
    return s.filter.Check(request)
}