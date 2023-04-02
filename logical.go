package mongque

import "go.mongodb.org/mongo-driver/bson"

type logical string

func (l logical) set(val any) bson.M {
	return bson.M{string(l): val}
}

const (
	and logical = "$and"
	not logical = "$not"
	nor logical = "$nor"
	or  logical = "or"
)

// AndBson create a new Field for a filter with the $and logical operator with a bson.D object
func AndBson(name string, values bson.D) Field[logical] {
	return Field[logical]{name, and, values}
}

// And create a new Field for a filter with the $and logical operator
func And[T Operable](name string, values ...Field[T]) Field[logical] {
	var filters bson.D
	for _, v := range values {
		filters = append(filters, v.FilterE())
	}
	return AndBson(name, filters)
}

// NotBson create a new Field for a filter with the $not logical operator with a bson.D object
func NotBson(name string, val bson.D) Field[logical] {
	return Field[logical]{name, not, val}
}

// Not create a new Field for a filter with the $and logical operator
func Not[T Operable](name string, values ...Field[T]) Field[logical] {
	var filters bson.D
	for _, v := range values {
		filters = append(filters, v.FilterE())
	}
	return NotBson(name, filters)
}

// NorBson create a new Field for a filter with the $nor logical operator with a bson.D object
func NorBson(name string, values bson.D) Field[logical] {
	return Field[logical]{name, nor, values}
}

// NorOperable create a new Field for a filter with the $and logical operator
func Nor[T Operable](name string, values ...Field[T]) Field[logical] {
	var filters bson.D
	for _, v := range values {
		filters = append(filters, v.FilterE())
	}
	return NorBson(name, filters)
}

// OrBson create a new Field for a filter with the $or logical operator with a bson.D object
func OrBson(name string, values ...bson.D) Field[logical] {
	return Field[logical]{name, or, values}
}

// Or create a new Field for a filter with the $and logical operator
func Or[T Operable](name string, values ...Field[T]) Field[logical] {
	var filters bson.D
	for _, v := range values {
		filters = append(filters, v.FilterE())
	}
	return OrBson(name, filters)
}
