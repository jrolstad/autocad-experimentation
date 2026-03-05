package geom

import (
	"testing"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

func TestPolygonAreaAndCentroid(t *testing.T) {
	p := []model.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}
	area := PolygonArea(p)
	if area != 100 {
		t.Fatalf("expected area 100, got %v", area)
	}
	c := PolygonCentroid(p)
	if c.X != 5 || c.Y != 5 {
		t.Fatalf("expected centroid (5,5), got (%v,%v)", c.X, c.Y)
	}
}

func TestPointInPolygon(t *testing.T) {
	poly := []model.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}
	if !PointInPolygon(model.Point{X: 5, Y: 5}, poly) {
		t.Fatal("expected point to be inside")
	}
	if PointInPolygon(model.Point{X: 15, Y: 5}, poly) {
		t.Fatal("expected point to be outside")
	}
}
