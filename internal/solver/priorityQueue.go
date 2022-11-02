package solver

import (
	"Ricochet-Robot-Solver/internal/types"
	"sort"
)

// An item is something we manage in a priority queue.
type item struct {
	Value      types.BoardState
	HAndGScore uint8
}

type priorityQueue []item

func (pq *priorityQueue) len() int { return len(*pq) }

func (pq *priorityQueue) push(item item) {
	*pq = append(*pq, item)
	pq.sortQueue()
}

func (pq *priorityQueue) sortQueue() {
	// We want pop to give us the lowest, not highest, priority, so we use less than here.
	sort.Slice(*pq, func(i, j int) bool {
		return calcPriority((*pq)[i].HAndGScore) < calcPriority((*pq)[j].HAndGScore)
	})
}

func (pq *priorityQueue) pop() item {
	item := (*pq)[0]
	*pq = append((*pq)[:0], (*pq)[1:]...)
	return item
}

func combineHAndGScore(gScore uint8, hScore uint8) uint8 {
	return (hScore << 5) | gScore
}

func calcPriority(hAndGScore byte) uint8 {
	return getHScore(hAndGScore) + getGScore(hAndGScore)
}

func getHScore(hAndGScore byte) uint8 {
	return (hAndGScore & (7 << 5)) >> 5
}

func getGScore(hAndGScore byte) uint8 {
	return hAndGScore & 31
}
