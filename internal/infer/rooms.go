package infer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jrolstad/autocad-experimentation/internal/geom"
	"github.com/jrolstad/autocad-experimentation/internal/model"
)

type Config struct {
	Epsilon         float64
	MinRoomAreaSqFt float64
	MaxRoomAreaSqFt float64
	LayerDenylist   map[string]struct{}
	NearLabelDist   float64
}

func DefaultConfig() Config {
	return Config{
		Epsilon:         0.02,
		MinRoomAreaSqFt: 25,
		MaxRoomAreaSqFt: 2000,
		LayerDenylist: map[string]struct{}{
			"DIM":       {},
			"ANNO":      {},
			"HATCH":     {},
			"DEFPOINTS": {},
		},
		NearLabelDist: 15.0,
	}
}

func InferRooms(d *model.Drawing, cfg Config) []model.Room {
	var rooms []model.Room
	for _, pl := range d.Polylines {
		if isDeniedLayer(pl.Layer, cfg.LayerDenylist) {
			continue
		}
		if len(pl.Points) < 3 {
			continue
		}
		if !pl.Closed && !geom.IsClosed(pl.Points, cfg.Epsilon) {
			continue
		}
		poly := geom.EnsureClosed(pl.Points, cfg.Epsilon)
		area := geom.PolygonArea(poly)
		if area < cfg.MinRoomAreaSqFt || area > cfg.MaxRoomAreaSqFt {
			continue
		}
		c := geom.PolygonCentroid(poly)
		label := nearestLabel(d.Texts, c, poly, cfg.NearLabelDist)
		conf := confidence(area, label, len(poly), cfg)
		rooms = append(rooms, model.Room{
			ID:         fmt.Sprintf("R%d", len(rooms)+1),
			AreaSqFt:   round2(area),
			Centroid:   c,
			Label:      label,
			Confidence: conf,
		})
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].AreaSqFt > rooms[j].AreaSqFt
	})
	for i := range rooms {
		rooms[i].ID = fmt.Sprintf("R%d", i+1)
	}
	return rooms
}

func nearestLabel(texts []model.Text, centroid model.Point, poly []model.Point, maxDist float64) string {
	best := ""
	bestDist := 1e18
	for _, t := range texts {
		if strings.TrimSpace(t.Value) == "" {
			continue
		}
		if geom.PointInPolygon(t.Anchor, poly) {
			return cleanLabel(t.Value)
		}
		d := geom.Dist(centroid, t.Anchor)
		if d <= maxDist && d < bestDist {
			bestDist = d
			best = t.Value
		}
	}
	return cleanLabel(best)
}

func confidence(area float64, label string, nPts int, cfg Config) float64 {
	c := 0.45
	if label != "" {
		c += 0.2
	}
	if nPts >= 4 {
		c += 0.15
	}
	mid := (cfg.MinRoomAreaSqFt + cfg.MaxRoomAreaSqFt) / 2
	if area > cfg.MinRoomAreaSqFt && area < cfg.MaxRoomAreaSqFt && area <= mid*1.25 {
		c += 0.2
	}
	if c > 0.99 {
		c = 0.99
	}
	return round2(c)
}

func isDeniedLayer(layer string, deny map[string]struct{}) bool {
	u := strings.ToUpper(strings.TrimSpace(layer))
	for d := range deny {
		if strings.Contains(u, d) {
			return true
		}
	}
	return false
}

func cleanLabel(s string) string {
	s = strings.ReplaceAll(s, "\\P", " ")
	s = strings.TrimSpace(s)
	return s
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
