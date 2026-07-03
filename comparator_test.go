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
