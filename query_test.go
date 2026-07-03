package mongque

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Field_Empty(t *testing.T) {
	assert.Equal(t, bson.M{"name": bson.M{}}, Field[string]("name").Filter())
}

func Test_Field_EmptyD(t *testing.T) {
	assert.Equal(t, bson.D{{Key: "name", Value: bson.D{}}}, Field[string]("name").FilterD())
}
