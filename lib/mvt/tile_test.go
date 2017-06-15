package mvt

import (
	"fmt"
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/stretchr/testify/assert"
)

func TestTileName(t *testing.T) {
	for _, tc := range []struct {
		p        spatial.Point
		expected TileID
	}{
		{
			p:        spatial.Point{13.73630, 51.05377},
			expected: TileID{X: 8817, Y: 5481, Z: 14},
		},
		{
			p:        spatial.Point{18.39856, -33.90184},
			expected: TileID{X: 9029, Y: 9833, Z: 14},
		},
		{
			p:        spatial.Point{-54.59123, -25.59547},
			expected: TileID{X: 5707, Y: 9397, Z: 14},
		},
		{
			p:        spatial.Point{-21.94073, 64.14607},
			expected: TileID{X: 7193, Y: 4354, Z: 14},
		},
		{
			p:        spatial.Point{-31.16580, 83.65691},
			expected: TileID{X: 6773, Y: 648, Z: 14},
		},
		{
			p:        spatial.Point{-64.45649, -85.04438},
			expected: TileID{X: 5258, Y: 16380, Z: 14},
		},
	} {
		t.Run(fmt.Sprintf("%v_%v", tc.expected.X, tc.expected.Y), func(t *testing.T) {
			var fail bool
			ti := TileName(tc.p, 14)
			if ti.X != tc.expected.X {
				fail = true
			}
			if ti.Y != tc.expected.Y {
				fail = true
			}
			if ti.Z != tc.expected.Z {
				fail = true
			}
			if fail {
				t.Fatalf("invalid conversion, expected %v, got %v", tc.expected, ti)
			}
		})
	}
}

func TestTileBBox(t *testing.T) {
	for _, tc := range []struct {
		tid      TileID
		expected bbox
	}{
		{
			tid:      TileID{0, 0, 0},
			expected: bbox{spatial.Point{-180, -85.05112878}, spatial.Point{180, 85.05112878}},
		},
		{
			tid:      TileID{0, 0, 1},
			expected: bbox{spatial.Point{-180, 0}, spatial.Point{0, 85.05112878}},
		},
		{
			tid:      TileID{0, 1, 1},
			expected: bbox{spatial.Point{-180, -85.05112878}, spatial.Point{0, 0}},
		},
		{
			tid:      TileID{1, 2, 2},
			expected: bbox{spatial.Point{-90, -66.51326044}, spatial.Point{0, 0}},
		},
	} {
		t.Run(fmt.Sprintf("%v_%v_%v", tc.tid.X, tc.tid.Y, tc.tid.Z), func(t *testing.T) {
			sw, ne := tc.tid.BBox()
			bb := bbox{sw.RoundedCoords(), ne.RoundedCoords()}
			assert.Equal(t, tc.expected, bb)
		})
	}

}
