package dxfread

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jrolstad/autocad-experimentation/internal/model"
)

type pair struct {
	code int
	val  string
}

func ReadFile(path string) (*model.Drawing, error) {
	pairs, err := readPairs(path)
	if err != nil {
		return nil, err
	}
	d := &model.Drawing{
		Layers:         map[string]struct{}{},
		Blocks:         map[string]struct{}{},
		EntitiesByType: map[string]int{},
	}

	parseHeaderVars(pairs, d)
	parseSections(pairs, d)
	return d, nil
}

func readPairs(path string) ([]pair, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open dxf: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	var lines []string
	for sc.Scan() {
		lines = append(lines, strings.TrimRight(sc.Text(), "\r\n"))
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read dxf: %w", err)
	}
	if len(lines)%2 != 0 {
		lines = lines[:len(lines)-1]
	}

	pairs := make([]pair, 0, len(lines)/2)
	for i := 0; i < len(lines); i += 2 {
		code, err := strconv.Atoi(strings.TrimSpace(lines[i]))
		if err != nil {
			continue
		}
		pairs = append(pairs, pair{code: code, val: strings.TrimSpace(lines[i+1])})
	}
	return pairs, nil
}

func parseHeaderVars(pairs []pair, d *model.Drawing) {
	for i := 0; i < len(pairs)-1; i++ {
		if pairs[i].code == 9 && strings.EqualFold(pairs[i].val, "$ACADVER") {
			for j := i + 1; j < len(pairs); j++ {
				if pairs[j].code == 1 {
					d.FileVersion = pairs[j].val
					break
				}
				if pairs[j].code == 9 {
					break
				}
			}
		}
		if pairs[i].code == 9 && strings.EqualFold(pairs[i].val, "$EXTMIN") {
			minX, minY := readXYVar(pairs, i+1)
			d.Extents.MinX = minX
			d.Extents.MinY = minY
		}
		if pairs[i].code == 9 && strings.EqualFold(pairs[i].val, "$EXTMAX") {
			maxX, maxY := readXYVar(pairs, i+1)
			d.Extents.MaxX = maxX
			d.Extents.MaxY = maxY
		}
	}
}

func readXYVar(pairs []pair, start int) (float64, float64) {
	var x, y float64
	for i := start; i < len(pairs); i++ {
		if pairs[i].code == 9 {
			break
		}
		switch pairs[i].code {
		case 10:
			x = parseFloat(pairs[i].val)
		case 20:
			y = parseFloat(pairs[i].val)
		}
	}
	return x, y
}

func parseSections(pairs []pair, d *model.Drawing) {
	for i := 0; i < len(pairs); i++ {
		if pairs[i].code == 0 && strings.EqualFold(pairs[i].val, "SECTION") {
			secName := ""
			if i+1 < len(pairs) && pairs[i+1].code == 2 {
				secName = strings.ToUpper(pairs[i+1].val)
			}
			j := i + 2
			for ; j < len(pairs); j++ {
				if pairs[j].code == 0 && strings.EqualFold(pairs[j].val, "ENDSEC") {
					break
				}
			}
			parseSection(secName, pairs[i+2:j], d)
			i = j
		}
	}
}

func parseSection(name string, pairs []pair, d *model.Drawing) {
	switch name {
	case "TABLES":
		parseTablesSection(pairs, d)
	case "BLOCKS":
		parseBlocksSection(pairs, d)
	case "ENTITIES":
		parseEntitiesSection(pairs, d)
	}
}

func parseTablesSection(pairs []pair, d *model.Drawing) {
	recs := records(pairs)
	for _, r := range recs {
		if r.typ == "LAYER" {
			layer := attr(r.attrs, 2)
			if layer != "" {
				d.Layers[layer] = struct{}{}
			}
		}
	}
}

func parseBlocksSection(pairs []pair, d *model.Drawing) {
	recs := records(pairs)
	for _, r := range recs {
		if r.typ == "BLOCK" {
			name := attr(r.attrs, 2)
			if name != "" {
				d.Blocks[name] = struct{}{}
			}
		}
	}
}

