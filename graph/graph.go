package graph

import (
	"fmt"
	"sync"
)

type Node struct{
	Name string
}

func (n *Node) String() string{
	return fmt.Sprintf("%v", n.Name)
}

type Graph struct {
	EntryNode Node
	nodes     map[Node] error
	Edges     map[Node] []*Node
	 sync.RWMutex
}


func (g *Graph) AddNode(n Node){
	g.Lock()
	if g.nodes == nil{
		g.nodes = make(map[Node] error)
	}
	g.nodes[n] = nil
	g.Unlock()
}

func (g *Graph) AddEdge(fromNode, toNode *Node){
	g.Lock()
	if g.Edges == nil{
		g.Edges = make(map[Node] []*Node)
	}
	g.Edges[*fromNode] = append(g.Edges[*fromNode], toNode)
	g.Unlock()
}



func (g *Graph) String() {
	g.RLock()
	s := ""
	for k, _ :=range g.nodes{
		s += k.String() + " -> "
		near := g.Edges[k]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.RUnlock()
}