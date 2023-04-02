package queryop

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Field[T Operable] struct {
	Name  string
	op    T
	Value any
}

type Operable interface {
	set(val any) bson.M
	comparator
}

func (q Field[Operable]) Filter() bson.M {
	return bson.M{
		q.Name: q.op.set(q.Value),
	}
}
