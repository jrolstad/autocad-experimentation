# dwgsum

`dwgsum` is a Go CLI that summarizes AutoCAD DWG files by converting them to DXF and then parsing core CAD entities.

## Features

- Dependency diagnostics with `doctor`
- DWG to DXF conversion via `dwg2dxf` (LibreDWG)
- Summary output in JSON or text:
  - file/version info
  - layer/block/entity counts
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
```

## Notes

- Room inference is heuristic and currently relies on closed polyline boundaries.
- Drawings without closed boundaries will still return entity-level summaries with warnings.

