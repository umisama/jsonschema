package json

import (
	"testing"
)

type TestCaseJsonType struct {
	jsontype JsonType
	ret      string
}

var TestCasesJsonType_String = []TestCaseJsonType{
	TestCaseJsonType{JsonType_Bool, "booleans"},
	TestCaseJsonType{JsonType_Nil, "null"},
	TestCaseJsonType{JsonType_Number, "numbers"},
	TestCaseJsonType{JsonType_Array, "arrays"},
	TestCaseJsonType{JsonType_Object, "objects"},
	TestCaseJsonType{JsonType_String, "strings"},
}

func Test_Test(t *testing.T) {
	for k, v := range TestCasesJsonType_String {
		if v.jsontype.String() != v.ret {
			t.Error("fail on", k)
		}
	}
	return
}
