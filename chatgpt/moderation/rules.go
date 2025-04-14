package moderation

type Rules struct {
    filter *Filter
}

func NewRules() *Rules {
    return &Rules{filter: NewFilter()}
}

func (r *Rules) Apply(text string) bool {
    return r.filter.IsSafe(text)
}