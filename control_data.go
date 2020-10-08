package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"

	"github.com/blakesmith/ar"
)

// Returns contents of 'control' file within 'control_data.tar.gz'
// archive in a deb package
func ReadControlDataBytes(reader io.Reader) ([]byte, int64, error) {
	var found bool
	var num int64
	var err error
	var ctrBuf bytes.Buffer
	var debBuf bytes.Buffer

	debReader := ar.NewReader(reader)

	for {
		header, err := debReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, num, err
		}

		if header.Name == "control.tar.gz/" {
			num, err = io.Copy(&debBuf, debReader)
		}
	}

	if num == 0 {
		return nil, num, errors.New("no control archive present")
	}

	num = 0

	gzipReader, err := gzip.NewReader(bytes.NewReader(debBuf.Bytes()))

	if err != nil {
		return nil, 0, err
	}

	tarReader := tar.NewReader(gzipReader)

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