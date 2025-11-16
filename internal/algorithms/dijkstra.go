package algorithms

import (
	"container/heap"
	"math"
)

type Edge struct {
	To   int     `json:"to"`
	Cost float64 `json:"cost"`
}

type Graph struct {
	N     int       `json:"n"`
	Edges [][]Edge  `json:"edges"`
}

type Result struct {
	Dist []float64 `json:"dist"`
}

type item struct {
	node int
	dist float64
}
type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func Dijkstra(g Graph, src int) Result {
	dist := make([]float64, g.N)
	for i := 0; i < g.N; i++ {
		dist[i] = math.Inf(1)
	}
	dist[src] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, item{node: src, dist: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		if cur.dist > dist[cur.node] {
			continue
		}
		for _, e := range g.Edges[cur.node] {
			nd := cur.dist + e.Cost
			if nd < dist[e.To] {
				dist[e.To] = nd
				heap.Push(pq, item{node: e.To, dist: nd})
			}
		}
	}
	return Result{Dist: dist}
}