package mongque

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Operable defines the acceptable type of operations to set a field
type Operable interface {
	set(val any) bson.M
	comparator | logical
}

// Field is the base object type to generate the filters.
type Field[T Operable] struct {
	name  string
	op    T
	value any
}

// Filter converts the field to a bson.M filter
func (f Field[Operable]) Filter() bson.M {
	return bson.M{
		f.name: f.op.set(f.value),
	}
}

// FilterD converts the field to bson.D filter
func (f Field[Operable]) FilterD() bson.D {
	return bson.D{
		f.FilterE(),
	}
}

// FilterE converts the field to a bson.E object. #Warning this is not a valid filter and not to be used by itself. Use within a bson.D object.
func (f Field[Operable]) FilterE() bson.E {
	return bson.E{
		Key: f.name, Value: f.op.set(f.value),
	}
}

// SetName sets a the field name for a filter.
func (f Field[Operable]) SetName(name string) Field[Operable] {
	f.name = name
	return f
}

// SetValue sets the field value for a filter
func (f Field[Operable]) SetValue(value string) Field[Operable] {
	f.value = value
	return f
}

// NewFilter generates a new filter with fields. Returns a bson.M object
func NewFilter[T Operable](fields ...Field[T]) bson.M {
	filter := make(bson.M)
	for _, f := range fields {
		filter[f.name] = f.op.set(f.value)
	}
	return filter
}

// NewFilterD generates a new filter with fields. Returns a bson.D object
func NewFilterD[T Operable](fields ...Field[T]) bson.D {
	filter := bson.D{}
	for _, f := range fields {
		filter = append(filter, f.FilterE())
	}
	return filter
}
