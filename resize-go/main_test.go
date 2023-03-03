package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsJpegExtension(t *testing.T) {

	invalidFilenames := []string{
		"foo.txt",
		"foo.png",
		"subfolder/foo.png",
	}

	for _, v := range invalidFilenames {
		t.Run(v, func(t *testing.T) {
			res := isJpegExtension(v)

			assert.False(t, res)
		})
	}

	validFilenames := []string{
		"foo.jpeg",
		"foo.jpg",
		"subfolder/foo.jpeg",
	}

	for _, v := range validFilenames {
		t.Run(v, func(t *testing.T) {
			res := isJpegExtension(v)

			assert.True(t, res)
		})
	}

}

func TestMakeTmpPath(t *testing.T) {
	p, err := resizeCopyImg("test.jpg", "resized.jpg")
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	expectedPath := makeTmpPath("resized.jpg")

	assert.Equal(t, expectedPath, p)
	assert.FileExists(t, expectedPath)

	// clean up
	os.Remove(p)
}
