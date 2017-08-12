// +build !golangclip

package spatial

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/paulsmith/gogeos/geos"
)

func (p Polygon) clipToBBox(b BBox) []Geom {
	gpoly := p.geos()

	var bboxLine []geos.Coord
	for _, pt := range NewLinesFromSegments(BBoxBorders(b.SW, b.NE))[0] {
		spew.Dump(pt)
		bboxLine = append(bboxLine, geos.NewCoord(pt.X(), pt.Y()))
	}
	bboxPoly := geos.Must(geos.NewPolygon(bboxLine))

	res, _ := bboxPoly.Intersection(gpoly)

	var resGeoms []Geom
	for _, poly := range geosToPolygons(res) {
		resGeoms = append(resGeoms, MustNewGeom(poly))
	}
	return resGeoms
}

func (p Polygon) geos() *geos.Geometry {
	rings := [][]geos.Coord{}
	for _, ring := range p {
		var rg []geos.Coord
		for _, pt := range ring {
			rg = append(rg, geos.NewCoord(pt.X(), pt.Y()))
		}
		rg = append(rg, rg[0])
		rings = append(rings, rg)
	}
	var gpoly *geos.Geometry
	if len(rings) > 1 {
		gpoly = geos.Must(geos.NewPolygon(rings[0], rings[1:]...))
	} else {
		gpoly = geos.Must(geos.NewPolygon(rings[0]))
	}
	return gpoly
}

func geosToPolygons(g *geos.Geometry) []Polygon {
	ty, _ := g.Type()
	if ty == geos.POLYGON {
		return []Polygon{geosToPolygon(g)}
	}
	var polys []Polygon
	nmax, err := g.NGeometry()
	if err != nil {
		panic(err)
	}
	for n := 0; n < nmax; n++ {
		polys = append(polys, geosToPolygon(geos.Must(g.Geometry(n))))
	}
	return polys
}

func geosToPolygon(g *geos.Geometry) Polygon {
	var (
		p    Polygon
		ring []Point
	)

	fmt.Print("{{{{ ")
	spew.Dump(g.Type())

	crds, _ := geos.Must(g.Shell()).Coords()
	for _, crd := range crds {
		ring = append(ring, Point{crd.X, crd.Y})
	}
	p = append(p, ring)

	holes, _ := g.Holes()
	for _, hole := range holes {
		crds, _ = hole.Coords()
		ring = []Point{}
		for _, crd := range crds {
			ring = append(ring, Point{crd.X, crd.Y})
		}
		p = append(p, ring)
	}
	return p
}
