# go-mongque #

(Pronounced mong-key) Simple utility to generate MongoDB Filter for queries using generics.

## Usage ##

Import

```go
import "github.com/doechyeah/go-mongque"
```

Install

```sh
go get github.com/doechyeah/go-mongque
```

### Example ###

```go
type Document struct {
    Id      primitive.ObjectID `bson:"_id"`
    Name    string             `bson:"name"`
    Summary string             `bson:"summary"`
    Score   int                `bson:"score"`
}

client, _ := mongo.Connect(context.TODO(), clientOptions)

coll := client.Database("test").Collection("sample")

filter := mongque.NewFilter( mongoque.Eq("name", "John"), mongque.Lte("score", 60) )
/**
bson.M {
    "summary" : bson.M{"$eq": "sample"},
    "score": bson.M{"$lte": 60},
}
*/

var Doc Document
_ = coll.FindOne(context.TODO(), filter).Decode(&doc)
```

## Features ##

Initial version: `0.1.0`

Operators are based on the MongoDB query operations <https://www.mongodb.com/docs/manual/reference/operator/query/>

Currently supports the following query types:

- Comparator
- Logical

Additional operators are planned to be developed alongside support for building aggregation pipelines.

### To Be Added ###

Priority is listed in order.

Operators:

- [ ] Geospatial
- [ ] Projection
- [ ] Array
- [ ] Element
- [ ] Evaluation
- [ ] Bitwise

Miscellaneous operators such as `$comment` and `$rand` will not be added.
