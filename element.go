package mongque

// Exists appends an $exists predicate: {field: {$exists: b}}.
func (f FieldExpr[V]) Exists(b bool) FieldExpr[V] {
	return f.add(Op[V]{"$exists", b})
}

// Type appends a $type predicate over one or more BSON type aliases.
// A single type renders as a scalar ({$type: "string"}); multiple render
// as an array ({$type: ["string", "int"]}); none is a no-op.
func (f FieldExpr[V]) Type(types ...string) FieldExpr[V] {
	switch len(types) {
	case 0:
		return f
	case 1:
		return f.add(Op[V]{"$type", types[0]})
	default:
		return f.add(Op[V]{"$type", types})
	}
}
