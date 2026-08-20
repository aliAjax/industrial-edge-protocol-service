package site

import "sync"

type Node struct {
	ID     string
	Parent string
	Kind   string
	Name   string
}
type Graph struct {
	mu    sync.RWMutex
	nodes map[string]Node
	edges map[string][]string
}

func New() *Graph { return &Graph{nodes: map[string]Node{}, edges: map[string][]string{}} }
func (g *Graph) Add(node Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.nodes[node.ID]; ok {
		return &Conflict{node.ID}
	}
	if node.Parent != "" {
		if _, ok := g.nodes[node.Parent]; !ok {
			return &Missing{node.Parent}
		}
		g.edges[node.Parent] = append(g.edges[node.Parent], node.ID)
	}
	g.nodes[node.ID] = node
	return nil
}
func (g *Graph) Children(id string) []Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := []Node{}
	for _, child := range g.edges[id] {
		out = append(out, g.nodes[child])
	}
	return out
}
func (g *Graph) Path(id string) []Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := []Node{}
	for id != "" {
		n, ok := g.nodes[id]
		if !ok {
			break
		}
		out = append([]Node{n}, out...)
		id = n.Parent
	}
	return out
}

type Conflict struct{ ID string }

func (e *Conflict) Error() string { return "node exists: " + e.ID }

type Missing struct{ ID string }

func (e *Missing) Error() string { return "parent missing: " + e.ID }
