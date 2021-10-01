package pathetic

import (
	"reflect"
	"testing"
)

func TestPointsAroundCenter(t *testing.T) {
	tests := []struct {
		points         []Point
		pointsToSelect []int
		expected       Point
	}{
		{
			points: []Point{
				{-2.5, -3},
				{-2, -2.5},
				{0, -4},
				{0, 0},
				{0.5, -1.5},
				{1, 0},
			},
			pointsToSelect: []int{3, 1, 0, 0, 0, 0},
			expected:       Point{-2, -2.5},
		},
	}
	for _, test := range tests {
		selectPoint = func(points []Point) int {
			next := test.pointsToSelect[0]
			test.pointsToSelect = test.pointsToSelect[1:]
			return next
		}
		vp := BuildVantagePointTree(test.points)
		got := vp.Search(Point{-2, -2.6})
		if test.expected != got {
			t.Errorf("expected %v, got %v", test.expected, got)
		}
	}
}

func TestPathetic(t *testing.T) {
	graph := Graph{
		Nodes: []Node{
			{
				Edges: []Edge{
					{
						Weight: 2,
						Target: 1,
					},
					{
						Weight: 4,
						Target: 2,
					},
				},
			},
			{
				Edges: []Edge{
					{
						Weight: 7,
						Target: 3,
					},
					{
						Weight: 1,
						Target: 2,
					},
				},
			},
			{
				Edges: []Edge{
					{
						Weight: 3,
						Target: 4,
					},
				},
			},
			{
				Edges: []Edge{
					{
						Weight: 1,
						Target: 5,
					},
				},
			},
			{
				Edges: []Edge{
					{
						Weight: 2,
						Target: 3,
					},
					{
						Weight: 5,
						Target: 5,
					},
				},
			},
			{
				Edges: []Edge{},
			},
		},
	}
	var source, target int = 0, 5
	distance, path := graph.ShortestPath(source, target)
	expectedDistance := 9
	expectedPath := []int{5, 3, 4, 2, 1, 0}
	if distance != expectedDistance {
		t.Errorf("got %d, expected %d", distance, expectedDistance)
	}
	if !reflect.DeepEqual(path, expectedPath) {
		t.Errorf("got %v, expected %v", path, expectedPath)
	}
}
