package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vidhill/terraform-lambda-play/util"
)

func TestIsJpegExtension(t *testing.T) {

	invalidFilenames := []string{
		"foo.txt",
		"foo.png",
		"subfolder/foo.png",
	}

	for _, v := range invalidFilenames {
		t.Run(v, func(t *testing.T) {
			res := util.IsJpegExtension(v)

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
			res := util.IsJpegExtension(v)

			assert.True(t, res)
		})
	}

}
