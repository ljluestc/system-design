package db

import (
    "github.com/gocql/gocql"
    "messenger/internal/auth"
    "messenger/internal/chat"
)

type Cassandra struct {
    session *gocql.Session
}

func NewCassandraSession(url string) (*Cassandra, error) {
    cluster := gocql.NewCluster(url)
    cluster.Keyspace = "messenger"
    cluster.Consistency = gocql.Quorum
    session, err := cluster.CreateSession()
    if err != nil {
        return nil, err
    }
    return &Cassandra{session: session}, nil
}

func (c *Cassandra) Close() {
    c.session.Close()
}

func (c *Cassandra) SaveUser(user auth.User) error {
    return c.session.Query("INSERT INTO users (id, email, password, name) VALUES (?, ?, ?, ?)",
        user.ID, user.Email, user.Password, user.Name).Exec()
}

func (c *Cassandra) GetUser(userID string) (auth.User, error) {
    var user auth.User
    err := c.session.Query("SELECT id, email, password, name FROM users WHERE id = ?", userID).
        Scan(&user.ID, &user.Email, &user.Password, &user.Name)
    return user, err
}

func (c *Cassandra) SaveMessage(msg chat.Message) error {
    return c.session.Query("INSERT INTO messages (id, sender_id, recipient_id, group_id, content, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
        msg.ID, msg.SenderID, msg.RecipientID, msg.GroupID, msg.Content, msg.Timestamp).Exec()
}