package game

import (
	"errors"
	"math/rand"
)

// LevelNode represents a single level node in a graph.
type LevelNode struct {
	Level *Level
	// Neighbors is a map from the node to the edge weight (distance) in tiles
	Neighbors map[*LevelNode]int
}

func NewLevelNode(level *Level) *LevelNode {
	return &LevelNode{
		Level:     level,
		Neighbors: make(map[*LevelNode]int),
	}
}

// LevelGraph is an Erdos Renyi Randomized Graph
type LevelGraph map[*LevelNode][]*LevelNode

func NewLevelGraph(nNodes int, seed int) (*LevelGraph, error) {
	if nNodes <= 1 {
		return nil, errors.New("nNodes must be greater than 1")
	}

	probability := float64(1 / nNodes)

	digraph := make(LevelGraph)

	// Generate the nodes
	nodes := make([]*LevelNode, 0)
	for i := 0; i < nNodes; i++ {
		node := NewLevelNode(GrassLevel)
		nodes = append(nodes, node)
		digraph[node] = make([]*LevelNode, 0)
	}

	for i := 0; i < nNodes; i++ {
		src := nodes[i]
		for j := 0; j < nNodes; j++ {
			if j == i {
				continue
			} else if j == i+1 {
				// Nodes always connect to their direct neighbor
				digraph[src] = append(digraph[src], nodes[j])
			} else {
				// Otherwise, connect with some small probability
				if rand.Float64() < probability {
					
				}
			}
		}
	}
}
