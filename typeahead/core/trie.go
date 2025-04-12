package core

// Trie represents a basic Trie data structure
type Trie struct {
    root *TrieNode
}

// TrieNode represents a node in the Trie
type TrieNode struct {
    children    map[rune]*TrieNode
    isEndOfWord bool
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
    return &Trie{
        root: &TrieNode{children: make(map[rune]*TrieNode)},
    }
}

// Insert adds a word to the Trie
func (t *Trie) Insert(word string) {
    node := t.root
    for _, char := range word {
        if _, exists := node.children[char]; !exists {
            node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
        }
        node = node.children[char]
    }
    node.isEndOfWord = true
}

// GetSuggestions retrieves suggestions for a prefix
func (t *Trie) GetSuggestions(prefix string, limit int) ([]string, error) {
    node := t.root
    for _, char := range prefix {
        if _, exists := node.children[char]; !exists {
            return []string{}, nil // Prefix not found, return empty list
        }
        node = node.children[char]
    }
    suggestions := []string{}
    t.collectSuggestions(node, prefix, &suggestions)
    if len(suggestions) > limit {
        return suggestions[:limit], nil
    }
    return suggestions, nil
}

// collectSuggestions gathers all words under a node
func (t *Trie) collectSuggestions(node *TrieNode, prefix string, suggestions *[]string) {
    if node.isEndOfWord {
        *suggestions = append(*suggestions, prefix)
    }
    for char, child := range node.children {
        t.collectSuggestions(child, prefix+string(char), suggestions)
    }
}