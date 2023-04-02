package mongque

import "go.mongodb.org/mongo-driver/bson"

// var ErrInvalidOP = errors.New("invalid operator used to build query")

type comparator string

func (c comparator) set(v any) bson.M {
	return bson.M{string(c): v}
}

const (
	eq  comparator = "$eq"
	neq comparator = "$ne"
	lte comparator = "$lte"
	lt  comparator = "$lt"
	gte comparator = "$gte"
	gt  comparator = "$gt"
	in  comparator = "$in"
	nin comparator = "$nin"
)

// Eq create a new Field for a filter with the $eq comparator
func Eq(name string, val any) Field[comparator] {
	return Field[comparator]{name, eq, val}
}

// Neq create a new Field for a filter with the $neq comparator
func Neq(name string, val any) Field[comparator] {
	return Field[comparator]{name, neq, val}
}

// Lte create a new Field for a filter with the $lte comparator
func Lte(name string, val any) Field[comparator] {
	return Field[comparator]{name, lte, val}
}

// Lt create a new Field for a filter with the $lt comparator
func Lt(name string, val any) Field[comparator] {
	return Field[comparator]{name, lt, val}
}

// Gte create a new Field for a filter with the $gte comparator
func Gte(name string, val any) Field[comparator] {
	return Field[comparator]{name, gte, val}
}

// Gt create a new Field for a filter with the $lt comparator
func Gt(name string, val any) Field[comparator] {
	return Field[comparator]{name, gte, val}
}

// In create a new Field for a filter with the $in comparator
func In(name string, val any) Field[comparator] {
	return Field[comparator]{name, in, val}
}

// Nin create a new Field for a filter with the $nin comparator
func Nin(name string, val any) Field[comparator] {
	return Field[comparator]{name, nin, val}
}
