package db

import (
    "context"
    "github.com/calelin/messenger/config"
    "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
)

// MongoDB manages database operations
type MongoDB struct {
    client *mongo.Client
    log    *logrus.Logger
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg *config.Config, log *logrus.Logger) *MongoDB {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    return &MongoDB{client: client, log: log}
}

// SaveURL stores a URL mapping
func (m *MongoDB) SaveURL(shortKey, originalURL string, expiryDate *time.Time) error {
    collection := m.client.Database("tinyurl").Collection("urls")
    _, err := collection.InsertOne(context.Background(), map[string]interface{}{
        "short_key":    shortKey,
        "original_url": originalURL,
        "expiry_date":  expiryDate,
        "created_at":   time.Now(),
    })
    if err != nil {
        m.log.Errorf("Failed to save URL %s: %v", shortKey, err)
        return err
    }
    m.log.Infof("Saved URL %s", shortKey)
    return nil
}

// GetOriginalURL retrieves the original URL
func (m *MongoDB) GetOriginalURL(shortKey string) (string, error) {
    collection := m.client.Database("tinyurl").Collection("urls")
    var result struct {
        OriginalURL string `bson:"original_url"`
    }
    err := collection.FindOne(context.Background(), map[string]string{"short_key": shortKey}).Decode(&result)
    if err != nil {
        m.log.Errorf("Failed to find URL %s: %v", shortKey, err)
        return "", err
    }
    return result.OriginalURL, nil
}

// DeleteURL removes a URL mapping
func (m *MongoDB) DeleteURL(shortKey string) error {
    collection := m.client.Database("tinyurl").Collection("urls")
    _, err := collection.DeleteOne(context.Background(), map[string]string{"short_key": shortKey})
    if err != nil {
        m.log.Errorf("Failed to delete URL %s: %v", shortKey, err)
        return err
    }
    m.log.Infof("Deleted URL %s", shortKey)
    return nil
}

// ReserveCustomAlias reserves a custom alias
func (m *MongoDB) ReserveCustomAlias(alias, originalURL string, expiryDate *time.Time) error {
    collection := m.client.Database("tinyurl").Collection("urls")
    _, err := collection.InsertOne(context.Background(), map[string]interface{}{
        "short_key":    alias,
        "original_url": originalURL,
        "expiry_date":  expiryDate,
        "created_at":   time.Now(),
    })
    if err != nil {
        m.log.Errorf("Failed to reserve alias %s: %v", alias, err)
        return err
    }
    m.log.Infof("Reserved alias %s", alias)
    return nil
}

// GetNextID generates a new unique ID
func (m *MongoDB) GetNextID() (uint64, error) {
    collection := m.client.Database("tinyurl").Collection("counters")
    result := collection.FindOneAndUpdate(
        context.Background(),
        map[string]string{"_id": "url_id"},
        map[string]interface{}{"$inc": map[string]int{"seq": 1}},
        options.FindOneAndUpdate().SetUpsert(true),
    )
    var counter struct {
        Seq uint64 `bson:"seq"`
    }
    if err := result.Decode(&counter); err != nil {
        m.log.Errorf("Failed to get next ID: %v", err)
        return 0, err
    }
    return counter.Seq, nil
}