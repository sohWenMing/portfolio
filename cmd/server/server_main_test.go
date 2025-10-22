package main

import "testing"

func TestServerRun(t *testing.T) {
	_, _, err := Run(true, "../../.env")
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
	}
}
