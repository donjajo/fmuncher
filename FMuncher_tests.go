package fmuncher

import (
	"fmt"
	"os"
)

func TestSplit(source, dest string) error {
	source_file, openErr := os.Open(source)

	if openErr != nil {
		return openErr
	}

	fmuncher := Munch(source_file)

	fmt.Println(fmuncher.Split())

	return nil
}
