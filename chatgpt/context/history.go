package context

type HistoryManager struct {
    store *Store
}

func NewHistoryManager() *HistoryManager {
    return &HistoryManager{store: NewStore()}
}

func (hm *HistoryManager) Update(userID, prompt, response string) {
    history, _ := hm.store.GetHistory(userID)
    newHistory := history + "\n" + prompt + " -> " + response
    hm.store.SaveHistory(userID, newHistory)
}

func (hm *HistoryManager) Get(userID string) string {
    h, _ := hm.store.GetHistory(userID)
    return h
}