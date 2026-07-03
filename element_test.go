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
