package inference

import "fmt"

type Model struct {
    ID int
}

func NewModel(shardID int) *Model {
    return &Model{ID: shardID}
}

func (m *Model) Generate(prompt, context string) string {
    return fmt.Sprintf("Shard %d response to '%s' (context: %s)", m.ID, prompt, context)
}