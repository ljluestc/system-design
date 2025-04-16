package scrubber

type Filter struct {
    bannedPatterns []string
}

func NewFilter() *Filter {
    return &Filter{bannedPatterns: []string{"malicious"}}
}

func (f *Filter) Check(request string) bool {
    for _, pattern := range f.bannedPatterns {
        if pattern == request {
            return false
        }
    }
    return true
}