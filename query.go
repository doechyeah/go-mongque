package mongque

import "go.mongodb.org/mongo-driver/bson"

// Expr is the common interface implemented by every filter expression —
// field predicates and top-level logical combinators alike. It is what
// Filter()/FilterD() produce and what And/Or/Nor consume.
type Expr interface {
	Filter() bson.M
	FilterD() bson.D
}

// Op is a single field-level operator expression, e.g. {$gt: 5}.
// The V type parameter is a phantom that constrains which FieldExpr[V]
// an operator may attach to; value holds the operand as any, so element
// and geospatial operands whose type differs from V can share the type.
type Op[V any] struct {
	key   string
	value any
}

// FieldExpr is the fluent builder for one field. Operator methods each
// return a new FieldExpr with the operator appended (immutable chaining).
type FieldExpr[V any] struct {
	name string
	ops  []Op[V]
}

// Field starts a fluent builder for the given field name. The type
// parameter V is the field's value type; comparison methods are checked
// against it.
func Field[V any](name string) FieldExpr[V] {
	return FieldExpr[V]{name: name}
}

// add returns a copy of f with op appended to a freshly allocated slice,
// so branched builders never share a backing array.
func (f FieldExpr[V]) add(op Op[V]) FieldExpr[V] {
	ops := make([]Op[V], len(f.ops)+1)
	copy(ops, f.ops)
	ops[len(f.ops)] = op
	return FieldExpr[V]{name: f.name, ops: ops}
}

// Filter renders the field predicate as a bson.M:
// {name: {op1: v1, op2: v2, ...}}.
func (f FieldExpr[V]) Filter() bson.M {
	m := make(bson.M, len(f.ops))
	for _, o := range f.ops {
		m[o.key] = o.value
	}
	return bson.M{f.name: m}
}

// FilterD renders the field predicate as a bson.D, preserving operator
// insertion order.
func (f FieldExpr[V]) FilterD() bson.D {
	d := make(bson.D, len(f.ops))
	for i, o := range f.ops {
		d[i] = bson.E{Key: o.key, Value: o.value}
	}
	return bson.D{{Key: f.name, Value: d}}
}

// NewFilter and NewFilterD assume each Expr renders exactly one top-level
// key per Filter()/FilterD() call, which every current Expr satisfies.
// NewFilter merges several expressions into one document (implicit AND
// over distinct fields). On a duplicate top-level key it falls back to
// And(...) so no clause is silently dropped.
func NewFilter(exprs ...Expr) bson.M {
	out := make(bson.M)
	for _, e := range exprs {
		for k, v := range e.Filter() {
			if _, dup := out[k]; dup {
				return And(exprs...).Filter()
			}
			out[k] = v
		}
	}
	return out
}

// NewFilterD merges several expressions into one ordered document.
// On a duplicate key it falls back to And(...).
func NewFilterD(exprs ...Expr) bson.D {
	seen := make(map[string]bool)
	var out bson.D
	for _, e := range exprs {
		for _, el := range e.FilterD() {
			if seen[el.Key] {
				return And(exprs...).FilterD()
			}
			seen[el.Key] = true
			out = append(out, el)
		}
	}
	return out
}
