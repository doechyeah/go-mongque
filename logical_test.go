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
