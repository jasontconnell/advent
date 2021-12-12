package common

import (
	"io"
	"log"
	"os"
)

func TeeOutput(w io.Writer) io.Writer {
	f, err := os.OpenFile("result.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(f, w)
	return mw
}
