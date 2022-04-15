package messages

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testJSONStruct struct {
	Int      int                    `json:"int" msgpack:"int"`
	Uint     uint                   `json:"uint" msgpack:"uint"`
	String   string                 `json:"string" msgpack:"string"`
	Map      map[string]interface{} `json:"map" msgpack:"map"`
	Ignored  string                 `json:"-" msgpack:"-"`
	Ignored2 string
}

var testJSON = "{" +
	"\"int\": -420," +
	"\"uint\": 69," +
	"\"string\": \"yep\"," +
	"\"map\": {" +
	"\"key0\": \"value\"," +
	"\"key1\": 1\n  }," +
	"\"ignored\": \"pls ignore me\"" +
	"}"

func TestParseMessage(t *testing.T) {
	t.Parallel()

	var raw gin.H
	_ = json.Unmarshal([]byte(testJSON), &raw)

	var res testJSONStruct
	err := parseDynamicStruct(raw, &res)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, int(-420), res.Int)
	assert.Equal(t, uint(69), res.Uint)
	assert.Equal(t, "yep", res.String)
	assert.Equal(t, map[string]interface{}{
		"key0": "value",
		// JSON does only know float. So this is what we have to expect
		"key1": float64(1),
	}, res.Map)
	assert.Equal(t, "", res.Ignored)
	assert.Equal(t, "", res.Ignored2)
}

var testJSON_InvalidStruct = "{" +
	"\"int\": -420," +
	"\"uint\": 69," +
	"\"should_be_string_but_oh_well\": \"yep\"," +
	"\"map\": {" +
	"\"key0\": \"value\"," +
	"\"key1\": 1\n  }," +
	"\"ignored\": \"pls ignore me\"" +
	"}"

func TestParseInvalidMessage(t *testing.T) {
	t.Parallel()

	var raw gin.H
	_ = json.Unmarshal([]byte(testJSON_InvalidStruct), &raw)

	var res testJSONStruct
	err := parseDynamicStruct(raw, &res)
	if !assert.Error(t, err) {
		t.FailNow()
	}
}
