package astar

import (
	"container/heap"
	"math"
)

// Interface defines the problem space for the A* algorithm. Whatever the
// implemention, every point in the problem space must have a unique integer
// label.
//
type Interface interface {

	// Adjacent returns a map of the adjacent points with their associated
	// cost.
	Adjacent(int) (map[int]float64, error)

	// Estimate returns the estimated cost of moving from the given point to the
	// goal. Known as the heuristic function (h) in the literature.
	Estimate(int, int) float64
}

type point struct {
	label    int     // The label assigned in the problem space.
	index    int     // The position in the open set or -1 if not.
	cost     float64 // The cost from the start until this point (g).
	estimate float64 // The estimated cost from the start through to the goal (f).
}

type minheap []*point

func (h minheap) Len() int {
	return len(h)
}
func (h minheap) Less(i, j int) bool {
	return h[i].estimate < h[j].estimate
}
func (h minheap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *minheap) Push(x interface{}) {
	p := x.(*point)
	p.index = len(*h)
	*h = append(*h, p)
}
func (h *minheap) Pop() interface{} {
	slice := *h
	n := len(slice) - 1
	p := slice[n]
	*h = slice[:n]
	p.index = -1
	return p
}

// Find the route from the start to the goal. If no route is feasible then
// return nil.
//
func Find(problem Interface, start, goal int) []int {
	//
	// The open set is a min heap of points based on the estimated cost.
	//
	open := make(minheap, 0)
	heap.Init(&open)
	//
	// All the known points, so that each point struct is only made once,
	// when it is first seen.
	//
	known := make(map[int]*point)
	//
	// The route as a map of point labels referring to the predecessor.
	//
	predecessor := make(map[int]int)
	//
	// Make the starting point.
	//
	startPoint := &point{
		label:    start,
		estimate: problem.Estimate(start, goal),
	}
	known[start] = startPoint
	heap.Push(&open, startPoint)
	//
	// Until the open set is empty ...
	//
	for open.Len() > 0 {
		//
		// Remove the lowest estimate point from the open set.
		//
		current := heap.Pop(&open).(*point)
		if current.label == goal {
			return route(goal, predecessor)
		}
		//
		// Examine each neighbour to find the lowest cost movement from
		// the current point.
		//
		adjacent, _ := problem.Adjacent(current.label)
		for k, v := range adjacent {
			p, ok := known[k]
			if !ok {
				p = &point{label: k, index: -1, cost: math.MaxFloat64}
				known[k] = p
			}
			cost := v + current.cost
			if cost < p.cost {
				predecessor[k] = current.label
				p.cost = cost
				p.estimate = p.cost + problem.Estimate(k, goal)
				if p.index == -1 {
					heap.Push(&open, p)
				} else {
					heap.Fix(&open, p.index)
				}
			}
		}
	}
	return nil
}

func route(goal int, predecessor map[int]int) []int {
	//
	// Unwind the route backwards from the goal.
	//
	backwards := []int{goal}
	label, ok := predecessor[goal]
	for ok {
		backwards = append(backwards, label)
		label, ok = predecessor[label]
	}
	//
	// Reverse this.
	//
	forwards := make([]int, len(backwards))
	j := len(backwards)
	for i := range backwards {
		j--
		forwards[i] = backwards[j]
	}
	return forwards
}
