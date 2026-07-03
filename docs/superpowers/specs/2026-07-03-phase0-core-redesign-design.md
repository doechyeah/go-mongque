# Phase 0 — Core Field-Centric Redesign

**Date:** 2026-07-03
**Status:** Approved design, pending implementation plan
**Scope:** Architectural foundation for feature-completeness with the MongoDB query operator surface.

## Motivation

go-mongque today implements 12 of ~31 MongoDB query-predicate operators, but two of its three "done" categories emit **invalid MongoDB queries**:

1. **Logical operators are nested under a field name.** `And("test", …).Filter()` produces `{"test":{"$and":[…]}}`. Valid MQL is top-level: `{"$and":[…]}`. `$and/$or/$nor` operate on arrays of full query expressions, not a named field — the `name` parameter is conceptually wrong for the whole category.
2. **`$not` is modeled as a top-level operator.** In MQL `$not` inverts a *single field's* operator expression: `{age:{$not:{$gt:5}}}`. The current code shapes it like `$and`, which MongoDB rejects.
3. **No field-level operator composition.** `{age:{$gt:5,$lt:20}}` is idiomatic MQL, but `NewFilter(Gt("age",5),Lt("age",20))` writes the same `bson.M` key twice → last-wins collision that silently drops a clause.
4. **Near-zero test coverage.** Only `Test_Eq` and `Test_And` exist — which is why the above bugs plus two copy-paste operator bugs (`Gt`→`gte`, `NearSphere`→`near`, since fixed) went unnoticed.

Every operator added on the current model inherits these defects, so the architecture must be corrected first.

## Decisions (locked during brainstorming)

- **Backward compatibility:** Break freely. This is a 0.x beta; prioritize a correct, clean API. Every current call site will change.
- **API shape:** Field-centric fluent builder — a field builder accumulates operators via chaining.
- **Type safety:** Typed field, generic value — `Field[V any](name)` parameterizes the field's value type; comparison methods are compile-time checked against `V`; operators with their own argument types (`Exists`→bool, `Size`→int, `Regex`→string, geospatial→geometry) take those types directly.
- **Phase 0 contents:** Core model + port existing comparators and geospatial + element operators (`$exists`, `$type`) + full test-per-operator coverage. Array, evaluation, and bitwise categories are deferred to later phases.

## Architecture

### Core types

**`Op[V any]`** — a single field-level operator expression, e.g. `{$gt:5}`.

```go
type Op[V any] struct {
    key   string // "$gt"
    value any    // 5
}
```

`V` is a phantom type parameter: it constrains which field an operator may attach to (so `$not` and, later, `$elemMatch` stay type-safe), while `value` remains `any` so the same struct can carry element/geospatial operands whose type is unrelated to `V`.

Standalone constructors (used as building blocks, e.g. inside `Not`):

```go
func Eq[V any](v V) Op[V]      { return Op[V]{"$eq", v} }
func Ne[V any](v V) Op[V]      { return Op[V]{"$ne", v} }
func Gt[V any](v V) Op[V]      { return Op[V]{"$gt", v} }
func Gte[V any](v V) Op[V]     { return Op[V]{"$gte", v} }
func Lt[V any](v V) Op[V]      { return Op[V]{"$lt", v} }
func Lte[V any](v V) Op[V]     { return Op[V]{"$lte", v} }
func In[V any](vs ...V) Op[V]  { return Op[V]{"$in", vs} }
func Nin[V any](vs ...V) Op[V] { return Op[V]{"$nin", vs} }
```

**`FieldExpr[V any]`** — the fluent builder; accumulates operators for one field.

```go
type FieldExpr[V any] struct {
    name string
    ops  []Op[V]
}

func Field[V any](name string) FieldExpr[V] {
    return FieldExpr[V]{name: name}
}
```

Methods return a **new copy** with one operator appended (immutable chaining, mirroring the current `SetName`/`SetValue` value-receiver style). Each `add` allocates a fresh `ops` slice so chained/branched builders never share backing arrays.

Method set:

- **Comparison** (checked against `V`): `Eq(V) Ne(V) Gt(V) Gte(V) Lt(V) Lte(V) In(...V) Nin(...V)` — each delegates to the standalone `Op[V]` constructor.
- **Element:** `Exists(bool)`, `Type(...string)`.
- **Geospatial:** `GeoWithin`, `GeoIntersects`, `Near`, `NearSphere` (see Geospatial section).
- **Field-level logical:** `Not(Op[V])` — wraps as `{$not:{<key>:<value>}}`; type-checked against `V` so `Field[string]("name").Not(Gt(5))` will not compile. Phase 0 accepts only comparison `Op[V]` here; `$not` over element operators is deferred.

**`Expr`** — the common interface that lets heterogeneous fields compose.

```go
type Expr interface {
    Filter() bson.M
    FilterD() bson.D
}
```

`FieldExpr[V]` and the logical result type both implement `Expr`. A generic type satisfying a non-generic interface is precisely what allows `Field[string]("status")` and `Field[int]("score")` to sit inside the same `Or(...)`.

