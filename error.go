package jsonpath_writer

import "errors"

var (
	ErrorJSONPathNotSetAble = errors.New("jsonpath is not set able")
	ErrorJSONPathNotExisted = errors.New("jsonpath not found in source")
	// ErrorSourceLeafNotMap is returned when  final step of jsonpath in  source is not a map
	ErrorSourceLeafNotMap = errors.New("source leaf data should be map")
	// ErrorSourceLeafNotArray is returned when  final step of jsonpath in  source is not an array
	ErrorSourceLeafNotArray    = errors.New("source leaf data should be array")
	ErrorSourceIndexOutOfRange = errors.New("source data index out of range")
)
