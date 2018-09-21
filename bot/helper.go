package bot

import (
	"bytes"
	"io"
	"io/ioutil"
)

func printReader(r io.Reader) (*bytes.Buffer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
