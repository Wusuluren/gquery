package gquery_test

import (
	"github.com/wusuluren/gquery"
	"testing"
)

func TestReStrCmp(t *testing.T) {
	testData := [][]string{
		{"abcd", "ab*"},
		{"abcdef", "ab*f"},
		{"abcdef", "*ef"},
	}
	for i, data := range testData {
		ok := gquery.ReStrCmp(data[0], data[1])
		if !ok {
			t.Fatal(i, ok)
		}
		t.Log(i, ok)
	}
}
