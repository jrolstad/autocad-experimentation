package infer

import (
	"testing"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

func TestCountDoorsFromInserts(t *testing.T) {
	d := &model.Drawing{
		Inserts: []model.Insert{
			{Name: "A-DOOR-SINGLE"},
			{Name: "A-DOOR-DOUBLE"},
			{Name: "WINDOW"},
		},
		Texts: []model.Text{
			{Value: "DOOR TAG"},
		},
	}

	got := CountDoors(d)
	if got != 2 {
		t.Fatalf("expected 2 doors from inserts, got %d", got)
	}
}

func TestCountDoorsFallbackToText(t *testing.T) {
	d := &model.Drawing{
		Texts: []model.Text{
			{Value: "DOOR 101"},
			{Value: "INTERIOR DOOR"},
			{Value: "OUTDOOR PATIO"},
		},
	}

	got := CountDoors(d)
	if got != 2 {
		t.Fatalf("expected 2 doors from text fallback, got %d", got)
	}
}
