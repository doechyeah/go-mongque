# Phase 0 Core Field-Centric Redesign — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace go-mongque's core query builder with a typed, field-centric fluent API that emits valid MongoDB query documents.

**Architecture:** A generic `Field[V]` builder accumulates operators via method chaining and renders to `bson.M`/`bson.D`. Operators are `Op[V]` values (phantom-typed to the field's value type). Top-level `And`/`Or`/`Nor` combine full expressions through a shared `Expr` interface; `$not` is a field-level method. Everything is immutable value semantics.

**Tech Stack:** Go 1.19+ (generics), `go.mongodb.org/mongo-driver/bson`, `testify/assert`.

## Global Constraints

- Go 1.19+ (generics required).
- Only dependencies: `go.mongodb.org/mongo-driver`, `github.com/stretchr/testify` (test only). Add no new dependencies.
- Package `mongque` (root), subpackage `geojson`.
- Operator string literals live only in constructor bodies; do not export them.
- Every operator renders the exact `bson.M` MongoDB expects. Tests assert exact shapes.
- Break freely from v0.1.0 — no backward-compatibility shims.
- Reference the spec: `docs/superpowers/specs/2026-07-03-phase0-core-redesign-design.md`.

## File Map

- `query.go` — REWRITE. `Expr` interface, `Op[V]`, `FieldExpr[V]`, `Field`, `add`, `Filter`/`FilterD`, and (Task 3) `NewFilter`/`NewFilterD`.
- `comparator.go` — REWRITE. Standalone `Op[V]` comparison constructors + `FieldExpr` comparison methods.
- `logical.go` — REWRITE. `logicalExpr`, `And`/`Or`/`Nor`, field-level `Not`.
- `geospatial.go` — REWRITE. Geospatial `FieldExpr` methods.
- `element.go` — CREATE. `Exists`/`Type` methods.
- `geojson/geometry.go` — MODIFY. Add `geo()` marker method + `GeometryArg` interface.
- `query_test.go`, `comparator_test.go`, `logical_test.go`, `geospatial_test.go`, `element_test.go` — test files (one per source file).
- `README.md`, `CLAUDE.md` — MODIFY (docs).

---

### Task 1: Core types and rendering

Establishes the new foundation and removes the old `Operable` machinery. Old and new APIs share names (`Field`, `NewFilter`), so the old operator files are deleted here; operators are re-added in later tasks. Strict red-first TDD is not possible for this one foundational swap (the old types occupy the new names), so we implement the core, then assert its behavior.

**Files:**
- Rewrite: `query.go`
- Delete: `comparator.go`, `logical.go`, `geospatial.go` (old contents; re-created in later tasks)
- Rewrite: `query_test.go`

**Interfaces:**
- Produces:
  - `type Expr interface { Filter() bson.M; FilterD() bson.D }`
  - `type Op[V any] struct { key string; value any }`
  - `type FieldExpr[V any] struct { name string; ops []Op[V] }`
  - `func Field[V any](name string) FieldExpr[V]`
  - `func (f FieldExpr[V]) add(op Op[V]) FieldExpr[V]` (unexported; used by all operator methods)
  - `func (f FieldExpr[V]) Filter() bson.M`
  - `func (f FieldExpr[V]) FilterD() bson.D`

- [ ] **Step 1: Delete the old operator files**

```bash
git rm comparator.go logical.go geospatial.go
```

- [ ] **Step 2: Replace `query.go` with the new core**

```go
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
```

- [ ] **Step 3: Replace `query_test.go`**

```go
package mongque

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Field_Empty(t *testing.T) {
	assert.Equal(t, bson.M{"name": bson.M{}}, Field[string]("name").Filter())
}

func Test_Field_EmptyD(t *testing.T) {
	assert.Equal(t, bson.D{{Key: "name", Value: bson.D{}}}, Field[string]("name").FilterD())
}
```

- [ ] **Step 4: Run the tests**

Run: `go test ./...`
Expected: PASS (root package compiles with the new core; `geojson` reports no test files).

- [ ] **Step 5: Vet**

Run: `go vet ./...`
Expected: no output.

- [ ] **Step 6: Commit**

```bash
git add query.go query_test.go
git commit -m "$(printf 'refactor!: replace core with typed field-centric builder\n\nBREAKING CHANGE: removes Field[Operable] and the old operator\nconstructors. Introduces Expr, Op[V], and FieldExpr[V].\n\nCo-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>')"
```

---

### Task 2: Comparison operators

**Files:**
- Create: `comparator.go`
- Create: `comparator_test.go`

**Interfaces:**
- Consumes: `Op[V]`, `FieldExpr[V]`, `(FieldExpr[V]).add` (Task 1).
- Produces:
  - Standalone: `func Eq[V any](v V) Op[V]`, and identically `Ne`, `Gt`, `Gte`, `Lt`, `Lte`; `func In[V any](vs ...V) Op[V]`, `func Nin[V any](vs ...V) Op[V]`.
  - Methods: `(FieldExpr[V]).Eq(v V) FieldExpr[V]`, and identically `Ne`, `Gt`, `Gte`, `Lt`, `Lte`, `In(vs ...V)`, `Nin(vs ...V)`.

- [ ] **Step 1: Write the failing test**

Create `comparator_test.go`:

```go
package mongque

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Comparators(t *testing.T) {
	tests := []struct {
		name string
		got  bson.M
		want bson.M
	}{
		{"eq", Field[string]("f").Eq("x").Filter(), bson.M{"f": bson.M{"$eq": "x"}}},
		{"ne", Field[int]("f").Ne(1).Filter(), bson.M{"f": bson.M{"$ne": 1}}},
		{"gt", Field[int]("f").Gt(1).Filter(), bson.M{"f": bson.M{"$gt": 1}}},
		{"gte", Field[int]("f").Gte(1).Filter(), bson.M{"f": bson.M{"$gte": 1}}},
		{"lt", Field[int]("f").Lt(1).Filter(), bson.M{"f": bson.M{"$lt": 1}}},
		{"lte", Field[int]("f").Lte(1).Filter(), bson.M{"f": bson.M{"$lte": 1}}},
		{"in", Field[int]("f").In(1, 2).Filter(), bson.M{"f": bson.M{"$in": []int{1, 2}}}},
		{"nin", Field[int]("f").Nin(1, 2).Filter(), bson.M{"f": bson.M{"$nin": []int{1, 2}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

func Test_MultiOp(t *testing.T) {
	got := Field[int]("age").Gte(18).Lt(65).Filter()
	assert.Equal(t, bson.M{"age": bson.M{"$gte": 18, "$lt": 65}}, got)
}

func Test_FilterD_Order(t *testing.T) {
	got := Field[int]("age").Gte(18).Lt(65).FilterD()
	assert.Equal(t, bson.D{{Key: "age", Value: bson.D{
		{Key: "$gte", Value: 18},
		{Key: "$lt", Value: 65},
	}}}, got)
}

func Test_Immutable(t *testing.T) {
	base := Field[int]("age").Gt(1)
	a := base.Lt(10).Filter()
	b := base.Lt(20).Filter()
	assert.Equal(t, bson.M{"age": bson.M{"$gt": 1, "$lt": 10}}, a)
	assert.Equal(t, bson.M{"age": bson.M{"$gt": 1, "$lt": 20}}, b)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'Test_Comparators|Test_MultiOp|Test_FilterD_Order|Test_Immutable' ./...`
Expected: FAIL — compile error, `.Eq` / `.Gt` / etc. undefined on `FieldExpr`.

- [ ] **Step 3: Write the implementation**

Create `comparator.go`:

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add comparator.go comparator_test.go
git commit -m "$(printf 'feat: add typed comparison operators\n\nCo-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>')"
```

---

### Task 3: Logical operators and NewFilter

**Files:**
- Create: `logical.go`
- Modify: `query.go` (append `NewFilter`/`NewFilterD`)
- Create: `logical_test.go`

**Interfaces:**
- Consumes: `Expr`, `Op[V]`, `FieldExpr[V]`, `(FieldExpr[V]).add`, `(FieldExpr[V]).Filter`, `(FieldExpr[V]).FilterD` (Tasks 1–2).
- Produces:
  - `func And(exprs ...Expr) Expr`, `func Or(exprs ...Expr) Expr`, `func Nor(exprs ...Expr) Expr`.
  - `func (f FieldExpr[V]) Not(op Op[V]) FieldExpr[V]`.
  - `func NewFilter(exprs ...Expr) bson.M`, `func NewFilterD(exprs ...Expr) bson.D`.

- [ ] **Step 1: Write the failing test**

Create `logical_test.go`:

```go
package mongque

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Or_Heterogeneous(t *testing.T) {
	got := Or(
		Field[string]("status").Eq("active"),
		Field[int]("score").Gt(90),
	).Filter()
	assert.Equal(t, bson.M{"$or": bson.A{
		bson.M{"status": bson.M{"$eq": "active"}},
		bson.M{"score": bson.M{"$gt": 90}},
	}}, got)
}

func Test_And(t *testing.T) {
	got := And(Field[int]("a").Eq(1), Field[int]("b").Eq(2)).Filter()
	assert.Equal(t, bson.M{"$and": bson.A{
		bson.M{"a": bson.M{"$eq": 1}},
		bson.M{"b": bson.M{"$eq": 2}},
	}}, got)
}

func Test_Nor(t *testing.T) {
	got := Nor(Field[int]("a").Eq(1)).Filter()
	assert.Equal(t, bson.M{"$nor": bson.A{
		bson.M{"a": bson.M{"$eq": 1}},
	}}, got)
}

func Test_Not(t *testing.T) {
	got := Field[int]("age").Not(Gt(5)).Filter()
	assert.Equal(t, bson.M{"age": bson.M{"$not": bson.M{"$gt": 5}}}, got)
}

func Test_NewFilter_Merge(t *testing.T) {
	got := NewFilter(Field[string]("name").Eq("John"), Field[int]("score").Lte(60))
	assert.Equal(t, bson.M{
		"name":  bson.M{"$eq": "John"},
		"score": bson.M{"$lte": 60},
	}, got)
}

func Test_NewFilter_CollisionFallsBackToAnd(t *testing.T) {
	got := NewFilter(Field[int]("age").Gt(5), Field[int]("age").Lt(20))
	assert.Equal(t, bson.M{"$and": bson.A{
		bson.M{"age": bson.M{"$gt": 5}},
		bson.M{"age": bson.M{"$lt": 20}},
	}}, got)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'Test_Or_Heterogeneous|Test_And|Test_Nor|Test_Not|Test_NewFilter' ./...`
Expected: FAIL — compile error, `And`/`Or`/`Nor`/`Not`/`NewFilter` undefined.

- [ ] **Step 3: Create `logical.go`**

```go
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
	arr := make(bson.A, len(l.exprs))
	for i, e := range l.exprs {
		arr[i] = e.Filter()
	}
	return bson.M{l.op: arr}
}

// FilterD renders the logical combinator as a bson.D.
func (l logicalExpr) FilterD() bson.D {
	arr := make(bson.A, len(l.exprs))
	for i, e := range l.exprs {
		arr[i] = e.FilterD()
	}
	return bson.D{{Key: l.op, Value: arr}}
}

// And joins expressions with logical AND: {$and: [...]}.
func And(exprs ...Expr) Expr { return logicalExpr{"$and", exprs} }

// Or joins expressions with logical OR: {$or: [...]}.
func Or(exprs ...Expr) Expr { return logicalExpr{"$or", exprs} }

// Nor joins expressions with logical NOR: {$nor: [...]}.
func Nor(exprs ...Expr) Expr { return logicalExpr{"$nor", exprs} }

// Not inverts a single operator on this field: {field: {$not: {op: v}}}.
// It accepts a comparison Op[V], keeping the negation type-checked
// against the field's value type.
func (f FieldExpr[V]) Not(op Op[V]) FieldExpr[V] {
	return f.add(Op[V]{"$not", bson.M{op.key: op.value}})
}
```

- [ ] **Step 4: Append `NewFilter`/`NewFilterD` to `query.go`**

Add to the end of `query.go`:

```go
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
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add logical.go query.go logical_test.go
git commit -m "$(printf 'feat: add top-level logical operators, field-level Not, and NewFilter\n\nCo-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>')"
```

---

### Task 4: Geospatial operators

**Files:**
- Modify: `geojson/geometry.go` (add marker)
- Create: `geospatial.go`
- Create: `geospatial_test.go`

**Interfaces:**
- Consumes: `Op[V]`, `FieldExpr[V]`, `(FieldExpr[V]).add` (Task 1); `geojson.Geometry[G]`, `geojson.NewPoint`, `geojson.NewPolygon` (existing).
- Produces:
  - `geojson`: `func (g Geometry[G]) geo() {}` and `type GeometryArg interface { geo() }`.
  - `mongque`: `(FieldExpr[V]).GeoWithin`, `GeoIntersects`, `Near`, `NearSphere`, each `func(geojson.GeometryArg) FieldExpr[V]`.

- [ ] **Step 1: Write the failing test**

Create `geospatial_test.go`:

```go
package mongque

import (
	"testing"

	"github.com/doechyeah/go-mongque/geojson"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Geospatial(t *testing.T) {
	poly := geojson.NewPolygon([][][]float64{{{0, 0}, {0, 1}, {1, 1}, {0, 0}}})
	pt := geojson.NewPoint([]float64{1, 2})

	assert.Equal(t,
		bson.M{"loc": bson.M{"$geoWithin": bson.M{"$geometry": poly}}},
		Field[any]("loc").GeoWithin(poly).Filter())

	assert.Equal(t,
		bson.M{"loc": bson.M{"$geoIntersects": bson.M{"$geometry": poly}}},
		Field[any]("loc").GeoIntersects(poly).Filter())

	assert.Equal(t,
		bson.M{"loc": bson.M{"$near": bson.M{"$geometry": pt}}},
		Field[any]("loc").Near(pt).Filter())

	assert.Equal(t,
		bson.M{"loc": bson.M{"$nearSphere": bson.M{"$geometry": pt}}},
		Field[any]("loc").NearSphere(pt).Filter())
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run Test_Geospatial ./...`
Expected: FAIL — compile error, `.GeoWithin` etc. undefined.

- [ ] **Step 3: Add the marker to `geojson/geometry.go`**

Insert after the `Geometry` type's `SetBBox` method:

```go
// geo marks a type as a GeoJSON geometry. It is unexported so only types
// in this package can satisfy GeometryArg.
func (g Geometry[G]) geo() {}

// GeometryArg is any GeoJSON geometry. It lets non-generic consumers
// (e.g. field builder methods, which cannot themselves be generic in Go)
// accept any Geometry[G] with compile-time safety.
type GeometryArg interface {
	geo()
}
```

- [ ] **Step 4: Create `geospatial.go`**

```go
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
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add geojson/geometry.go geospatial.go geospatial_test.go
git commit -m "$(printf 'feat: port geospatial operators to the fluent builder\n\nAdds a GeometryArg marker to geojson so non-generic field methods\naccept any Geometry[G].\n\nCo-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>')"
```

---

### Task 5: Element operators

**Files:**
- Create: `element.go`
- Create: `element_test.go`

**Interfaces:**
- Consumes: `Op[V]`, `FieldExpr[V]`, `(FieldExpr[V]).add` (Task 1).
- Produces: `(FieldExpr[V]).Exists(bool) FieldExpr[V]`, `(FieldExpr[V]).Type(types ...string) FieldExpr[V]`.

- [ ] **Step 1: Write the failing test**

Create `element_test.go`:

```go
package mongque

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Exists(t *testing.T) {
	assert.Equal(t,
		bson.M{"name": bson.M{"$exists": true}},
		Field[string]("name").Exists(true).Filter())
}

func Test_Type_Single(t *testing.T) {
	assert.Equal(t,
		bson.M{"name": bson.M{"$type": "string"}},
		Field[string]("name").Type("string").Filter())
}

func Test_Type_Multiple(t *testing.T) {
	assert.Equal(t,
		bson.M{"name": bson.M{"$type": []string{"string", "int"}}},
		Field[string]("name").Type("string", "int").Filter())
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'Test_Exists|Test_Type' ./...`
Expected: FAIL — compile error, `.Exists` / `.Type` undefined.

- [ ] **Step 3: Create `element.go`**

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add element.go element_test.go
git commit -m "$(printf 'feat: add element operators (exists, type)\n\nCo-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>')"
```

---

### Task 6: Documentation and full verification

**Files:**
- Modify: `README.md`
- Modify: `CLAUDE.md`

**Interfaces:** none (docs only).

- [ ] **Step 1: Update the README usage example**

In `README.md`, replace the `### Example ###` code block (the `mongque.NewFilter(...)` snippet) with:

````markdown
```go
filter := mongque.NewFilter(
    mongque.Field[string]("name").Eq("John"),
    mongque.Field[int]("score").Lte(60),
)
/*
bson.M{
    "name":  bson.M{"$eq": "John"},
    "score": bson.M{"$lte": 60},
}
*/
```
````

- [ ] **Step 2: Update the README features/roadmap**

In `README.md`, replace the "Currently supports the following query types" list with exactly:

```markdown
- Comparator
- Logical
- Geospatial
- Element
```

Then replace the entire "### To Be Added ###" operators checklist with exactly:

```markdown
- [ ] Evaluation
- [ ] Array
- [ ] Bitwise
- [ ] Projection
```

- [ ] **Step 3: Update `CLAUDE.md` architecture section**

In `CLAUDE.md`, replace the `## Architecture` section body so it describes the new model. Use this text:

```markdown
The library is built around a typed, field-centric fluent builder:

- **`FieldExpr[V any]`** (`query.go`) is the builder: `Field[V](name)` starts
  it, and operator methods (`.Eq`, `.Gt`, `.In`, `.Exists`, `.GeoWithin`, `.Not`,
  ...) each return a new `FieldExpr` with one operator appended (immutable value
  semantics). `V` is the field's value type; comparison methods are compile-time
  checked against it.
- **`Op[V any]`** (`query.go`) is a single operator expression (`{key: value}`).
  Standalone constructors (`Eq`, `Gt`, `In`, ...) build `Op[V]` values used as
  building blocks, notably inside field-level `Not`.
- **`Expr`** (`query.go`) is the interface (`Filter() bson.M` / `FilterD() bson.D`)
  implemented by both `FieldExpr[V]` and the logical combinators — that shared type
  is what lets heterogeneous fields (`Field[string]`, `Field[int]`) compose inside
  `And`/`Or`/`Nor`.
- **`And`/`Or`/`Nor`** (`logical.go`) are top-level combinators over full `Expr`s,
  emitting `{$and: [...]}` etc. **`Not`** is a field-level method emitting
  `{field: {$not: {...}}}`. **`NewFilter`/`NewFilterD`** merge several `Expr`s into
  one document (implicit AND), falling back to `$and` on a field-name collision.

Operator categories live in their own files: `comparator.go`, `logical.go`,
`geospatial.go`, `element.go`. The **`geojson` subpackage** defines
`Geometry[G GeometryTypes]` and a `GeometryArg` marker interface so the
(non-generic) geospatial field methods accept any geometry type.

### Adding a new operator category

1. Add a `<category>.go` with `FieldExpr[V]` methods that call `f.add(Op[V]{"$op", value})`.
2. For operators usable inside `Not`, also add a standalone `func Op[V](...) Op[V]` constructor.
3. Add a `<category>_test.go` asserting the exact `bson.M` for each operator.
```

Also update the `### Adding a new operator category` heading in the existing CLAUDE.md if it duplicates — ensure only one such section remains.

- [ ] **Step 4: Full verification**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: build clean, no vet output, all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add README.md CLAUDE.md
git commit -m "$(printf 'docs: update README and CLAUDE.md for the fluent API\n\nCo-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>')"
```

---

## Self-Review Notes

- **Spec coverage:** Core types (Task 1); comparators (Task 2); top-level logical + field-level `Not` + `NewFilter` collision→`$and` (Task 3); geospatial + geojson marker (Task 4); element `Exists`/`Type` (Task 5); README + CLAUDE.md (Task 6). All spec sections mapped.
- **Type consistency:** `Op[V]{key, value}`, `FieldExpr[V]{name, ops}`, `add`, `Filter`/`FilterD`, `Expr`, `logicalExpr`, `GeometryArg` names are used identically across tasks.
- **Deferred:** evaluation, array, bitwise operators; geospatial legacy shapes/distance modifiers; projection/aggregation — all out of Phase 0 per the spec.
- **TDD note:** Task 1 (foundational type swap) implements-then-tests because the old API occupies the new names; Tasks 2–5 are strict red-first.