### Rendering

`FieldExpr[V].Filter()` merges all accumulated operators into one document — the fix for the silent key-collision bug:

```go
Field[int]("age").Gte(18).Lt(65).Filter()
// bson.M{"age": bson.M{"$gte": 18, "$lt": 65}}
```

`FilterD()` renders the same content in insertion order (`bson.D`), giving deterministic output for order-sensitive operators and for tests. A single-operator field still renders explicitly (`{name:{$eq:"John"}}`) to match existing behavior.

### Top-level logical combinators

Take full expressions, emit at the top level (fixing the nested-under-a-name bug):

```go
func And(exprs ...Expr) Expr
func Or(exprs ...Expr) Expr
func Nor(exprs ...Expr) Expr
```

Backed by an internal `logicalExpr{ op string; exprs []Expr }` implementing `Expr`:

```go
Or(Field[string]("status").Eq("active"), Field[int]("score").Gt(90)).Filter()
// bson.M{"$or": bson.A{
//     bson.M{"status": bson.M{"$eq": "active"}},
//     bson.M{"score":  bson.M{"$gt": 90}},
// }}
```

Array elements are rendered via each child's `Filter()` (or `FilterD()` for the `FilterD` path) into a `bson.A`.

### Implicit-AND helpers

`NewFilter(...Expr) bson.M` remains the ergonomic implicit AND over **distinct** fields, merging each expression's document into one:

```go
NewFilter(Field[string]("name").Eq("John"), Field[int]("score").Lte(60))
// bson.M{"name": bson.M{"$eq":"John"}, "score": bson.M{"$lte":60}}
```

On a duplicate top-level key it transparently falls back to `And(exprs...).Filter()` rather than dropping a clause. `NewFilterD(...Expr) bson.D` returns the ordered form.

### Geospatial & geojson

Go forbids generic methods, so the current `GeoWithinGeoJSON[G]` cannot become a fluent method. Resolution: add an unexported marker method to the `geojson` package so a single non-generic method accepts any geometry type with compile-time safety.

```go
// geojson/geometry.go
func (g Geometry[G]) geo() {}         // implemented by every Geometry[G]
type GeometryArg interface { geo() }  // only geojson geometries satisfy it
```

Geospatial methods on `FieldExpr[V]` then take `geojson.GeometryArg` and wrap in `$geometry` (matching current behavior):

```go
Field[any]("loc").GeoWithin(geojson.NewPolygon(coords)).Filter()
// bson.M{"loc": bson.M{"$geoWithin": bson.M{"$geometry": {…}}}}
```

This collapses the current plain + `…GeoJSON` constructor pairs into one method per operator. Legacy shape specifiers (`$box`, `$center`, `$centerSphere`, `$polygon`) and distance modifiers (`$maxDistance`, `$minDistance`) remain deferred to Phase 4.

## Files

**Rewrite:**
- `query.go` — `Op[V]`, `FieldExpr[V]`, `Expr`, `logicalExpr`, `NewFilter`, `NewFilterD`, `Field`.
- `comparator.go` — standalone `Op[V]` comparison constructors + corresponding `FieldExpr` methods.
- `logical.go` — top-level `And`/`Or`/`Nor` + field-level `Not`.
- `geospatial.go` — `GeoWithin`/`GeoIntersects`/`Near`/`NearSphere` as `FieldExpr` methods.

**Add:**
- `element.go` — `Exists`/`Type` methods.

**Edit:**
- `geojson/geometry.go` — add `geo()` marker + `GeometryArg` interface.
- `query_test.go` — rewrite for the new API (and split into per-file test files as coverage grows).
- `README.md` — update the usage example to the fluent API.

**Delete:**
- `Field[Operable]`, `SetName`/`SetValue`, the `Operable`/`comparator`/`logical`/`geospatial` string-type constraint machinery, and all `*Bson` / `*GeoJSON` constructor variants.

## Testing

Table-driven tests asserting the exact `bson.M` produced, one case per operator, plus:

- Multi-operator composition on one field (`Gte(18).Lt(65)`).
- Heterogeneous top-level `Or`/`And`/`Nor` (mixed `Field[string]`/`Field[int]`).
- Field-level `Not` shape (`{field:{$not:{$gt:5}}}`).
- `NewFilter` collision → `$and` fallback.
- `FilterD` insertion ordering.
- Element operators (`Exists`, `Type`) and geospatial `$geometry` wrapping.

## Out of scope (later phases)

- **Phase 1:** Evaluation operators (`$expr`, `$jsonSchema`, `$mod`, `$regex`, `$text`, `$where`).
- **Phase 2:** Array operators (`$all`, `$elemMatch`, `$size`).
- **Phase 3:** Bitwise operators (`$bitsAllClear/Set`, `$bitsAnyClear/Set`).
- **Phase 4:** Geospatial legacy shapes and distance modifiers; geojson JSON marshaller and `GeometryCollection`.
- Projection operators and aggregation-pipeline construction (separate builder surface / module).
- `$comment`, `$rand` — explicitly excluded per README.
