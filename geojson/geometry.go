package geojson

type Geometry[G GeometryTypes] struct {
	Type   string    `json:"type" bson:"type"`
	Coords G         `json:"coordinates" bson:"coordinates"`
	BBox   []float64 `json:"bbox,omitempty" bson:"bbox,omitempty"`
}

func (g Geometry[G]) SetBBox(bbox []float64) Geometry[G] {
	g.BBox = bbox
	return g
}

// TODO implement json marshaller/unmarshaller

const (
	point           = "Point"
	multiPoint      = "MultiPoint"
	lineString      = "LineString"
	multiLineString = "MultiLineString"
	polygon         = "Polygon"
	multiPolygon    = "MultiPolygon"
)

type GeometryTypes interface {
	Point | MultiPoint | LineString | MultiLineString |
		Polygon | MultiPolygon
}

type Point []float64
type MultiPoint [][]float64
type LineString [][]float64
type MultiLineString [][][]float64
type Polygon [][][]float64
type MultiPolygon [][][][]float64

func NewPoint(coords []float64) Geometry[Point] {
	return Geometry[Point]{
		Type:   point,
		Coords: coords,
	}
}

func NewMultiPoint(coords [][]float64) Geometry[MultiPoint] {
	return Geometry[MultiPoint]{
		Type:   multiPoint,
		Coords: coords,
	}
}

func NewLineString(coords [][]float64) Geometry[LineString] {
	return Geometry[LineString]{
		Type:   lineString,
		Coords: coords,
	}
}

func NewMultiLineString(coords [][][]float64) Geometry[MultiLineString] {
	return Geometry[MultiLineString]{
		Type:   multiLineString,
		Coords: coords,
	}
}

func NewPolygon(coords [][][]float64) Geometry[Polygon] {
	return Geometry[Polygon]{
		Type:   polygon,
		Coords: coords,
	}
}

func NewMultiPolygon(coords [][][][]float64) Geometry[MultiPolygon] {
	return Geometry[MultiPolygon]{
		Type:   multiPolygon,
		Coords: coords,
	}
}
