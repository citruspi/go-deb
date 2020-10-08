package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"strings"

	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

func normalizeArchivedFileName(s string) string {
	s = strings.TrimPrefix(s, "/")

	return strings.TrimSuffix(s, "/")
}

// Returns contents of 'control' file within 'control.tar.(gz|xz)'
// archive in a deb package
func ReadControlDataBytes(reader io.Reader) ([]byte, int64, error) {
	var found bool
	var num int64
	var err error
	var ctrBuf bytes.Buffer
	var debBuf bytes.Buffer
	var compressionMode string

	debReader := ar.NewReader(reader)

	for {
		header, err := debReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, num, err
		}

		name := normalizeArchivedFileName(header.Name)

		if strings.HasPrefix(name, "control.") {
			if strings.HasSuffix(name, ".tar.gz") {
				compressionMode = "tar-gz"
			} else if strings.HasSuffix(name, ".tar.xz") {
				compressionMode = "tar-xz"
			} else {
				return nil, num, errors.New("unsupported control archive compression")
			}

			num, err = io.Copy(&debBuf, debReader)
		}
	}

	if num == 0 {
		return nil, num, errors.New("no control archive present")
	}

	num = 0

	var cmpReader io.Reader

	switch compressionMode {
	case "tar-gz":
		cmpReader, err = gzip.NewReader(bytes.NewReader(debBuf.Bytes()))
	case "tar-xz":
		cmpReader, err = xz.NewReader(bytes.NewReader(debBuf.Bytes()))
	}

	if err != nil {
		return nil, 0, err
	}

	tarReader := tar.NewReader(cmpReader)

	for {
		hdr, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, err
		}

		if hdr.Name == "./control" {
			num, err = io.Copy(&ctrBuf, tarReader)
			found = true
		}
	}

	if !found {
		return nil, num, errors.New("no control data present")
	}

	return ctrBuf.Bytes(), num, err
}