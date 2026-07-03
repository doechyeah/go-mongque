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

The whole library is built around one generic type and one constraint:

- **`Operable`** (`query.go`) is a type-set constraint: `comparator | logical | geospatial`. Each operand type is a `string` (e.g. `comparator = "$eq"`) that implements `set(val any) bson.M` — that method is what wraps a value into its Mongo operator document.
- **`Field[T Operable]`** (`query.go`) is the universal builder: `{name, op, value}`. It renders to a filter via `Filter()` (→ `bson.M`), `FilterD()` (→ `bson.D`), or `FilterE()` (→ `bson.E`, only valid inside a `bson.D`).
- **`NewFilter` / `NewFilterD`** combine multiple `Field`s of the same operand type `T` into one filter.

Each operator category lives in its own file and follows the same pattern — a `string` type with a `set` method, a block of `const` operator strings, and exported constructor functions returning `Field[T]`:

- `comparator.go` — `$eq`, `$ne`, `$lte`, `$lt`, `$gte`, `$gt`, `$in`, `$nin` (`Eq`, `Neq`, `Lte`, ...).
- `logical.go` — `$and`, `$not`, `$nor`, `$or`. Each has two constructors: a variadic `And[T](name, ...Field[T])` that flattens sub-fields via `FilterE()`, and an `*Bson` variant (`AndBson`) taking a raw `bson.D`.
- `geospatial.go` — `$geoIntersects`, `$geoWithin`, `$near`, `$nearSphere`. Its `set` wraps values as `{"$op": {"$geometry": v}}`. Each op has a plain (`GeoWithin`) and GeoJSON-typed (`GeoWithinGeoJSON[G]`) constructor.

The **`geojson` subpackage** (`geojson/geometry.go`) defines `Geometry[G GeometryTypes]` where `GeometryTypes` is a constraint over `Point | MultiPoint | LineString | MultiLineString | Polygon | MultiPolygon` (each a nested `[]float64` slice type). Constructors like `NewPoint`, `NewPolygon` build these; they carry `bson`/`json` tags for encoding.

### Adding a new operator category

1. Create `<category>.go` with a `type <category> string` and a `func (c <category>) set(v any) bson.M`.
2. Add the type to the `Operable` union in `query.go`.
3. Add `const` operator strings and exported constructor functions returning `Field[<category>]`.

## Conventions

- Tests use `testify/assert` and verify the exact `bson.M`/`bson.D` shape produced by `.Filter()`. Reference operator constants (e.g. `string(eq)`) rather than hardcoding strings.
- Operator string constants are unexported; only the constructor functions are part of the public API.
