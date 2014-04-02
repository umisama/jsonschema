package jsonschema

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestSelector []string

var testlist = TestSelector{
	"remote ref",
	"fragment within remote ref",
	"ref within remote ref",
	"change resolution scope",
}

func (s TestSelector) IsSkip(str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type TestSuiteSchema struct {
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	Tests       []struct {
		Description string          `json:"description"`
		Data        json.RawMessage `json:"data"`
		Valid       bool            `json:"valid"`
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
	cases, err := loadTestCases(t, "./jsonSchemaTestSuite/tests/draft4")
	if err != nil {
		return
	}

	case_count, fail_count, skip_count := 0, 0, 0
	for _, v := range cases {
		if !testlist.IsSkip(v.Description) {
			case_count = case_count + v.Count()
			validator, err := newValidator(v.Schema)
			if err != nil {
				t.Error("fail on (", v.Description, ") with", err)
				continue
			}

			for ki, vi := range v.Tests {
				valid, err := validator.IsValid(vi.Data)
				if err != nil {
					fail_count = fail_count + 1
					t.Error("fail on (", v.Description, ") -", ki, "with", err)
					continue
				}

				if valid != vi.Valid {
					fail_count = fail_count + 1
					t.Error("fail on (", v.Description, ") -", ki)
					continue
				}
			}
		} else {
			skip_count = skip_count + 1
			t.Log("skipped:", v.Description)
		}
	}

	t.Log("fail count:", fail_count, "/", case_count, "cases")
	t.Log("skip count:", skip_count)

	return
}
