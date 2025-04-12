package db

func Replicate(key, value string, replicationFactor int) error {
    primaryNode := GetResponsibleNode(key)
    targetNodes := getReplicaNodes(primaryNode, replicationFactor)
    for _, node := range targetNodes {
        if err := node.Storage.Put(key, value); err != nil {
            return err
        }
    }
    return nil
}

func getReplicaNodes(primary *Node, count int) []*Node {
    replicas := []*Node{primary}
    for i := 1; i < count && i < len(nodes); i++ {
        replicas = append(replicas, nodes[(i+int(primary.ID[0]))%len(nodes)])
    }
    return replicas
}