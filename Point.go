package nci_geo_lib

import (
	"math"
)

type Point struct {
	X float64
	Y float64
}

//返回两点之间的距离
func (a *Point) Distance(b *Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

//点p离线段[a-b]距离最短的点和最短距离
func MinDistancePoint(a, b, p *Point) (Point, float64) {
	var minDistPoint Point
	var minDist float64
	t := ((a.Y-p.Y)*(b.Y-a.Y) - (a.X-p.X)*(a.X-b.X)) / ((a.X-b.X)*(a.X-b.X) - (a.Y-b.Y)*(b.Y-a.Y))
	x := (a.X-b.X)*t + a.X
	y := (a.Y-b.Y)*t + a.Y
	//理论最短距离点
	minDistPoint = Point{x, y}

	//超出线段范围，没有在a-b上的垂直点,则最短距离为a or b
	if minDistPoint.X > math.Max(a.X, b.X) || minDistPoint.X < math.Min(a.X, b.X) || minDistPoint.Y > math.Max(a.Y, b.Y) || minDistPoint.Y < math.Min(a.Y, b.Y) {
		if a.Distance(p) < b.Distance(p) {
			minDistPoint = *a
		} else {
			minDistPoint = *b
		}
	}
	minDist = p.Distance(&minDistPoint)
	return minDistPoint, minDist
}

//点到线段最短距离
func DistanceFromSegment(P, A, B *Point) float64 {
	if A.X == B.X && A.Y == B.Y {
		return P.Distance(A)
	}
	len2 := (B.X-A.X)*(B.X-A.X) + (B.Y-A.Y)*(B.Y-A.Y)
	r := ((P.X-A.X)*(B.X-A.X) + (P.Y-A.Y)*(B.Y-A.Y)) / len2

	if r <= 0.0 {
		return P.Distance(A)
	}
	if r >= 1.0 {
		return P.Distance(B)
	}
	s := ((A.Y-P.Y)*(B.X-A.X) - (A.X-P.X)*(B.Y-A.Y)) / len2
	return math.Abs(s) * math.Sqrt(len2)
}

//点P在线段A-B的投影点
func ProjectPoint(P, A, B *Point) Point {
	if *P == *A || *P == *B {
		return *P
	}
	r := ProjectionFactor(P, A, B)
	return Point{A.X + r*(B.X-A.X), A.Y + r*(B.Y-A.Y)}
}

func ProjectionFactor(P, P0, P1 *Point) float64 {
	if *P == *P0 {
		return 0.0
	}
	if *P == *P1 {
		return 1.0
	}
	dx := P1.X - P0.X
	dy := P1.Y - P0.Y
	len := dx*dx + dy*dy
	r := ((P.X-P0.X)*dx + (P.Y-P0.Y)*dy) / len
	return r
}

//点到线段最短距离的点,P靠A更近
func ClosestPoint(P, A, B *Point) (Point,bool) {
	factor := ProjectionFactor(P, A, B)
	if factor > 0 && factor < 1 {
		return ProjectPoint(P, A, B),P.Distance(A) < P.Distance(B)
	}
	dist0 := A.Distance(P)
	dist1 := B.Distance(P)
	if dist0 < dist1 {
		return *A,true
	} else {
		return *B,false
	}
}

//点P在线段A-B的右侧?
func IsOnRight(P, A, B *Point) bool {
	res := (B.X-A.X)*(P.Y-A.Y) - (P.X-A.X)*(B.Y-A.Y)
	return res < 0
}
