package moderation

import "strings"

type Filter struct {
    bannedWords []string
}

func NewFilter() *Filter {
    return &Filter{bannedWords: []string{"harmful", "offensive"}}
}

func (f *Filter) IsSafe(text string) bool {
    for _, word := range f.bannedWords {
        if strings.Contains(strings.ToLower(text), word) {
            return false
        }
    }
    return true
}