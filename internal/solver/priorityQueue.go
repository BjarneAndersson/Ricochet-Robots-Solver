package solver

import (
	"Ricochet-Robot-Solver/internal/types"
	"sort"
)

// An item is an element of a priority queue.
type item struct {
	Value      types.BoardState
	HAndGScore uint8
}

type priorityQueue []item

// len Returns the item count of the priority queue
func (pq *priorityQueue) len() int { return len(*pq) }

// push Appends the item to the priorityQueue and runs sortQueue
func (pq *priorityQueue) push(item item) {
	*pq = append(*pq, item)
	pq.sortQueue()
}

// sortQueue Sorts the priorityQueue according to the sum of item.HAndGScore
func (pq *priorityQueue) sortQueue() {
	// We want pop to give us the lowest, not highest, priority, so we use less than here.
	sort.Slice(*pq, func(i, j int) bool {
		return calcPriority((*pq)[i].HAndGScore) < calcPriority((*pq)[j].HAndGScore)
	})
}

// pop Returns the item with the least priority from the queue and removes it
func (pq *priorityQueue) pop() item {
	item := (*pq)[0]
	*pq = append((*pq)[:0], (*pq)[1:]...)
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
