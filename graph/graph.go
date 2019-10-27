package graph

import (
	"fmt"
	"sync"
)

type Node struct{
	name string

}

func (n *Node) String() string{
	return fmt.Sprintf("%v", n.name)
}

type Graph struct {
	nodes []*Node
	edges map[Node] []*Node
	 sync.RWMutex
}


func (g *Graph) AddNode(n *Node){
	g.Lock()
	g.nodes = append(g.nodes, n)
	g.Unlock()
}

func (g *Graph) AddEdge(fromNode, toNode *Node){
	g.Lock()
	if g.edges == nil{
		g.edges = make(map[Node] []*Node)
	}
	g.edges[*fromNode] = append(g.edges[*fromNode], toNode)
	g.Unlock()
}



func (g *Graph) String() {
	g.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.RUnlock()
}