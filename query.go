package queryop

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Operable interface {
	set(val any) bson.M
	comparator | logical
}

type Field[T Operable] struct {
	name  string
	op    T
	value any
}

func (f Field[Operable]) Filter() bson.M {
	return bson.M{
		f.name: f.op.set(f.value),
	}
}

func (f Field[Operable]) FilterD() bson.D {
	return bson.D{
		f.FilterE(),
	}
}

func (f Field[Operable]) FilterE() bson.E {
	return bson.E{
		Key: f.name, Value: f.op.set(f.value),
	}
}

func (f Field[Operable]) SetName(name string) Field[Operable] {
	f.name = name
	return f
}

func (f Field[Operable]) SetField(value string) Field[Operable] {
	f.value = value
	return f
}
