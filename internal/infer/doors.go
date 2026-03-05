package infer

import (
	"strings"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

// CountDoors estimates door count from block inserts and text labels.
// It prefers INSERT references because they map more directly to instances.
func CountDoors(d *model.Drawing) int {
	count := 0
	for _, ins := range d.Inserts {
		if looksLikeDoor(ins.Name) {
			count++
		}
	}
	if count > 0 {
		return count
	}

	for _, t := range d.Texts {
		if looksLikeDoor(t.Value) {
			count++
		}
	}
	return count
}

func looksLikeDoor(v string) bool {
	u := strings.ToUpper(strings.TrimSpace(v))
	if u == "" {
		return false
	}
	if strings.Contains(u, "OUTDOOR") {
		return false
	}
	return strings.Contains(u, "DOOR")
}
