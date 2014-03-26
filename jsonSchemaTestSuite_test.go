package jsonschema

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestSuiteSchema struct {
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	Tests       []struct {
		Description string      `json:"description"`
		Data        interface{} `json:"data"`
		Valid       bool        `json:"valid"`
	} `json:"tests"`
}

func (t TestSuiteSchema) Count() int {
	return len(t.Tests)
}

func loadTestCases(t *testing.T, dir string) (testcases []TestSuiteSchema, err error) {
	dirf, err := os.Open(dir)
	if err != nil {
		t.Error("fail on load with ", err)
		return
	}

	finfos, err := dirf.Readdir(0)
	if err != nil {
		t.Error("fail on load with ", err)
		return
	}

	testcases = make([]TestSuiteSchema, 0)
	for _, v := range finfos {
		if !v.IsDir() && path.Ext(v.Name()) == ".json" {
			f, ierr := os.Open(dir + "/" + v.Name())
			if ierr != nil {
				t.Error("fail on load with ", ierr)
				err = ierr
				return
			}

			buf, ierr := ioutil.ReadAll(f)
			if ierr != nil {
				t.Error("fail on load with ", err)
				err = ierr
				return
			}

			v := make([]TestSuiteSchema, 0)
			err = json.Unmarshal(buf, &v)
			if err != nil {
				t.Error("fail on load with ", err)
				return
			}

			testcases = append(testcases, v...)
		}
	}
	return
}

func Test_jsonSchemaTestSuiteDraft4(t *testing.T) {
	t.Skip("now skip")
	cases, err := loadTestCases(t, "./test/tests/draft4")
	if err != nil {
		return
	}

	count := 0
	for _, v := range cases {
		count = count + v.Count()
	}
	t.Log("load", count, "cases")


	return
}
