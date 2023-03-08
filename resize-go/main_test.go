package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	localFile "github.com/vidhill/terraform-lambda-play/localfile"
)

func TestDownloadResizeUpload(t *testing.T) {
	outFile := "out.jpeg"
	localServ := localFile.NewLocalServ(outFile)

	assert := assert.New(t)

	if err := downloadResizeUpload("dummyBucket", "testdata/test.jpg", localServ); err != nil {
		assert.FailNow(err.Error())
		return
	}

	info, err := os.Stat(outFile)

	// check file exists
	if err != nil {
		assert.FailNow(err.Error())
		return
	}

	// ensure is not an empty file
	assert.NotEqual(info.Size(), int64(0))

	// clean up
	err = localServ.CleanUp()

	assert.NoError(err)
}
