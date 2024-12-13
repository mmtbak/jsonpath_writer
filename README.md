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
    add_data :=
    sourcejson, _ := json.Umashall([]byte(source))
    dest, _ := jsonpath_writer.Set(sourcejson, jsonpath,  add_data)
    fmt.Println(dest.String)
    // Output
    // {
    //     "a":1,
    //     "b":2,
    //     "c":3
    // }
    //
}


```
