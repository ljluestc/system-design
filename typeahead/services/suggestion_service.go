package services

import (
    "typeahead_system/core"

    "github.com/sirupsen/logrus"
)

// SuggestionService provides suggestion functionality
type SuggestionService struct {
    trie  *core.Trie
    cache *core.LRUCache
    log   *logrus.Logger
}

// NewSuggestionService creates a new SuggestionService
func NewSuggestionService(trie *core.Trie, cache *core.LRUCache, log *logrus.Logger) *SuggestionService {
    return &SuggestionService{
        trie:  trie,
        cache: cache,
        log:   log,
    }
}

// GetSuggestions retrieves suggestions for a prefix
func (s *SuggestionService) GetSuggestions(prefix string, limit int) ([]string, error) {
    if suggestions, exists := s.cache.Get(prefix); exists {
        s.log.Debugf("Cache hit for prefix: %s", prefix)
        return suggestions, nil
    }
    suggestions, err := s.trie.GetSuggestions(prefix, limit)
    if err != nil {
        s.log.Errorf("Error retrieving suggestions: %v", err)
        return nil, err
    }
    s.cache.Put(prefix, suggestions)
    s.log.Infof("Fetched and cached suggestions for prefix: %s", prefix)
    return suggestions, nil
}