package mongque

// Exists appends an $exists predicate: {field: {$exists: b}}.
func (f FieldExpr[V]) Exists(b bool) FieldExpr[V] {
	return f.add(Op[V]{"$exists", b})
}

// Type appends a $type predicate over one or more BSON type aliases.
// A single type renders as a scalar ({$type: "string"}); multiple render
// as an array ({$type: ["string", "int"]}).
func (f FieldExpr[V]) Type(types ...string) FieldExpr[V] {
	var v any
	if len(types) == 1 {
		v = types[0]
	} else {
		v = types
	}
	return f.add(Op[V]{"$type", v})
}
