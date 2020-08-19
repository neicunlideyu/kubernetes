package util

import (
	"fmt"
)

type NodeScore struct {
	Host  string
	Score int
}

type MinHeap struct {
	size          int
	priorityNodes []NodeScore
}

func NewMinHeap() *MinHeap {
	return &MinHeap{}
}

func Parent(i int) int {
	if i == 0 {
		return 0
	}
	return (i - 1) / 2
}

func LeftChild(i int) int {
	return 2*i + 1
}

func (heap *MinHeap) IsEmpty() bool {
	return heap.size == 0
}

func (heap *MinHeap) GetSize() int {
	return heap.size
}

func (heap *MinHeap) GetMinVal() (NodeScore, error) {
	if heap.IsEmpty() {
		return NodeScore{}, fmt.Errorf("GetMinVal failed, MinHeap is empty")
	}
	return heap.priorityNodes[0], nil
}

func siftDown(heap *MinHeap, parI int) {
	var minI int
	for {
		leftI := LeftChild(parI)
		switch {
		case leftI+1 > heap.size:
			return
		case leftI+2 > heap.size:
			if heap.priorityNodes[parI].Score > heap.priorityNodes[leftI].Score {
				heap.priorityNodes[parI], heap.priorityNodes[leftI] = heap.priorityNodes[leftI],
					heap.priorityNodes[parI]
			}
			return
		case heap.priorityNodes[leftI].Score <= heap.priorityNodes[leftI+1].Score:
			minI = leftI
		case heap.priorityNodes[leftI].Score > heap.priorityNodes[leftI+1].Score:
			minI = leftI + 1
		}

		if heap.priorityNodes[parI].Score > heap.priorityNodes[minI].Score {
			heap.priorityNodes[parI], heap.priorityNodes[minI] = heap.priorityNodes[minI],
				heap.priorityNodes[parI]
		}
		parI = minI
	}
}

func (heap *MinHeap) SiftUp(priorityNode NodeScore) {
	heap.priorityNodes = append(heap.priorityNodes, priorityNode)
	parI := Parent(heap.size)
	childI := heap.size

	for heap.priorityNodes[parI].Score > heap.priorityNodes[childI].Score {
		heap.priorityNodes[parI], heap.priorityNodes[childI] = heap.priorityNodes[childI],
			heap.priorityNodes[parI]
		childI = parI
		parI = Parent(parI)
	}
	heap.size++
}

func (heap *MinHeap) SiftDown() (NodeScore, error) {
	minVal, err := heap.GetMinVal()
	if err != nil {
		return minVal, err
	}

	heap.size--
	heap.priorityNodes[0], heap.priorityNodes[heap.size] = heap.priorityNodes[heap.size], NodeScore{}

	siftDown(heap, 0)
	return minVal, nil
}

func (heap *MinHeap) Replace(priorityNode NodeScore) error {
	_, err := heap.GetMinVal()
	if err != nil {
		return fmt.Errorf("GetMinVal error: %+v", err)
	}
	heap.priorityNodes[0] = priorityNode

	siftDown(heap, 0)
	return nil
}

// Get the top m highest score nodes
func TopHighestScoreNodes(priorityNodes []NodeScore, m int) (resp []NodeScore) {
	h := NewMinHeap()

	for i := 0; i < m; i++ {
		h.SiftUp(priorityNodes[i])
	}

	minVal, _ := h.GetMinVal()
	for _, v := range priorityNodes[m:] {
		if v.Score > minVal.Score {
			h.Replace(v)
			minVal, _ = h.GetMinVal()
		}
	}
	for i := 0; i < m; i++ {
		v, _ := h.SiftDown()
		resp = append(resp, v)
	}
	return
}
