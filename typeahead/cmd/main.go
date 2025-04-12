package main

import (
    "typeahead_system/api"
    "typeahead_system/core"
    "typeahead_system/services"
    "typeahead_system/utils"

    "github.com/sirupsen/logrus"
)

func main() {
    log := utils.NewLogger()
    // Initialize components
    trie := core.NewTrie()
    cache := core.NewLRUCache(100) // Cache with capacity of 100
    suggestionService := services.NewSuggestionService(trie, cache, log)

    // Set up API
    suggestionAPI := api.NewSuggestionAPI(suggestionService, log)

    // Start server on port 8080
    api.StartServer("8080", suggestionAPI.GetSuggestionsHandler, nil)
}