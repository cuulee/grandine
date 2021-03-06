package geojson

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/thomersch/grandine/lib/spatial"
)

type Codec struct{}

func (c *Codec) Decode(r io.Reader, fc *spatial.FeatureCollection) error {
	var gjfc featureColl
	err := json.NewDecoder(r).Decode(&gjfc)
	if err != nil {
		return err
	}
	fc.Features = append(fc.Features, gjfc.Features...)
	if srid, ok := ogcSRID[gjfc.CRS.Properties.Name]; ok {
		if len(fc.SRID) == 0 {
			fc.SRID = srid
		}
		if len(srid) != 0 && fc.SRID != srid {
			return errors.New("incompatible projections: ")
		}
	}
	return nil
}

func (c *Codec) Encode(w io.Writer, fc *spatial.FeatureCollection) error {
	geojsonFC := featureColl{
		Type: "FeatureCollection",
		// TODO: set CRS
		Features: fc.Features,
	}
	return json.NewEncoder(w).Encode(&geojsonFC)
}

func (c *Codec) Extensions() []string {
	return []string{"geojson", "json"}
}

type featureColl struct {
	Type string `json:"type"`
	CRS  struct {
		Type       string
		Properties struct {
			Name string
		}
	} `json:"crs"`
	Features []spatial.Feature
}
