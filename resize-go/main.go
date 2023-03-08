package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nfnt/resize"

	awsService "github.com/vidhill/terraform-lambda-play/awsservice"
)

var (
	imageSize = uint(200)
	// bucket    = "vidhill-my-tf-test-bucket"
	// key       = "ainur-khasanov-WVkJxAqX1iQ-unsplash.jpg"
)

func main() {

	imageSize = setImageSize()

	// err := downloadResizeUpload(bucket, key)
	// if err != nil {
	// 	panic(err)
	// }
	lambda.Start(Handler)
}

func downloadResizeUpload(bucket, key string) error {
	newKey := "resized-" + key
	destBucket := bucket + "-resized"

	serv, err := awsService.NewService()

	if err != nil {
		return err
	}

	file, dlPath, err := serv.DownloadFileS3(bucket, key)
	if err != nil {
		return err
	}

	defer file.Close()

	r, err := streamResize(file)
	if err != nil {
		return err
	}

	if err = serv.UploadReaderContentS3(r, destBucket, newKey); err != nil {
		return err
	}

	return removeFiles(dlPath)
}

func streamResize(r io.Reader) (io.Reader, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}

	resizedImg := resizeImg(img)

	var buf bytes.Buffer

	// encode the resized image into the bytes buffer
	err = jpeg.Encode(&buf, resizedImg, nil)

	return &buf, err
}

func resizeImg(img image.Image) image.Image {
	return resize.Resize(imageSize, 0, img, resize.Bilinear)
}

func Handler(ctx context.Context, s3Event events.S3Event) error {

	for _, record := range s3Event.Records {
		srcKey := record.S3.Object.Key
		srcBucket := record.S3.Bucket.Name

		if !isJpegExtension(srcKey) {
			continue
		}

		err := downloadResizeUpload(srcBucket, srcKey)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func isJpegExtension(p string) bool {
	switch path.Ext(p) {
	case ".jpeg", ".jpg":
		return true
	default:
		return false
	}
}

func removeFiles(paths ...string) error {
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

// determine image size from environment variable
func setImageSize() uint {
	sizeString := os.Getenv("IMAGE_SIZE")
	i, err := strconv.Atoi(sizeString)
	if err != nil && i != 0 {
		return uint(i)
	}
	return imageSize
}
