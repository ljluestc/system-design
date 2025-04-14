package context

type Cache struct {
    store *Store
}

func NewCache() *Cache {
    return &Cache{store: NewStore()}
}

func (c *Cache) Set(userID, value string) {
    c.store.SaveHistory(userID, value)
}

func (c *Cache) Get(userID string) (string, bool) {
    return c.store.GetHistory(userID)
}