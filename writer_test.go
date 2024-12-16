package jsonpath_writer

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"gopkg.in/go-playground/assert.v1"
)

func Test_Use_Ojg_JP(t *testing.T) {
	jsoncontent := `{
        "a":[
            {"x":1,"y":2,"z":3},
            {"x":2,"y":4,"z":6}
        ]
    }`
	// obj1 for ojg unmarshal
	obj1, err := oj.ParseString(jsoncontent)
	assert.Equal(t, err, nil)
	// obj2 for encoding/json unmarshal
	var obj2 interface{}
	decoder := json.NewDecoder(strings.NewReader(jsoncontent))
	err = decoder.Decode(&obj2)
	assert.Equal(t, err, nil)
	// obj1 and obj2 are not equal, because ojg unmarshal will convert all numbers to int64 and json unmarshal will convert all numbers to float64
	assert.NotEqual(t, obj1, obj2)
	x, err := jp.ParseString("$.a[0].y")
	assert.Equal(t, err, nil)

	expect := []interface{}{int64(2)}
	ys1 := x.Get(obj1)
	ys2 := x.Get(obj2)
	for i := range ys2 {
		ys2[i] = int64(ys2[i].(float64))
	}
	ys2[0] = ys2[0].(int64)
	// returns [4]
	assert.Equal(t, ys1, expect)
	// returns [4]
	assert.Equal(t, ys2, expect)
}

func TestParseJSONPath(t *testing.T) {
	testcases := []struct {
		jsonpath    string
		finalstep   jp.Frag
		setable     bool
		expectError error
	}{
		{
			jsonpath:  "$.a[0].y.0",
			finalstep: jp.Child('0'),
			setable:   true,
		},
		{
			jsonpath:  "$.a[0].y[3]",
			finalstep: jp.Nth(3),
			setable:   true,
		},
		{
			jsonpath:  "$",
			finalstep: jp.Root('$'),
			setable:   true,
		},
		{
			jsonpath:  "$.a[?(@.x > 1)].y",
			finalstep: jp.Child('y'),
			setable:   true,
		},
		{
			jsonpath:  "$[10].y",
			finalstep: jp.Child('y'),
			setable:   true,
		},
		{
			jsonpath:  "$.a[-10]",
			finalstep: jp.Nth(-10),
			setable:   true,
		},
		{
			jsonpath:    "",
			finalstep:   jp.Nth(-10),
			setable:     false,
			expectError: ErrorJSONPathInvalid,
		},
		{
			jsonpath:    "a[-10]",
			expectError: ErrorJSONPathInvalid,
		},
	}

	for idx, tt := range testcases {
		t.Run(fmt.Sprintf("index: %d, path: %v ", idx, tt.jsonpath), func(t *testing.T) {
			jpc, err := ParseJSONPathString(tt.jsonpath)
			if err != nil {
				assert.Equal(t, err, tt.expectError)
				return
			}
			assert.Equal(t, jpc.finalstep, tt.finalstep)
			assert.Equal(t, jpc.SetAble(), tt.setable)
		})
	}
}

func TestJSONPathSetValue(t *testing.T) {
	var testcases = []struct {
		jsonpath    string
		source      string
		value       string
		expect      string
		expectError error
	}{
		{
			jsonpath:    "$.a[0].y.0",
			source:      `{"a":[{"y":[1,2,3]}]}`,
			value:       "4",
			expect:      `{"a":[{"y":[4,2,3]}]}`,
			expectError: ErrorSourceLeafNotMap,
		},
		{
			jsonpath:    "$.a[0].y[3]",
			source:      `{"a":[{"y":[1,2,3]}]}`,
			value:       "4",
			expect:      `{"a":[{"y":[1,2,3,4]}]}`,
			expectError: ErrorSourceIndexOutOfRange,
		},
		{
			jsonpath:    "$.a[0].z[3]",
			source:      `{"a":[{"y":[1,2,3]}]}`,
			value:       "4",
			expect:      `{"a":[{"y":[1,2,3,4]}]}`,
			expectError: ErrorJSONPathNotExisted,
		},
		{
			jsonpath:    "$.a[?(@.x > 1)]",
			source:      `{"a":[{"y":[1,2,3]}]}`,
			value:       "4",
			expect:      `{"a":[{"y":[1,2,3,4]}]}`,
			expectError: ErrorJSONPathNotSetAble,
		},
		{
			jsonpath:    "$.a[0].y[0]",
			source:      `null`,
			value:       "4",
			expect:      `{"a":[{"y":[4,2,3]}]}`,
			expectError: ErrorJSONPathNotExisted,
		},
		{
			jsonpath:    "$.a[0].y[0]",
			source:      `{"a":[{"y":[1,2,3]}]}`,
			value:       "4",
			expect:      `{"a":[{"y":[4,2,3]}]}`,
			expectError: ErrorSourceLeafNotMap,
		},
		{
			jsonpath:    "$.a[0].y",
			source:      `{"a":[{"y":[1,2,3]}]}`,
			value:       `{"z":4}`,
			expect:      `{"a":[{"y":{"z":4}}]}`,
			expectError: nil,
		},
	}
	for _, tt := range testcases {
		t.Run(fmt.Sprintf("jsonpath: %v, source: %v, value: %v", tt.jsonpath, tt.source, tt.value), func(t *testing.T) {
			jpobj, err := ParseJSONPathString(tt.jsonpath)
			assert.Equal(t, err, nil)
			var source, value, except interface{}
			err = json.Unmarshal([]byte(tt.source), &source)
			assert.Equal(t, err, nil)
			err = json.Unmarshal([]byte(tt.value), &value)
			assert.Equal(t, err, nil)
			err = json.Unmarshal([]byte(tt.expect), &except)
			assert.Equal(t, err, nil)

			dest, err := jpobj.SetValue(source, value)
			if err != nil {
				fmt.Println(err)
				assert.Equal(t, err, tt.expectError)
				return
			}
			assert.Equal(t, dest, except)
		})
	}
}

func TestSourceNilSetValue(t *testing.T) {
	path := "$.a[0].y"
	dest, err := SetValue(nil, path, 1)
	assert.Equal(t, dest, nil)
	assert.Equal(t, err, ErrorJSONPathNotExisted)

}
