package infer

import (
	"testing"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

func TestInferRoomsFromClosedPolyline(t *testing.T) {
	d := &model.Drawing{
		Polylines: []model.Polyline{
			{
				Layer:  "A-WALL",
				Closed: true,
				Points: []model.Point{
					{X: 0, Y: 0},
					{X: 20, Y: 0},
					{X: 20, Y: 15},
					{X: 0, Y: 15},
				},
			},
		},
		Texts: []model.Text{
			{
				Layer:  "A-ANNO",
				Value:  "BEDROOM",
				Anchor: model.Point{X: 10, Y: 8},
			},
		},
	}
	cfg := DefaultConfig()
	cfg.MinRoomAreaSqFt = 10
	cfg.MaxRoomAreaSqFt = 1000
	rooms := InferRooms(d, cfg)
	if len(rooms) != 1 {
		t.Fatalf("expected 1 room, got %d", len(rooms))
	}
	if rooms[0].Label != "BEDROOM" {
		t.Fatalf("expected room label BEDROOM, got %q", rooms[0].Label)
	}
}

func TestInferRoomsDeniedLayer(t *testing.T) {
	d := &model.Drawing{
		Polylines: []model.Polyline{
			{
				Layer:  "ANNO-ROOM",
				Closed: true,
				Points: []model.Point{
					{X: 0, Y: 0},
					{X: 20, Y: 0},
					{X: 20, Y: 20},
					{X: 0, Y: 20},
				},
			},
		},
	}
	cfg := DefaultConfig()
	cfg.MinRoomAreaSqFt = 10
	cfg.MaxRoomAreaSqFt = 1000
	rooms := InferRooms(d, cfg)
	if len(rooms) != 0 {
		t.Fatalf("expected 0 rooms on denied layer, got %d", len(rooms))
	}
}
