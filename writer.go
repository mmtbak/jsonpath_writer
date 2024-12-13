package jsonpath_writer

import (
	"github.com/ohler55/ojg/jp"
)

var (
	RootFrag = jp.Root('$')
)

func ParseJSONPathString(path string) (*JSONPathCompiled, error) {
	jpexpr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	finalstep := jpexpr[len(jpexpr)-1]

	jpc := &JSONPathCompiled{
		compiled:  jpexpr,
		finalstep: finalstep,
	}
	return jpc, nil
}

// JSONPathCompiled is a compiled JSONPath expression.
type JSONPathCompiled struct {
	// compiled is the compiled JSONPath expression.
	compiled jp.Expr
	// finalstep is the final step in the expression.
	finalstep jp.Frag
}

func (jpc JSONPathCompiled) SetAble() bool {

	switch jpc.finalstep.(type) {
	// only support these types can be set, other types only support to get and read data
	case jp.Child, jp.Nth, jp.Root:
		return true
	default:
		return false
	}
}

// SetValue sets the value object in the source object by the compiled JSONPath expression.
func (jpc JSONPathCompiled) SetValue(source interface{}, value interface{}) (interface{}, error) {

	// var jsonMap map[string]interface{}
	// var ok bool

	// var curobj = obj
	// var len = len(rc.compiled.steps)
	// // 正向操作有不可预测的风险，反向合并就好了

	// for idx, s := range rc.compiled.steps {
	// 	// "key", "idx"
	// 	switch s.op {
	// 	case "key":
	// 		if jsonMap, ok = curobj.(map[string]interface{}); !ok {
	// 			return nil, fmt.Errorf("data error :  key '%s' data should be map", s.key)
	// 		}
	// 		// 最后一个key
	// 		if idx+1 == len {
	// 			jsonMap[s.key] = value
	// 		}
	// 		// 不是最后一个key
	// 		val, exists := jsonMap[s.key]
	// 		if !exists {
	// 			// 如果不存在， 创建新值
	// 			jsonMap[s.key] = map[string]interface{}{}
	// 			curobj = jsonMap[s.key]
	// 			// return nil, fmt.Errorf("key error: '%s' not found in object", s.key)
	// 		} else {
	// 			curobj = val
	// 		}

	// 	}
	// }

	if !jpc.SetAble() {
		return nil, ErrorJSONPathNotSetAble
	}

	if jpc.finalstep == RootFrag {
		return value, nil
	}

	prefixsteps := jpc.compiled[0 : len(jpc.compiled)-1]
	prefixdata := prefixsteps.Get(source)
	if len(prefixdata) == 0 {
		return nil, ErrorJSONPathNotExisted
	}
	for _, data := range prefixdata {

		switch jpc.finalstep.(type) {
		case jp.Child:
			datamap, ok := data.(map[string]interface{})
			if !ok {
				return nil, ErrorSourceLeafNotMap
			}
			key := string(jpc.finalstep.(jp.Child))
			datamap[key] = value
		case jp.Nth:
			dataarray, ok := data.([]interface{})
			if !ok {
				return nil, ErrorSourceLeafNotArray
			}
			idx := int(jpc.finalstep.(jp.Nth))
			var realidx = 0
			// support negative index
			if idx < 0 {
				realidx = len(dataarray) + idx
				if realidx < 0 {
					return nil, ErrorSourceIndexOutOfRange
				}
			} else {
				// support positive index
				realidx = idx
				if realidx >= len(dataarray) {
					return nil, ErrorSourceIndexOutOfRange
				}
			}
			dataarray[realidx] = value
		}
	}

	return source, nil

}

func SetValue(source interface{}, path string, value interface{}) (data interface{}, err error) {
	jpc, err := ParseJSONPathString(path)
	if err != nil {
		return nil, err
	}
	result, err := jpc.SetValue(source, value)
	if err != nil {
		return nil, err
	}
	return result, nil
}