func parseEntitiesSection(pairs []pair, d *model.Drawing) {
	recs := records(pairs)
	var activeLegacy *model.Polyline

	for _, r := range recs {
		d.EntitiesByType[r.typ]++
		switch r.typ {
		case "LINE":
			ln := model.Line{
				Layer: attr(r.attrs, 8),
				Start: model.Point{X: floatAttr(r.attrs, 10), Y: floatAttr(r.attrs, 20)},
				End:   model.Point{X: floatAttr(r.attrs, 11), Y: floatAttr(r.attrs, 21)},
			}
			if ln.Layer != "" {
				d.Layers[ln.Layer] = struct{}{}
			}
			d.Lines = append(d.Lines, ln)
		case "LWPOLYLINE":
			pl := model.Polyline{
				Layer:  attr(r.attrs, 8),
				Closed: intAttr(r.attrs, 70)&1 == 1,
				Points: readLWPoints(r.attrs),
			}
			if pl.Layer != "" {
				d.Layers[pl.Layer] = struct{}{}
			}
			if len(pl.Points) >= 2 {
				d.Polylines = append(d.Polylines, pl)
			}
		case "POLYLINE":
			pl := &model.Polyline{
				Layer:  attr(r.attrs, 8),
				Closed: intAttr(r.attrs, 70)&1 == 1,
			}
			activeLegacy = pl
		case "VERTEX":
			if activeLegacy != nil {
				activeLegacy.Points = append(activeLegacy.Points, model.Point{
					X: floatAttr(r.attrs, 10),
					Y: floatAttr(r.attrs, 20),
				})
			}
		case "SEQEND":
			if activeLegacy != nil && len(activeLegacy.Points) >= 2 {
				if activeLegacy.Layer != "" {
					d.Layers[activeLegacy.Layer] = struct{}{}
				}
				d.Polylines = append(d.Polylines, *activeLegacy)
			}
			activeLegacy = nil
		case "TEXT", "MTEXT":
			t := model.Text{
				Layer:  attr(r.attrs, 8),
				Value:  strings.TrimSpace(attr(r.attrs, 1)),
				Anchor: model.Point{X: floatAttr(r.attrs, 10), Y: floatAttr(r.attrs, 20)},
			}
			if t.Layer != "" {
				d.Layers[t.Layer] = struct{}{}
			}
			if t.Value != "" {
				d.Texts = append(d.Texts, t)
			}
		case "INSERT":
			ins := model.Insert{
				Layer: attr(r.attrs, 8),
				Name:  strings.TrimSpace(attr(r.attrs, 2)),
				Point: model.Point{X: floatAttr(r.attrs, 10), Y: floatAttr(r.attrs, 20)},
			}
			if ins.Layer != "" {
				d.Layers[ins.Layer] = struct{}{}
			}
			if ins.Name != "" {
				d.Inserts = append(d.Inserts, ins)
			}
		}
	}
}

type record struct {
	typ   string
	attrs []pair
}

func records(pairs []pair) []record {
	var out []record
	i := 0
	for i < len(pairs) {
		if pairs[i].code != 0 {
			i++
			continue
		}
		typ := strings.ToUpper(strings.TrimSpace(pairs[i].val))
		j := i + 1
		for ; j < len(pairs); j++ {
			if pairs[j].code == 0 {
				break
			}
		}
		out = append(out, record{typ: typ, attrs: pairs[i+1 : j]})
		i = j
	}
	return out
}

func attr(attrs []pair, code int) string {
	for _, a := range attrs {
		if a.code == code {
			return a.val
		}
	}
	return ""
}

func floatAttr(attrs []pair, code int) float64 {
	for _, a := range attrs {
		if a.code == code {
			return parseFloat(a.val)
		}
	}
	return 0
}

func intAttr(attrs []pair, code int) int {
	for _, a := range attrs {
		if a.code == code {
			n, _ := strconv.Atoi(strings.TrimSpace(a.val))
			return n
		}
	}
	return 0
}

func readLWPoints(attrs []pair) []model.Point {
	var xs, ys []float64
	for _, a := range attrs {
		switch a.code {
		case 10:
			xs = append(xs, parseFloat(a.val))
		case 20:
			ys = append(ys, parseFloat(a.val))
		}
	}
	n := len(xs)
	if len(ys) < n {
		n = len(ys)
	}
	points := make([]model.Point, 0, n)
	for i := 0; i < n; i++ {
		points = append(points, model.Point{X: xs[i], Y: ys[i]})
	}
	return points
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}
