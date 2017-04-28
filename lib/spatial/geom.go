package spatial

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var endianness = binary.LittleEndian

type Point [2]float64

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) Y() float64 {
	return p[1]
}

type GeomType uint32

const (
	GeomTypePoint      GeomType = 1
	GeomTypeLineString          = 2
	GeomTypePolygon             = 3
	GeomTypeInvalid
)

type Geom struct {
	typ GeomType
	g   interface{}
}

func NewGeom(g interface{}) (Geom, error) {
	switch g.(type) {
	case [2]float64:
		return Geom{typ: GeomTypePoint, g: g}, nil
	case [][2]float64:
		return Geom{typ: GeomTypeLineString, g: g}, nil
	case [][][2]float64:
		return Geom{typ: GeomTypePolygon, g: g}, nil
	default:
		return Geom{}, errors.New("unknown input geom type")
	}
}

func (g *Geom) UnmarshalJSON(buf []byte) error {
	wg := struct {
		Type        string
		Coordinates json.RawMessage
	}{}
	err := json.Unmarshal(buf, &wg)
	if err != nil {
		return err
	}

	switch strings.ToLower(wg.Type) {
	case "point":
		g.typ = GeomTypePoint
		var p Point
		if err = json.Unmarshal(wg.Coordinates, &p); err != nil {
			return err
		}
		g.g = p
	case "linestring":
		g.typ = GeomTypeLineString
		var ls []Point
		if err = json.Unmarshal(wg.Coordinates, &ls); err != nil {
			return err
		}
		g.g = ls
	case "polygon":
		g.typ = GeomTypePolygon
		var poly [][]Point
		if err = json.Unmarshal(wg.Coordinates, &poly); err != nil {
			return err
		}
		g.g = poly
	default:
		return fmt.Errorf("unsupported geometry type: %s", wg.Type)
	}
	return nil
}

func (g *Geom) Typ() GeomType {
	return g.typ
}

func (g *Geom) Point() (Point, error) {
	geom, ok := g.g.(Point)
	if !ok {
		return Point{}, errors.New("geometry is not a Point")
	}
	return geom, nil
}

func (g *Geom) LineString() ([]Point, error) {
	geom, ok := g.g.([]Point)
	if !ok {
		return nil, errors.New("geometry is not a LineString")
	}
	return geom, nil
}

func (g *Geom) Polygon() ([][]Point, error) {
	geom, ok := g.g.([][]Point)
	if !ok {
		return nil, errors.New("geometry is not a Polygon")
	}
	return geom, nil
}

type Feature struct {
	Type     string
	Props    map[string]interface{} `json:"properties"`
	Geometry Geom
}

func (f *Feature) MarshalWKB() ([]byte, error) {
	if endianness != binary.LittleEndian {
		return nil, errors.New("only little endian is supported")
	}
	var buf bytes.Buffer
	binary.Write(&buf, endianness, uint8(1))         // little endian
	binary.Write(&buf, endianness, f.Geometry.Typ()) // geometry type

	switch f.Geometry.Typ() {
	case GeomTypePoint:
		p, err := f.Geometry.Point()
		if err != nil {
			return nil, err
		}
		wkbWritePoint(&buf, p)
	case GeomTypeLineString:
		ls, err := f.Geometry.LineString()
		if err != nil {
			return nil, err
		}
		wkbWriteLineString(&buf, ls)
	}
	return buf.Bytes(), nil
}

func (f *Feature) Properties() map[string]interface{} {
	return f.Props
}

type FeatureCollection struct {
	Features []Feature
}
