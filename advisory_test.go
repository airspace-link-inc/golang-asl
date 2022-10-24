package asl

/*
import (
	"fmt"
	"testing"
	"time"

	"github.com/peterstace/simplefeatures/geom"
	"github.com/stretchr/testify/assert"
)

func init() {
	ls, err := geom.NewLineString(geom.NewSequence([]float64{
		-85.38391174432391, 38.782187748582714,
		34.782107793927395, 32.085243181703234,
		-77.03652118394466, 38.897601427166194,
		-85.38391174432391, 38.782187748582714,
	}, geom.DimXY))

	if err != nil {
		panic(err)
	}

	exampleGeom1, err = geom.NewPolygon([]geom.LineString{ls})
	if err != nil {
		panic(err)
	}
}

var (
	exampleGeom1 geom.Polygon
)

func TestMarshalJSON(mainTest *testing.T) {
	randStr1, randStr2 := "sh08dajsid", "asjfasf"

	testCases := []struct {
		name     string
		arg      Advisory
		expected []byte
	}{
		{
			name:     "base case (geometry should never be null; just know GeometryCollection isn't acceptable)",
			expected: []byte(`{"type":"Feature","geometry":{"type":"GeometryCollection","geometries":[]},"properties":{"advisoryCategory":"AdvisoryCategoryType(0)","altitudeLower":0,"altitudeUpper":0,"contactEmail":null,"contactPhone":null,"countryGeoID":"","createdBy":"","endTime":"0001-01-01T00:00:00Z","geoID":"","id":"","lastEditedBy":"","name":"","ovn":"","published":false,"referenceNumber":null,"startTime":"0001-01-01T00:00:00Z","tags":null,"timezoneName":"","url":null,"version":0}}`),
		},
		{
			name: "just geom",
			arg: Advisory{
				Geometry: exampleGeom1.AsGeometry(),
			},
			expected: []byte(`{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-85.38391174432391,38.782187748582714],[34.782107793927395,32.085243181703234],[-77.03652118394466,38.897601427166194],[-85.38391174432391,38.782187748582714]]]},"properties":{"advisoryCategory":"AdvisoryCategoryType(0)","altitudeLower":0,"altitudeUpper":0,"contactEmail":null,"contactPhone":null,"countryGeoID":"","createdBy":"","endTime":"0001-01-01T00:00:00Z","geoID":"","id":"","lastEditedBy":"","name":"","ovn":"","published":false,"referenceNumber":null,"startTime":"0001-01-01T00:00:00Z","tags":null,"timezoneName":"","url":null,"version":0}}`),
		},
		{
			name: "only fields",
			arg: Advisory{
				ID:               "heo2",
				GeoID:            "uqhroh3o",
				AdvisoryCategory: Admin,
				Name:             "oj2oiejqwo",
				Tags:             []string{"asdh8", "a9ud9"},
				AltitudeLower:    100,
				AltitudeUpper:    200,
				StartTime:        time.Date(1902, 10, 2, 3, 5, 6, 11, time.UTC),
				EndTime:          time.Date(2011, 11, 8, 1, 7, 3, 22, time.UTC),
				TimezoneName:     "ajisodjaosd",
				ContactEmail:     &randStr1,
				ContactPhone:     &randStr2,
				CountryGeoID:     "kasojdiad",
				URLString:        &randStr1,
				ReferenceNumber:  &randStr1,
				CreatedBy:        "asjidh8ajd0ip",
				LastEditedBy:     "h89123h1",
				Published:        true,
				OVN:              "128h3910jidoqwnoq",
				Version:          2,
			},
			expected: []byte(`{"type":"Feature","geometry":{"type":"GeometryCollection","geometries":[]},"properties":{"advisoryCategory":"admin","altitudeLower":100,"altitudeUpper":200,"contactEmail":"sh08dajsid","contactPhone":"asjfasf","countryGeoID":"kasojdiad","createdBy":"asjidh8ajd0ip","endTime":"2011-11-08T01:07:03.000000022Z","geoID":"uqhroh3o","id":"heo2","lastEditedBy":"h89123h1","name":"oj2oiejqwo","ovn":"128h3910jidoqwnoq","published":true,"referenceNumber":"sh08dajsid","startTime":"1902-10-02T03:05:06.000000011Z","tags":["asdh8","a9ud9"],"timezoneName":"ajisodjaosd","url":"sh08dajsid","version":2}}`),
		},
		{
			name: "both geom and fields",
			arg: Advisory{
				ID:               "heo2",
				GeoID:            "uqhroh3o",
				Geometry:         exampleGeom1.AsGeometry(),
				AdvisoryCategory: Admin,
				Name:             "oj2oiejqwo",
				Tags:             []string{"asdh8", "a9ud9"},
				AltitudeLower:    100,
				AltitudeUpper:    200,
				StartTime:        time.Date(1902, 10, 2, 3, 5, 6, 11, time.UTC),
				EndTime:          time.Date(2011, 11, 8, 1, 7, 3, 22, time.UTC),
				TimezoneName:     "ajisodjaosd",
				ContactEmail:     &randStr1,
				ContactPhone:     &randStr2,
				CountryGeoID:     "kasojdiad",
				URLString:        &randStr2,
				ReferenceNumber:  &randStr1,
				CreatedBy:        "asjidh8ajd0ip",
				LastEditedBy:     "h89123h1",
				Published:        true,
				OVN:              "128h3910jidoqwnoq",
				Version:          2,
			},
			expected: []byte(`{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-85.38391174432391,38.782187748582714],[34.782107793927395,32.085243181703234],[-77.03652118394466,38.897601427166194],[-85.38391174432391,38.782187748582714]]]},"properties":{"advisoryCategory":"admin","altitudeLower":100,"altitudeUpper":200,"contactEmail":"sh08dajsid","contactPhone":"asjfasf","countryGeoID":"kasojdiad","createdBy":"asjidh8ajd0ip","endTime":"2011-11-08T01:07:03.000000022Z","geoID":"uqhroh3o","id":"heo2","lastEditedBy":"h89123h1","name":"oj2oiejqwo","ovn":"128h3910jidoqwnoq","published":true,"referenceNumber":"sh08dajsid","startTime":"1902-10-02T03:05:06.000000011Z","tags":["asdh8","a9ud9"],"timezoneName":"ajisodjaosd","url":{"Scheme":"https","Opaque":"","User":null,"Host":"airspacelink.com:199","Path":"/asdasd","RawPath":"","OmitHost":false,"ForceQuery":false,"RawQuery":"","Fragment":"","RawFragment":""},"version":2}}`),
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		actual, actualErr := tc.arg.MarshalJSON()
		fmt.Println(string(actual))
		if t.Nil(actualErr, tc.name+" should never return an error") {
			t.Equal(tc.expected, actual, tc.name)
		}
	}
}

*/
