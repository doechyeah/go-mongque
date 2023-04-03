package mongque

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Eq(t *testing.T) {
	f := Eq("test", 123).Filter()
	assert.Equal(t, bson.M{"test": bson.M{string(eq): 123}}, f)
}

func Test_And(t *testing.T) {
	f := And("test", Eq("test", 123), Neq("test", 321)).Filter()
	if andFilter, ok := f["test"]; ok {
		assert.EqualValues(t, bson.M{string(and): bson.D{
			{Key: "test", Value: bson.M{string(eq): 123}},
			{Key: "test", Value: bson.M{string(neq): 321}},
		}}, andFilter)
	} else {
		t.Logf("could not find andFilter in filter: %v", f)
		t.Fail()
	}
}
