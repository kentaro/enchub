package main

import (
	"testing"
)

func TestReplaceCharset(t *testing.T) {
	{
		meta := `<meta charset="utf-8">`
		actual := replaceCharset(meta, "EUC-JP")
		expected := `<meta charset="EUC-JP">`

		if actual != expected {
			t.Errorf("\nGot: %v\nExpected: %v", actual, expected)
		}
	}

	{
		meta := `<meta charset='utf-8'>`
		actual := replaceCharset(meta, "EUC-JP")
		expected := `<meta charset='EUC-JP'>`

		if actual != expected {
			t.Errorf("\nGot: %v\nExpected: %v", actual, expected)
		}
	}

	{
		meta := `<form accept-charset="utf-8" action="/" method="get">`
		actual := replaceCharset(meta, "EUC-JP")
		expected := `<form accept-charset="EUC-JP" action="/" method="get">`

		if actual != expected {
			t.Errorf("\nGot: %v\nExpected: %v", actual, expected)
		}
	}

	{
		meta := `<form accept-charset='utf-8' action="/" method="get">`
		actual := replaceCharset(meta, "EUC-JP")
		expected := `<form accept-charset='EUC-JP' action="/" method="get">`

		if actual != expected {
			t.Errorf("\nGot: %v\nExpected: %v", actual, expected)
		}
	}
}
