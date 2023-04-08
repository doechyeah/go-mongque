package mongque

import (
	"github.com/doechyeah/go-mongque/geojson"
	"go.mongodb.org/mongo-driver/bson"
)

type geospatial string

func (g geospatial) set(v any) bson.M {
	return bson.M{string(g): bson.M{"$geometry": v}}
}

const (
	geoIntersects geospatial = "$geoIntersects"
	geoWithin     geospatial = "$geoWithin"
	near          geospatial = "$near"
	nearSphere    geospatial = "$nearSphere"
)

func GeoIntersects(name string, val any) Field[geospatial] {
	return Field[geospatial]{name, geoIntersects, val}
}

func GeoIntersectsGeoJSON[G geojson.GeometryTypes](name string, geom geojson.Geometry[G]) Field[geospatial] {
	return Field[geospatial]{
		name:  name,
		op:    geoIntersects,
		value: geom,
	}
}

func GeoWithin(name string, val any) Field[geospatial] {
	return Field[geospatial]{name, geoWithin, val}
}

func GeoWithinGeoJSON[G geojson.GeometryTypes](name string, geom geojson.Geometry[G]) Field[geospatial] {
	return Field[geospatial]{
		name:  name,
		op:    geoWithin,
		value: geom,
	}
}

func Near(name string, val any) Field[geospatial] {
	return Field[geospatial]{name, near, val}
}

func NearGeoJSON[G geojson.GeometryTypes](name string, geom geojson.Geometry[G]) Field[geospatial] {
	return Field[geospatial]{
		name:  name,
		op:    near,
		value: geom,
	}
}

func NearSphere(name string, val any) Field[geospatial] {
	return Field[geospatial]{name, near, val}
}

func NearSphereGeoJSON[G geojson.GeometryTypes](name string, geom geojson.Geometry[G]) Field[geospatial] {
	return Field[geospatial]{
		name:  name,
		op:    nearSphere,
		value: geom,
	}
}
