package externalsort

import "container/heap"

// An Item is something we manage in a priority queue.
type mergePQItem struct {
	value  string
	reader LineReader
}

type mergePriorityQueue []mergePQItem

var (
	_ heap.Interface = &mergePriorityQueue{}
)

func (m mergePriorityQueue) Len() int {
	return len(m)
}

func (m mergePriorityQueue) Less(i, j int) bool {
	return m[i].value < m[j].value
}

func (m mergePriorityQueue) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m *mergePriorityQueue) Push(x interface{}) {
	*m = append(*m, x.(mergePQItem))
}

func (m *mergePriorityQueue) Pop() interface{} {
	n := len(*m)
	item := (*m)[n-1]
	*m = (*m)[:n-1]
	return item
}
