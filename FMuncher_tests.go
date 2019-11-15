package fmuncher

import (
	"fmt"
	"os"
)

// Tests splitting of file
//
// Takes file and print splits of offset and length
func TestSplit(source string) error {
	source_file, openErr := os.Open(source)

	if openErr != nil {
		return openErr
	}

	muncher := Munch(source_file)

	splits, splitErr := muncher.Split()

	if splitErr != nil {
		return splitErr
	}

	for index, split := range splits {
		fmt.Printf("%d. Offset: %d, Length: %d\n", index+1, split.Off, split.Len)
	}

	return nil
}

// Test for proper splitting and gathering back of file. At the end, the file is copied
//
// This takes the source file path and the destination path
func TestScatterGather(source, dest string) error {
	sourceOpen, sourceErr := os.Open(source)
	destOpen, destErr := os.OpenFile(dest, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if sourceErr != nil {
		return sourceErr
	}

	if destErr != nil {
		return destErr
	}

	muncher := Munch(sourceOpen)
	splits, splitErr := muncher.Split()

	if splitErr != nil {
		return splitErr
	}

	fmt.Println("Splitted:")
	for _, split := range splits {
		fmt.Printf("Offset: %d\nLength: %d\n\n", split.Off, split.Len)
	}

	fmt.Printf("Gathering back into %s\n", dest)

	for _, split := range splits {
		buf := make([]byte, split.Len)
		_, readErr := sourceOpen.ReadAt(buf, split.Off)

		if readErr != nil {
			return readErr
		}

		fmt.Printf("Writing %d bytes at offset %d into %s\n", split.Len, split.Off, dest)

		_, writeErr := destOpen.WriteAt(buf, split.Off)
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}
