package groudon

import (
	"bytes"
	"encoding/json"
	"testing"
)

type Something struct {
	Name   string `json:"name"`
	Friend string `json:"friend"`
	Age    int    `json:"age"`
}

func (some *Something) Types() (data map[string]string) {
	data = map[string]string{
		"name":   "string",
		"friend": "string",
		"age":    "number",
	}

	return
}

func (some *Something) Defaults() (data map[string]interface{}) {
	data = map[string]interface{}{
		"friend": "jim",
	}

	return
}

func Test_something(test *testing.T) {
	var data map[string]interface{} = map[string]interface{}{
		"name": "foo",
		"age":  20,
	}

	var err error
	var data_bytes []byte
	if data_bytes, err = json.Marshal(data); err != nil {
		test.Fatal(err)
	}

	var target Something
	var external error
	if err, external = SerializeBody(bytes.NewReader(data_bytes), &target); err != nil {
		test.Fatal(err)
	}

	if external != nil {
		test.Fatal(external)
	}

	if target.Name != data["name"].(string) {
		test.Errorf("data name mismatch! have: %s, want: %s", target.Name, data["name"])
	}

	var defaults map[string]interface{} = target.Defaults()
	if target.Friend != defaults["friend"].(string) {
		test.Errorf("default friend not taken! have: %s, want: %s", target.Friend, defaults["friend"])
	}

}
