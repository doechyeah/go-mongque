package mongque

import (
	"github.com/doechyeah/go-mongque/geojson"
	"go.mongodb.org/mongo-driver/bson"
)

// geo wraps a geometry as {op: {$geometry: g}} and appends it.
func (f FieldExpr[V]) geo(op string, g geojson.GeometryArg) FieldExpr[V] {
	return f.add(Op[V]{op, bson.M{"$geometry": g}})
}

// GeoWithin appends a $geoWithin predicate against a GeoJSON geometry.
func (f FieldExpr[V]) GeoWithin(g geojson.GeometryArg) FieldExpr[V] {
	return f.geo("$geoWithin", g)
}

// GeoIntersects appends a $geoIntersects predicate against a GeoJSON geometry.
func (f FieldExpr[V]) GeoIntersects(g geojson.GeometryArg) FieldExpr[V] {
	return f.geo("$geoIntersects", g)
}

// Near appends a $near predicate against a GeoJSON geometry.
func (f FieldExpr[V]) Near(g geojson.GeometryArg) FieldExpr[V] {
	return f.geo("$near", g)
}

// NearSphere appends a $nearSphere predicate against a GeoJSON geometry.
func (f FieldExpr[V]) NearSphere(g geojson.GeometryArg) FieldExpr[V] {
	return f.geo("$nearSphere", g)
}
