package priorityQueue

import (
	"../types"
	"sort"
)

// An Item is something we manage in a priority queue.
type Item struct {
	Value      types.BoardState
	RobotOrder uint8
	HAndGScore uint8
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
		return calcPriority((pq)[i].HAndGScore) < calcPriority((pq)[j].HAndGScore)
	})
}

func Pop(pq *PriorityQueue) Item {
	item := (*pq)[0]
	*pq = append((*pq)[:0], (*pq)[1:]...)
	return item
}

func CombineHAndGScore(gScore uint8, hScore uint8) uint8 {
	return uint8((hScore << 5) | gScore)
}

func calcPriority(hAndGScore byte) uint8 {
	return GetHScore(hAndGScore) + GetGScore(hAndGScore)
}

func GetHScore(hAndGScore byte) uint8 {
	return (hAndGScore & (7 << 5)) >> 5
}

func GetGScore(hAndGScore byte) uint8 {
	return hAndGScore & 31
}
