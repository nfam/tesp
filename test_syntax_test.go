package tesp

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"strings"
	"testing"
)

// [syntax, location?]

//go:embed test_syntax.json
var test_syntax_json []byte

func TestSyntax(t *testing.T) {
	var items [][]json.RawMessage
	err := json.Unmarshal(test_syntax_json, &items)
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range items {
		var e Expression
		var location *string

		if len(item) > 1 {
			var value string
			if err := json.Unmarshal(item[1], &value); err != nil {
				t.Error(err)
				continue
			}
			location = &value
		}
		if err := json.Unmarshal(item[0], &e); err != nil {
			if location != nil {
				if !strings.HasSuffix(err.Error(), *location) {
					t.Error("error location does not match.\nMessage: ", err, "\nExpected: ", *location)
				}
			} else {
				t.Error(err)
			}
			continue
		}
		if location != nil {
			t.Error("should fail with location: " + *location + "\n" + string(item[0]))
			continue
		}

		data, err := json.Marshal(&e)
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(data, item[0]) {
			t.Error("should sucessful stringify.\nExpected:\n" + string(item[0]) + "\n\nActual:\n" + string(data))
		}
	}
}
