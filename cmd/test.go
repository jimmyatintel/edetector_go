package main

import (
	"fmt"
	"github.com/google/uuid"
)

type Relation struct {
	UUID  string
	Child []string
}

func main() {
	agent := "123"
	RelationMap := make(map[string]map[int]Relation)
	RelationMap[agent] = make(map[int]Relation)
	parent := 0
	children := []int{1, 2, 3}
	_, exists := RelationMap[agent][parent]
	if !exists {
		RelationMap[agent][parent] = Relation{
			UUID:  uuid.NewString(),
			Child: []string{},
		}
		fmt.Println("uuid", RelationMap[agent][parent].UUID)
	}
	for _, child := range children {
		_, exists := RelationMap[agent][child]
		if !exists {
			RelationMap[agent][child] = Relation{
				UUID:  uuid.NewString(),
				Child: []string{},
			}
			fmt.Println("uuid", RelationMap[agent][child].UUID)
		}
		relation := RelationMap[agent][parent]
		relation.Child = append(relation.Child, RelationMap[agent][child].UUID)
		RelationMap[agent][parent] = relation

		fmt.Printf("%d: %d\n", parent, child)
		fmt.Println(relation.Child)
	}
}
