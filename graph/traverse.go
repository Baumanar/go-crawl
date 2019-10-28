package graph

// Traverse implements the BFS traversing algorithm
func (g *Graph) Traverse(f func(*Node)) {
	g.RLock()
	q := NodeQueue{}
	q.New()
	n := g.EntryNode
	q.Enqueue(n)
	visited := make(map[*Node]bool)
	for {
		if q.IsEmpty() {
		break
		}
		node := q.Dequeue()
		visited[node] = true
		near := g.Edges[*node]

		for i := 0; i < len(near); i++ {
			j := near[i]
			if !visited[j] {
				q.Enqueue(*j)
				visited[j] = true
			}
		}
		if f != nil {
				f(node)
			}
	}
	g.RUnlock()
	}

