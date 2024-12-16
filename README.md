# JSONPath Writer
A golang implementation of `JsonPath Writer` based on http://github.com/spyzhov/ajson 。

this libary can help  modify /update json data using JSONPath syntax.

**Golang Version Required**:  ``1.21``

## Install
```
go get github.com/mmtbak/jsonpath_writer
```

## Example

```Go
import "github.com/mmtbak/jsonpath_writer"

func main(){

	source := `
    {
        "a":1,
        "b":2
    }
    `
	jsonpath := "$.c"
	value := `
	{
		"key1": "val2",
		"key2":"va2"
	}`
	var sourcejson interface{}
	var valuejson interface{}
	json.Unmarshal([]byte(source), &sourcejson)
	json.Unmarshal([]byte(value), &valuejson)

	dest, _ := jsonpath_writer.SetValue(sourcejson, jsonpath, valuejson)
	fmt.Println(dest)
	// Output
	// {
	//     "a":1,
	//     "b":2,
	//     "c": {
	//         "key1": "val2",
	//         "key2": "va2"
	// 		}
	// }
	//
}

```
