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

func Test_NewFilterD_Merge(t *testing.T) {
	got := NewFilterD(Field[string]("name").Eq("John"), Field[int]("score").Lte(60))
	assert.Equal(t, bson.D{
		{Key: "name", Value: bson.D{{Key: "$eq", Value: "John"}}},
		{Key: "score", Value: bson.D{{Key: "$lte", Value: 60}}},
	}, got)
}

func Test_NewFilterD_CollisionFallsBackToAnd(t *testing.T) {
	got := NewFilterD(Field[int]("age").Gt(5), Field[int]("age").Lt(20))
	assert.Equal(t, bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 5}}}},
		bson.D{{Key: "age", Value: bson.D{{Key: "$lt", Value: 20}}}},
	}}}, got)
}

func Test_LogicalEmpty(t *testing.T) {
	assert.Equal(t, bson.M{}, And().Filter())
	assert.Equal(t, bson.M{}, Or().Filter())
	assert.Equal(t, bson.M{}, Nor().Filter())
	assert.Equal(t, bson.D{}, And().FilterD())
	assert.Equal(t, bson.D{}, Or().FilterD())
	assert.Equal(t, bson.D{}, Nor().FilterD())
}

func Test_Not_FilterD(t *testing.T) {
	got := Field[int]("age").Not(Gt(5)).FilterD()
	assert.Equal(t, bson.D{{Key: "age", Value: bson.D{
		{Key: "$not", Value: bson.M{"$gt": 5}},
	}}}, got)
}

func Test_Not_StandaloneIn(t *testing.T) {
	got := Field[int]("x").Not(In(1, 2)).Filter()
	assert.Equal(t, bson.M{"x": bson.M{"$not": bson.M{"$in": []int{1, 2}}}}, got)
}

func Test_And_Heterogeneous_FilterD(t *testing.T) {
	got := And(Field[string]("s").Eq("a"), Field[int]("n").Gt(1)).FilterD()
	assert.Equal(t, bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "s", Value: bson.D{{Key: "$eq", Value: "a"}}}},
		bson.D{{Key: "n", Value: bson.D{{Key: "$gt", Value: 1}}}},
	}}}, got)
}
