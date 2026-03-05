package report

import (
	"strings"
	"testing"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

func TestBuildSummaryWithDoorCount(t *testing.T) {
	d := &model.Drawing{
		Layers:         map[string]struct{}{"A-WALL": {}},
		Blocks:         map[string]struct{}{"B1": {}},
		EntitiesByType: map[string]int{"INSERT": 3},
	}
	doors := 2

	s := BuildSummary("in.dwg", 100, "tmp.dxf", d, nil, &doors)
	if s.Stats.Doors == nil {
		t.Fatal("expected doors count to be present")
	}
	if *s.Stats.Doors != 2 {
		t.Fatalf("expected doors=2, got %d", *s.Stats.Doors)
	}
}

func TestToTextIncludesDoorsWhenPresent(t *testing.T) {
	doors := 4
	var s model.Summary
	s.File.Path = "in.dwg"
	s.File.DWGVersion = "unknown"
	s.Stats.EntitiesByType = map[string]int{}
	s.Stats.Doors = &doors

	out := ToText(s)
	if !strings.Contains(out, "Doors: 4") {
		t.Fatalf("expected text output to include door count, got:\n%s", out)
	}
}
