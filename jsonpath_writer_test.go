package jsonpath_writer

import (
	"encoding/json"
	"testing"

	"github.com/yalp/jsonpath"
)

var jsondata = `
{
    "add": {
        "a": 1,
        "b": 3
    },
    "max": {
        "a": 1,
        "b": 2
    }
}

`

func TestJsonpath(t *testing.T) {
	var err error
	var bookstore interface{}
	err = json.Unmarshal([]byte(jsondata), &bookstore)
	if err != nil {
		t.Error(err)
		return
	}

	authors, err := jsonpath.Read(bookstore, "$.max.a")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(authors)

	jp, err := jsonpath.Prepare("$.max.a")
	if err != nil {
		t.Error(err)
		return
	}
	authors, err = jp(bookstore)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(authors)

}
