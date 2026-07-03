package mongque

// Standalone operator constructors. These build Op[V] values usable on
// their own (for example inside Not) and are the primitives the FieldExpr
// methods delegate to.

// Eq builds an $eq operator.
func Eq[V any](v V) Op[V] { return Op[V]{"$eq", v} }

// Ne builds a $ne operator.
func Ne[V any](v V) Op[V] { return Op[V]{"$ne", v} }

// Gt builds a $gt operator.
func Gt[V any](v V) Op[V] { return Op[V]{"$gt", v} }

// Gte builds a $gte operator.
func Gte[V any](v V) Op[V] { return Op[V]{"$gte", v} }

// Lt builds a $lt operator.
func Lt[V any](v V) Op[V] { return Op[V]{"$lt", v} }

// Lte builds a $lte operator.
func Lte[V any](v V) Op[V] { return Op[V]{"$lte", v} }

// In builds an $in operator over the given values.
func In[V any](vs ...V) Op[V] { return Op[V]{"$in", vs} }

// Nin builds a $nin operator over the given values.
func Nin[V any](vs ...V) Op[V] { return Op[V]{"$nin", vs} }

// Eq appends an $eq comparison against the field's value type.
func (f FieldExpr[V]) Eq(v V) FieldExpr[V] { return f.add(Eq(v)) }

// Ne appends a $ne comparison.
func (f FieldExpr[V]) Ne(v V) FieldExpr[V] { return f.add(Ne(v)) }

// Gt appends a $gt comparison.
func (f FieldExpr[V]) Gt(v V) FieldExpr[V] { return f.add(Gt(v)) }

// Gte appends a $gte comparison.
func (f FieldExpr[V]) Gte(v V) FieldExpr[V] { return f.add(Gte(v)) }

// Lt appends a $lt comparison.
func (f FieldExpr[V]) Lt(v V) FieldExpr[V] { return f.add(Lt(v)) }

// Lte appends a $lte comparison.
func (f FieldExpr[V]) Lte(v V) FieldExpr[V] { return f.add(Lte(v)) }

// In appends an $in comparison over the given values.
func (f FieldExpr[V]) In(vs ...V) FieldExpr[V] { return f.add(In(vs...)) }

// Nin appends a $nin comparison over the given values.
func (f FieldExpr[V]) Nin(vs ...V) FieldExpr[V] { return f.add(Nin(vs...)) }
