package astar

import (
	"fmt"
	"math"
	"testing"
)

// Labels for a 4x4 square are:
//
//       0  1  2  3
//       4  5  6  7
//       8  9 10 11
//      12 13 14 15

type square struct {
	size   int
	points [][]float64
}

func newSquare(size int) *square {
	s := &square{size: size}
	s.points = make([][]float64, size)
	for r := 0; r < size; r++ {
		s.points[r] = make([]float64, size)
	}
	return s
}

func (s *square) coordinates(label int) (r, c int) {
	r = label / s.size
	c = label - r*s.size
	return
}

func (s *square) contains(x int) bool { return x >= 0 && x < s.size }

func (s *square) label(r, c int) int { return c + r*s.size }

func (s *square) put(label int, value float64) {
	r, c := s.coordinates(label)
	s.points[r][c] = value
}

func (s *square) Adjacent(label int) map[int]float64 {
	adjacency := make(map[int]float64)
	r, c := s.coordinates(label)
	offsets := []int{-1, 0, 1}
	for i := range offsets {
		nr := r + offsets[i]
		if !s.contains(nr) {
			continue
		}
		for j := range offsets {
			nc := c + offsets[j]
			if !s.contains(nc) {
				continue
			}
			// Omit self.
			if r == nr && c == nc {
				continue
			}
			adjacency[s.label(nr, nc)] = s.points[nr][nc]
		}
	}
	return adjacency
}

func (s *square) Estimate(from, to int) float64 {
	fr, fc := s.coordinates(from)
	tr, tc := s.coordinates(to)
	// Simply return the manhattan distance regardless of cost.
	return math.Abs(float64(tr-fr)) + math.Abs(float64(tc-fc))
}

func compareToSlice(expected, given []int) bool {
	if len(given) != len(expected) {
		return false
	}
	for i, v := range expected {
		if given[i] != v {
			return false
		}
	}
	return true
}

func TestRouteThroughEmptySquare(t *testing.T) {
	s := newSquare(4)
	route := Find(s, 0, 15)
	if route == nil {
		t.Error()
	}
	if !compareToSlice([]int{0, 5, 10, 15}, route) {
		t.Error()
	}
}

func TestSquareRouteWithSingleBlockedPoint(t *testing.T) {
	s := newSquare(4)
	s.put(10, math.MaxFloat64)
	route := Find(s, 0, 15)
	if route == nil {
		t.Error()
	}
	if !compareToSlice([]int{0, 5, 6, 11, 15}, route) {
		fmt.Println(route)
		t.Error()
	}
}

func TestSquareRouteWithBarrier(t *testing.T) {
	s := newSquare(4)
	s.put(10, math.MaxFloat64)
	s.put(11, math.MaxFloat64)
	route := Find(s, 0, 15)
	if route == nil {
		t.Error()
	}
	if !compareToSlice([]int{0, 5, 9, 14, 15}, route) {
		t.Error()
	}
	s.put(9, math.MaxFloat64)
	route = Find(s, 0, 15)
	if route == nil {
		t.Error()
	}
	if !compareToSlice([]int{0, 5, 8, 13, 14, 15}, route) {
		t.Error()
	}
	s.put(8, math.MaxFloat64)
	if Find(s, 0, 15) != nil {
		t.Error()
	}
}

// Graph is:
//
//      0 --- 1 --- 2 --- 3
//              \-- 4 --/

type graph struct {
	kind      int                     // One of the constants.
	adjacency map[int]map[int]float64 // For each vertex, a map of its adjacent vertices.
}

const (
	directed int = iota
	undirected
)

func newGraph(kind int) (*graph, error) {
	return &graph{kind: kind, adjacency: make(map[int]map[int]float64)}, nil
}

func (g *graph) edge(i, j int, cost float64) {
	vertex, known := g.adjacency[i]
	if !known {
		vertex = make(map[int]float64)
		g.adjacency[i] = vertex
	}
	vertex[j] = cost
	// Ensure the j vertex is also in the adjacency map.
	vertex, known = g.adjacency[j]
	if !known {
		vertex = make(map[int]float64)
		g.adjacency[j] = vertex
	}
	if g.kind == undirected {
		// For undirected graphs the same cost applies from j to i.
		vertex[i] = cost
	}
}
func (g *graph) Adjacent(i int) map[int]float64 { return g.adjacency[i] }

func (g *graph) Estimate(i, j int) float64 { return 1 }

func TestGraphRoute(t *testing.T) {
	g, _ := newGraph(undirected)
	g.edge(0, 1, 1)
	g.edge(1, 2, 1)
	g.edge(2, 3, 1)
	g.edge(1, 4, 1)
	g.edge(4, 3, 1)
	if !compareToSlice([]int{0, 1, 2, 3}, Find(g, 0, 3)) {
		t.Error()
	}
	g.edge(1, 2, 5) // Cost seen at the branch.
	if !compareToSlice([]int{0, 1, 4, 3}, Find(g, 0, 3)) {
		t.Error()
	}
	g.edge(1, 2, 1)
	g.edge(2, 3, 5) // Cost only seen next to the goal.
	if !compareToSlice([]int{0, 1, 4, 3}, Find(g, 0, 3)) {
		t.Error()
	}
}
