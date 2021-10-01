package pathetic

import (
	"container/heap"
	"math"
	"math/rand"
)

type Graph struct {
	Nodes []Node
}

type Node struct {
	Edges []Edge
}

type Edge struct {
	Weight int
	Target int // index of node in graph
}

type Point struct {
	Lon float64
	Lat float64
}

type VantagePointTree struct {
	Left   *VantagePointTree
	Right  *VantagePointTree
	Center Point
	Radius float64
}

func (vp VantagePointTree) Search(p Point) Point {
	centerToPoint := distance(p, vp.Center)

	closest := vp.Center
	closestToPoint := centerToPoint

	if vp.Left != nil {
		closestLeft := vp.Left.Search(p)
		closestLeftToPoint := distance(p, closestLeft)
		if closestLeftToPoint < closestToPoint {
			closestToPoint = closestLeftToPoint
			closest = closestLeft
		}
	}

	pointToCircumference := vp.Radius - centerToPoint
	if pointToCircumference <= centerToPoint && vp.Right != nil {
		closestRight := vp.Right.Search(p)
		closestRightToPoint := distance(p, closestRight)
		if closestRightToPoint < closestToPoint {
			closest = closestRight
		}
	}
	return closest
}

type pointAroundCenter struct {
	point    Point
	distance float64
}

type pointsAroundCenter struct {
	center Point
	points []pointAroundCenter
}

func newPointsAroundCenter(center Point) pointsAroundCenter {
	return pointsAroundCenter{
		center: center,
	}
}

func distance(a Point, b Point) float64 {
	latDistance := a.Lat - b.Lat
	lonDistance := a.Lon - b.Lon
	return math.Sqrt(math.Pow(lonDistance, 2) + math.Pow(latDistance, 2))
}

func (ps *pointsAroundCenter) Len() int {
	return len(ps.points)
}

func (ps *pointsAroundCenter) Less(i, j int) bool {
	return ps.points[i].distance < ps.points[j].distance
}

func (ps *pointsAroundCenter) Swap(i, j int) {
	ps.points[i], ps.points[j] = ps.points[j], ps.points[i]
}

func (ps *pointsAroundCenter) Push(p interface{}) {
	ps.points = append(ps.points, p.(pointAroundCenter))
}

func (ps *pointsAroundCenter) Pop() interface{} {
	p := ps.points[len(ps.points)-1]
	ps.points = ps.points[:len(ps.points)-1]
	return p
}

func (ps *pointsAroundCenter) Add(point Point) {
	d := distance(ps.center, point)
	p := pointAroundCenter{
		point:    point,
		distance: d,
	}
	heap.Push(ps, p)
}

func (ps *pointsAroundCenter) Partition(point Point) ([]Point, []Point, float64) {
	// pop off heap until we have more than half the points
	left, right := []Point{}, []Point{}
	var radius float64
	if len(ps.points) == 0 {
		return nil, nil, 0
	}
	for i := 0; i < int(math.Ceil(float64(len(ps.points))/2.0)); i++ {
		next := heap.Pop(ps).(pointAroundCenter)
		radius = next.distance
		left = append(left, next.point)
	}
	// we need to keep going until we get something whose distance is greater
	for ps.Len() > 0 {
		next := heap.Pop(ps).(pointAroundCenter)
		if next.distance <= radius {
			left = append(left, next.point)
		} else {
			right = append(right, next.point)
		}
	}

	return left, right, radius
}

var selectPoint func(point []Point) int = selectRandomPoint

func selectRandomPoint(points []Point) int {
	return rand.Intn(len(points))
}

func BuildVantagePointTree(points []Point) VantagePointTree {
	c := selectPoint(points)
	center := points[c]
	// put the points into a heap sorted by closeness to center
	ps := newPointsAroundCenter(center)
	for i, p := range points {
		if c == i {
			continue
		}
		ps.Add(p)
	}
	// pop off until you end up with half - these go to left
	left, right, radius := ps.Partition(center)

	var leftTree, rightTree *VantagePointTree
	if len(left) > 0 {
		l := BuildVantagePointTree(left)
		leftTree = &l
	}
	if len(right) > 0 {
		r := BuildVantagePointTree(right)
		rightTree = &r
	}

	return VantagePointTree{
		Left:   leftTree,
		Right:  rightTree,
		Center: center,
		Radius: radius,
	}
}

func (g Graph) ShortestPath(s, t int) (int, []int) {
	// shortest distance from s to i
	weights := make([]int, len(g.Nodes))
	// previous[i] is the previous vertex along the shortest path from s to i
	previous := make([]int, len(g.Nodes))
	// whether i has been visited
	visited := make([]bool, len(g.Nodes))

	for i := range weights {
		weights[i] = math.MaxInt64
		previous[i] = -1
	}
	weights[s] = 0

	frontier := func() int {
		shortestDistanceFromS := math.MaxInt64
		closestToS := -1
		for i := range g.Nodes {
			if !visited[i] && weights[i] < shortestDistanceFromS {
				shortestDistanceFromS = weights[i]
				closestToS = i
			}
		}
		return closestToS
	}

	// while unvisited:
	// frontier = vertex with the shortest known distance from the start vertex
	// calculate distance of each neighbor of frontier from start vertex
	// for each neighbor, update shortest distance from s/previous if applicable
	for f := frontier(); f > -1; f = frontier() {
		n := g.Nodes[f]
		frontierWeight := weights[f]
		for _, edge := range n.Edges {
			pathWeight := frontierWeight + edge.Weight
			if pathWeight < weights[edge.Target] {
				weights[edge.Target] = pathWeight
				previous[edge.Target] = f
			}
		}
		visited[f] = true
	}

	pathWeight := weights[t]
	path := []int{t}
	for n := t; n != s; n = previous[n] {
		path = append(path, previous[n])
	}

	return pathWeight, path
}
