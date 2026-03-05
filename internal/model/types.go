package model

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Extents struct {
	MinX float64 `json:"min_x"`
	MinY float64 `json:"min_y"`
	MaxX float64 `json:"max_x"`
	MaxY float64 `json:"max_y"`
}

type Drawing struct {
	FileVersion    string
	Layers         map[string]struct{}
	Blocks         map[string]struct{}
	Extents        Extents
	EntitiesByType map[string]int
	Polylines      []Polyline
	Lines          []Line
	Texts          []Text
	Inserts        []Insert
	Warnings       []string
}

type Polyline struct {
	Layer  string
	Closed bool
	Points []Point
}

type Line struct {
	Layer string
	Start Point
	End   Point
}

type Text struct {
	Layer  string
	Value  string
	Anchor Point
}

type Insert struct {
	Layer string
	Name  string
	Point Point
}

type Room struct {
	ID         string  `json:"id"`
	AreaSqFt   float64 `json:"area_sqft"`
	Centroid   Point   `json:"centroid"`
	Label      string  `json:"label,omitempty"`
	Confidence float64 `json:"confidence"`
}

type Summary struct {
	File struct {
		Path         string `json:"path"`
		SizeBytes    int64  `json:"size_bytes"`
		DWGVersion   string `json:"dwg_version"`
		ConvertedDXF string `json:"converted_dxf"`
	} `json:"file"`
	Stats struct {
		Layers         int            `json:"layers"`
		Blocks         int            `json:"blocks"`
		EntitiesTotal  int            `json:"entities_total"`
		EntitiesByType map[string]int `json:"entities_by_type"`
		Doors          *int           `json:"doors,omitempty"`
	} `json:"stats"`
	Extents  Extents  `json:"extents"`
	Rooms    []Room   `json:"rooms"`
	Warnings []string `json:"warnings"`
}
