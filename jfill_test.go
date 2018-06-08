package jfill

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFill(t *testing.T) {
	testCases := []struct {
		inJSON, in, expect string
	}{
		{`{"age":10}`, "hoge", "hoge"},
		{`{"age":10}`, "{{age}}", "10"},
		{`{"age":"10"}`, "{{age}}", "10"},
		{`{"age":"10"}`, "hoge{{age}}", "hoge10"},
		{`{"age":"10"}`, "{{age}}fuga", "10fuga"},
		{`{"age":"10"}`, "{{/age}}", "10"},
		{`{"person":{"age":10}}`, "{{/person/age}}", "10"},
		{`{}`, "{{age:15}}", "15"},
		{`{"age":15,"name":"Songmu"}`, "I'm {{name}}. {{age}} years old", "I'm Songmu. 15 years old"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %s", tc.inJSON, tc.in), func(t *testing.T) {
			var tree interface{}
			err := json.Unmarshal([]byte(tc.inJSON), &tree)
			if err != nil {
				t.Errorf("failed to unmarshal json: %s", err)
			}
			out, err := fill(tc.in, tree)
			if err != nil {
				t.Errorf("failed to fill: %s", err)
			}
			if out != tc.expect {
				t.Errorf("failed to fill\n  out=%s\nexpect=%s", out, tc.expect)
			}
		})
	}

}
