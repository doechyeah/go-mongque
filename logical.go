package mongque

import "go.mongodb.org/mongo-driver/bson"

// logicalExpr renders a top-level logical operator over full expressions:
// {$and: [expr, ...]}.
type logicalExpr struct {
	op    string
	exprs []Expr
}

// Filter renders the logical combinator as a bson.M.
func (l logicalExpr) Filter() bson.M {
	if len(l.exprs) == 0 {
		return bson.M{}
	}
	arr := make(bson.A, len(l.exprs))
	for i, e := range l.exprs {
		arr[i] = e.Filter()
	}
	return bson.M{l.op: arr}
}

// FilterD renders the logical combinator as a bson.D.
func (l logicalExpr) FilterD() bson.D {
	if len(l.exprs) == 0 {
		return bson.D{}
	}
	arr := make(bson.A, len(l.exprs))
	for i, e := range l.exprs {
		arr[i] = e.FilterD()
	}
	return bson.D{{Key: l.op, Value: arr}}
}

// And joins expressions with logical AND: {$and: [...]}. With no
// arguments it renders an empty (match-all) filter.
func And(exprs ...Expr) Expr { return logicalExpr{"$and", exprs} }

// Or joins expressions with logical OR: {$or: [...]}. With no
// arguments it renders an empty (match-all) filter.
func Or(exprs ...Expr) Expr { return logicalExpr{"$or", exprs} }

// Nor joins expressions with logical NOR: {$nor: [...]}. With no
// arguments it renders an empty (match-all) filter.
func Nor(exprs ...Expr) Expr { return logicalExpr{"$nor", exprs} }

// Not inverts a single operator on this field: {field: {$not: {op: v}}}.
// It accepts a comparison Op[V], keeping the negation type-checked
// against the field's value type.
func (f FieldExpr[V]) Not(op Op[V]) FieldExpr[V] {
	return f.add(Op[V]{"$not", bson.M{op.key: op.value}})
}
