package service

import (
	"bytes"
	"testing"
)

func TestDefaultData(t *testing.T) {
	_, err := ReadData(bytes.NewReader(JSONData))
	if err != nil {
		t.Error(err)
	}
}
