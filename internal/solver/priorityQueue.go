package solver

import (
	"Ricochet-Robot-Solver/internal/types"
)

// An item is an element of a priority queue.
type item struct {
	Value      types.BoardState
	HAndGScore uint8
	index      int
}

type priorityQueue []*item

// Len Returns the item count of the priority queue
func (pq priorityQueue) Len() int { return len(pq) }

// Less compares the priority of two items
func (pq priorityQueue) Less(i, j int) bool {
	return calcPriority(pq[i].HAndGScore) < calcPriority(pq[j].HAndGScore)
}

// Swap switches the index of two items
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push Appends the item to the priorityQueue and runs sortQueue
func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and return the item with the lowest priority
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// combineHAndGScore Combines g score and h score into one byte - h score (bit: 7-5) | g score (bit: 4-0)
func combineHAndGScore(gScore uint8, hScore uint8) uint8 {
	return (hScore << 5) | gScore
}

// calcPriority Calculates the priority of the item as the sum of g and h score
func calcPriority(hAndGScore byte) uint8 {
	return getHScore(hAndGScore) + getGScore(hAndGScore)
}

// getHScore Returns the h score of an item (attributes)
func getHScore(hAndGScore byte) uint8 {
	return (hAndGScore & (7 << 5)) >> 5
}

// getGScore Returns the g score of an item (attributes)
func getGScore(hAndGScore byte) uint8 {
	return hAndGScore & 31
}
