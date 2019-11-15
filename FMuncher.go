package fmuncher

/*
#include <fcntl.h>
#if defined( __GNUC__ )
	int isGNU = 1;
#else
	int isGNU = 0;
#endif
*/
import "C"
import (
	"errors"
	"math"
	"log"
	"os"
	"syscall"
)

// Get block size of device
const BLKSIZE int64 = (128 * 1024) * 1024

type munch struct {
	sourceFile *os.File
	maxSplit   int64
	FileInfo   syscall.Stat_t
}

type Splits struct {
	Off int64
	Len int64
}

// FMuncher takes file pointer and manipulate as you like
// Accepts *os.File of file opened for reading
//	muncher := fmuncher.Munch( filePointer )
func Munch(file *os.File) *munch {
	fmuncher := &munch{}

	fmuncher.sourceFile = file
	FileInfo := syscall.Stat_t{}
	err := syscall.Fstat(int(file.Fd()), &FileInfo )
	if err != nil {
		log.Print(err)
	}
	fmuncher.FileInfo = FileInfo
	
	// Calculate number of splits based on the blksize
	fmuncher.maxSplit = int64(math.Ceil(float64(FileInfo.Size) / float64(BLKSIZE)))
	return fmuncher
}

// Splits a file into parts of Splits{}
// Returns array of Splits{} and error if unable to seek file position
//	muncher := fmuncher.Munch( filePointer )
//	splits, err := muncher.Split()
func (fmuncher munch) Split() ([]Splits, error) {
	size := fmuncher.FileInfo.Size
	var splits []Splits

	// Size of file is less than or same as block size, return only one part and the size.
	if size <= BLKSIZE {
		return append(splits, Splits{0, size}), nil
	}

	// Is this GNU? Tell the kernel this file will be accessed randomly
	if C.isGNU == 1 {
		C.posix_fadvise(C.int(fmuncher.sourceFile.Fd()), 0, 0, C.POSIX_FADV_RANDOM)
	}

	var i int64 = 0
	var lastOffset int64
	var seekErr error

	// Loop through number of splits
	for ; i < fmuncher.maxSplit; i++ {
		var split = Splits{}
		blksize := BLKSIZE

		// First split, set offset to 0 and end in block size
		if i == 0 {
			split.Off = 0
			split.Len = BLKSIZE

			splits = append(splits, split)
			lastOffset, seekErr = fmuncher.sourceFile.Seek(BLKSIZE, 0)
			if seekErr != nil {
				return splits, errors.New("Seek error")
			}
			continue
		}

		split.Off = lastOffset

		// last offset + block size is greater than the size? Lets get the remaining bytes instead. Leave no byte behind, don't add more
		if lastOffset+BLKSIZE >= size {
			blksize = size - lastOffset
		}

		lastOffset, seekErr = fmuncher.sourceFile.Seek(blksize, 1)
		if seekErr != nil {
			return splits, errors.New("Seek error")
		}

		split.Len = blksize
		splits = append(splits, split)
	}

	return splits, nil
}
