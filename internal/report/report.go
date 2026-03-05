package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

func BuildSummary(inPath string, size int64, convertedDXF string, d *model.Drawing, rooms []model.Room, doors *int) model.Summary {
	var s model.Summary
	s.File.Path = inPath
	s.File.SizeBytes = size
	s.File.DWGVersion = d.FileVersion
	if s.File.DWGVersion == "" {
		s.File.DWGVersion = "unknown"
	}
	s.File.ConvertedDXF = convertedDXF
	s.Stats.Layers = len(d.Layers)
	s.Stats.Blocks = len(d.Blocks)
	s.Stats.EntitiesByType = map[string]int{}
	total := 0
	for k, v := range d.EntitiesByType {
		s.Stats.EntitiesByType[k] = v
		total += v
	}
	s.Stats.EntitiesTotal = total
	s.Stats.Doors = doors
	s.Extents = d.Extents
	s.Rooms = rooms
	s.Warnings = append([]string{}, d.Warnings...)
	if len(rooms) == 0 {
		s.Warnings = append(s.Warnings, "No room-like closed polylines found in configured area thresholds.")
	}
	return s
}

func ToText(s model.Summary) string {
	var b strings.Builder
	fmt.Fprintf(&b, "File: %s\n", s.File.Path)
	fmt.Fprintf(&b, "Size: %d bytes\n", s.File.SizeBytes)
	fmt.Fprintf(&b, "DWG Version: %s\n", s.File.DWGVersion)
	fmt.Fprintf(&b, "Converted DXF: %s\n", s.File.ConvertedDXF)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Layers: %d\nBlocks: %d\nEntities: %d\n", s.Stats.Layers, s.Stats.Blocks, s.Stats.EntitiesTotal)
	if s.Stats.Doors != nil {
		fmt.Fprintf(&b, "Doors: %d\n", *s.Stats.Doors)
	}
	fmt.Fprintln(&b, "Entity Types:")
	keys := make([]string, 0, len(s.Stats.EntitiesByType))
	for k := range s.Stats.EntitiesByType {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&b, "- %s: %d\n", k, s.Stats.EntitiesByType[k])
	}
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Extents: min(%.2f, %.2f) max(%.2f, %.2f)\n", s.Extents.MinX, s.Extents.MinY, s.Extents.MaxX, s.Extents.MaxY)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Rooms: %d\n", len(s.Rooms))
	for _, r := range s.Rooms {
		if r.Label != "" {
			fmt.Fprintf(&b, "- %s: %.2f sqft, label=%q, confidence=%.2f\n", r.ID, r.AreaSqFt, r.Label, r.Confidence)
		} else {
			fmt.Fprintf(&b, "- %s: %.2f sqft, confidence=%.2f\n", r.ID, r.AreaSqFt, r.Confidence)
		}
	}
	if len(s.Warnings) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Warnings:")
		for _, w := range s.Warnings {
			fmt.Fprintf(&b, "- %s\n", w)
		}
	}
	return strings.TrimSpace(b.String())
}
