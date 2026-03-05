## Go DWG Summarizer MVP (Open-Source, Windows-First, Room Inference)

### Summary
Build a **CLI-first Go app** that summarizes AutoCAD DWG files by:
1. Converting `DWG -> DXF` via **LibreDWG** tools.
2. Parsing DXF entities in Go.
3. Producing:
   - reliable structural summary (metadata/entity/layer/block counts),
   - heuristic **room/space inference** for architectural floor plans.

This is the fastest feasible open-source path with acceptable DWG support and clear upgrade options later.

### Scope
In scope:
- Single executable CLI.
- Local file processing (no upload service).
- JSON + human-readable output.
- Heuristic room detection for 2D floor plans.
- Windows-first packaging/instructions.

Out of scope (v1):
- Perfect semantic understanding of all CAD conventions.
- 3D solids/BIM semantics.
- cloud service/UI.
- commercial SDK integration.

---

## Architecture

### Pipeline
1. **Input validation**
   - Accept `.dwg` path.
   - Verify file exists/readable.
2. **DWG conversion**
   - Run LibreDWG converter (`dwg2dxf`) to temp DXF.
   - Capture converter stderr/stdout for diagnostics.
3. **DXF parse stage**
   - Read entities, layers, blocks, polylines, lines, arcs, circles, text.
4. **Geometry normalization**
   - Convert entities to 2D segments/loops in model space.
   - Normalize units/epsilon tolerance.
5. **Room inference**
   - Detect closed loops likely representing rooms.
   - Filter by area thresholds and containment rules.
   - Attach nearest text labels (if present).
6. **Summary generation**
   - Metadata + counts + inferred rooms.
7. **Output**
   - JSON to stdout/file.
   - Optional text report.

---

## Project Layout (greenfield)
- `cmd/dwgsum/main.go` - CLI entrypoint.
- `internal/convert/` - LibreDWG invocation + error handling.
- `internal/dxfread/` - DXF parsing adapter.
- `internal/geom/` - segment graph, loop closure, area calc.
- `internal/infer/` - room inference heuristics.
- `internal/report/` - JSON/text formatting.
- `internal/model/` - shared types.
- `testdata/` - sample DWG/DXF fixtures + expected summaries.

---

## Public Interfaces (initial)

### CLI
- `dwgsum summarize --in <file.dwg> [--out summary.json] [--format json|text] [--keep-temp] [--verbose]`
- `dwgsum doctor`  
  Checks LibreDWG binaries and prints environment diagnostics.

### Output schema (`v1`)
```json
{
  "file": {
    "path": "...",
    "size_bytes": 0,
    "dwg_version": "unknown|Rxx",
    "converted_dxf": "..."
  },
  "stats": {
    "layers": 0,
    "blocks": 0,
    "entities_total": 0,
    "entities_by_type": {}
  },
  "extents": {
    "min_x": 0,
    "min_y": 0,
    "max_x": 0,
    "max_y": 0
  },
  "rooms": [
    {
      "id": "R1",
      "area_sqft": 0,
      "centroid": {"x": 0, "y": 0},
      "label": "optional",
      "confidence": 0.0
    }
  ],
  "warnings": []
}
```

---

## Room Inference Heuristics (decision-complete for v1)
1. Build a planar graph from line/polyline/arc approximations.
2. Find closed loops with tolerance snapping (`epsilon`).
3. Reject loops with:
   - area below minimum threshold,
   - extreme aspect ratio/noise signatures,
   - obvious non-room layers (configurable denylist).
4. Resolve nested loops:
   - outer loop as candidate space shell,
   - holes treated as exclusions.
5. Label matching:
   - nearest text/MTEXT point-in-polygon or within distance threshold.
6. Confidence scoring:
   - closure quality + area plausibility + label presence + boundary continuity.
7. Emit low-confidence candidates with warnings rather than dropping silently.

Config defaults (YAML/flags):
- `epsilon = 0.02` drawing units
- `min_room_area_sqft = 25`
- `max_room_area_sqft = 2000`
- `layer_denylist = ["DIM", "ANNO", "HATCH", "DEFPOINTS"]` (case-insensitive)

---

## Implementation Phases

1. **Bootstrap + diagnostics**
   - Initialize Go module.
   - Implement `doctor` command and converter process wrapper.
   - Add robust process timeout + actionable error messages.
2. **Reliable summary core**
   - Parse DXF into canonical entity model.
   - Implement stats/extents/layer/block reporting.
   - Ship stable JSON output.
3. **Room inference MVP**
   - Add loop detection + area computation.
   - Add text association and confidence scoring.
   - Add warnings for ambiguous drawings.
4. **Hardening**
   - Improve tolerance handling.
   - Add config overrides and `--verbose` traces.
   - Expand fixture set with edge cases.
5. **Windows-first packaging**
   - Document LibreDWG install paths and PATH requirements.
   - Provide PowerShell install/check script.
   - Add CI job for Windows tests.

---

## Testing and Acceptance Criteria

### Unit tests
- Converter command construction and failure modes.
- Geometry: closure, area, point-in-polygon, nesting.
- Label matching and confidence computation.

### Fixture/integration tests
- Known DWG floorplan -> expected room count range and key stats.
- Non-architectural DWG -> no false “high confidence” room output.
- Corrupted/unreadable DWG -> graceful error with remediation.

### CLI contract tests
- `--format json` outputs valid schema.
- Non-zero exit codes on conversion/parse failures.
- `doctor` reports missing dependencies correctly.

### Acceptance criteria
- For representative floorplan files, inferred room count is directionally correct and no crashes.
- Summary always returns deterministic stats even when room inference is low confidence.
- Clear diagnostics when LibreDWG unavailable or conversion fails.

---

## Risks and Mitigations
- DWG version incompatibility: fallback warnings + fixture coverage by version.
- Heuristic room accuracy variability: confidence scoring + explicit warnings.
- Windows tooling friction: `doctor` + install script + documented PATH checks.

---

## Assumptions and Defaults
- Input files are mostly **2D architectural floor plans**.
- Open-source-only constraint remains (no ODA/commercial SDK in v1).
- Windows is the primary runtime for first release.
- Conversion path `DWG -> DXF` is acceptable for MVP latency/accuracy.
- If layer naming is inconsistent, inference still runs but confidence may decrease.
