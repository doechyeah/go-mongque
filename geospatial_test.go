package mongque

import (
	"testing"

	"github.com/doechyeah/go-mongque/geojson"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Geospatial(t *testing.T) {
	poly := geojson.NewPolygon([][][]float64{{{0, 0}, {0, 1}, {1, 1}, {0, 0}}})
	pt := geojson.NewPoint([]float64{1, 2})

	assert.Equal(t,
		bson.M{"loc": bson.M{"$geoWithin": bson.M{"$geometry": poly}}},
		Field[any]("loc").GeoWithin(poly).Filter())

	assert.Equal(t,
		bson.M{"loc": bson.M{"$geoIntersects": bson.M{"$geometry": poly}}},
		Field[any]("loc").GeoIntersects(poly).Filter())

	assert.Equal(t,
		bson.M{"loc": bson.M{"$near": bson.M{"$geometry": pt}}},
		Field[any]("loc").Near(pt).Filter())

	assert.Equal(t,
		bson.M{"loc": bson.M{"$nearSphere": bson.M{"$geometry": pt}}},
		Field[any]("loc").NearSphere(pt).Filter())
}
