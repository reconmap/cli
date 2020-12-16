package api

import (
	"testing"
)

func TestRetrieveApi(t *testing.T) {
	if RetrieveData() == "" {
		t.Error("API data is empty")
	}
}
