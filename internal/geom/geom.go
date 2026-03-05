package geom

import (
	"math"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

func PolygonArea(pts []model.Point) float64 {
	if len(pts) < 3 {
		return 0
	}
	var s float64
	for i := 0; i < len(pts); i++ {
		j := (i + 1) % len(pts)
		s += pts[i].X*pts[j].Y - pts[j].X*pts[i].Y
	}
	return math.Abs(s) / 2
}

func PolygonCentroid(pts []model.Point) model.Point {
	if len(pts) == 0 {
		return model.Point{}
	}
	var sx, sy float64
	for _, p := range pts {
		sx += p.X
		sy += p.Y
	}
	return model.Point{X: sx / float64(len(pts)), Y: sy / float64(len(pts))}
}

func IsClosed(pts []model.Point, eps float64) bool {
	if len(pts) < 3 {
		return false
	}
	a := pts[0]
	b := pts[len(pts)-1]
	return Dist(a, b) <= eps
}

func EnsureClosed(pts []model.Point, eps float64) []model.Point {
	if len(pts) < 3 {
		return pts
	}
	if IsClosed(pts, eps) {
		return pts[:len(pts)-1]
	}
	return pts
}

func PointInPolygon(p model.Point, poly []model.Point) bool {
	if len(poly) < 3 {
		return false
	}
	inside := false
	j := len(poly) - 1
	for i := 0; i < len(poly); i++ {
		pi := poly[i]
		pj := poly[j]
		intersects := ((pi.Y > p.Y) != (pj.Y > p.Y)) &&
			(p.X < (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y+1e-12)+pi.X)
		if intersects {
			inside = !inside
		}
		j = i
	}
	return inside
}

func Dist(a, b model.Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}
