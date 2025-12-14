package main

import (
	"testing"
)

func TestHelloMessage(t *testing.T) {
	// 簡単なテストケース
	expected := "Hello"
	actual := "Hello"

	if expected != actual {
		t.Errorf("期待値: %s, 実際: %s", expected, actual)
	}
}
