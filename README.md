# dwgsum

`dwgsum` is a Go CLI that summarizes AutoCAD DWG files by converting them to DXF and then parsing core CAD entities.

## Features

- Dependency diagnostics with `doctor`
- DWG to DXF conversion via `dwg2dxf` (LibreDWG)
- Summary output in JSON or text:
  - file/version info
  - layer/block/entity counts
  - optional door count (`--count-doors`)
  - drawing extents
  - heuristic room inference from closed polylines

## Requirements

- Go 1.23+
- LibreDWG tools on `PATH`:
  - `dwg2dxf` (required)
  - `dwgread` (optional; checked by `doctor`)

## Build

```powershell
go build ./cmd/dwgsum
```

## Usage

```powershell
dwgsum doctor
dwgsum summarize --in "C:\path\file.dwg" --format json
dwgsum summarize --in "C:\path\file.dwg" --format text --out ".\summary.txt"
dwgsum summarize --in "C:\path\file.dwg" --format json --count-doors
```

## Door Count Option

Use `--count-doors` with `summarize` to estimate door count from door-like `INSERT` block names, with a text-label fallback.

```powershell
dwgsum summarize --in "C:\path\file.dwg" --count-doors --format json
```

When enabled, output includes `stats.doors` (JSON) or `Doors: <n>` (text).

## Notes

- Room inference is heuristic and currently relies on closed polyline boundaries.
- Drawings without closed boundaries will still return entity-level summaries with warnings.
