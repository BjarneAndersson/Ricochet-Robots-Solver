package priorityQueue

import (
	"Ricochet-Robot-Solver/internal/types"
	"sort"
)

// An Item is something we manage in a priority queue.
type Item struct {
	Value    types.BoardState
	Priority int
}

type PriorityQueue []Item

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Push(item Item) {
	*pq = append(*pq, item)
	sortQueue(*pq)
}

func sortQueue(pq PriorityQueue) {
	// We want Pop to give us the lowest, not highest, priority, so we use less than here.
	sort.Slice(pq, func(i, j int) bool {
		return (pq)[i].Priority < (pq)[j].Priority
	})
}

func Pop(pq *PriorityQueue) Item {
	item := (*pq)[0]
	*pq = append((*pq)[:0], (*pq)[1:]...)
	return item
}
