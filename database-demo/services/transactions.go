package services

func ExecuteTransaction(shardID int, fn func(*sql.Tx) error) error {
    db := dbPools[shardID]
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}