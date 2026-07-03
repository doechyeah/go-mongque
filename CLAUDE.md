# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

go-mongque (pronounced "mong-key") is a Go library that generates MongoDB query filters (`bson.M` / `bson.D`) using generics. It is a small, dependency-light utility mirroring the [MongoDB query operators](https://www.mongodb.com/docs/manual/reference/operator/query/). The module is `github.com/doechyeah/go-mongque` (package `mongque`), with a `geojson` subpackage.

## Commands

```sh
go test ./...              # run all tests (root + geojson)
go test -run Test_Eq       # run a single test by name
go test -v ./...           # verbose
go build ./...             # compile
go vet ./...               # static checks
```

Requires Go 1.19+ (uses generics).

## Architecture

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

## Conventions

- Tests use `testify/assert` and verify the exact `bson.M`/`bson.D` shape produced by `.Filter()`. Reference operator constants (e.g. `string(eq)`) rather than hardcoding strings.
- Operator string constants are unexported; only the constructor functions are part of the public API.
