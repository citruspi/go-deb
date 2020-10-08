package deb

import (
	"os"
	"testing"
)

var (
	pkg = os.Getenv("TEST_PKG")
)

func TestReadControlDataBytes(t *testing.T) {
	if pkg == "" {
		t.Fatalf("TEST_PKG undefined; set to path of package")
	}

	f, err := os.Open(pkg)

	if err != nil {
		t.Fatalf(err.Error())
	}

	_, _, err = ReadControlDataBytes(f)

	if err != nil {
		t.Fatalf(err.Error())
	}
}