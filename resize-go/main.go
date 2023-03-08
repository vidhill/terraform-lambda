package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nfnt/resize"

	awsService "github.com/vidhill/terraform-lambda-play/awsservice"
	"github.com/vidhill/terraform-lambda-play/util"
)

var (
	IMAGE_SIZE  = uint(200)
	DEST_BUCKET string
	// bucket    = "vidhill-my-tf-test-bucket"
	// key       = "ainur-khasanov-WVkJxAqX1iQ-unsplash.jpg"
)

type FilesProvider interface {
	LoadFile(bucket, key string) (io.ReadCloser, string, error)
	WriteReaderContent(r io.Reader, bucket, key string) error
	WriteFile(filePath, bucket, key string) error
	RemoveFiles(paths ...string) error
}

func init() {
	// load values from env vars
	IMAGE_SIZE = getImageSize()

	if b, err := getDestinationBucket(); err != nil {
		panic(err)
	} else {
		DEST_BUCKET = b
	}

}

func main() {

	// err := downloadResizeUpload(bucket, key)
	// if err != nil {
	// 	panic(err)
	// }
	lambda.Start(Handler)
}

func downloadResizeUpload(bucket, key string, serv FilesProvider) error {

	file, dlPath, err := serv.LoadFile(bucket, key)
	if err != nil {
		return err
	}

	defer file.Close()

	r, err := streamResize(file)
	if err != nil {
		return err
	}

	if err = serv.WriteReaderContent(r, DEST_BUCKET, key); err != nil {
		return err
	}

	return serv.RemoveFiles(dlPath)
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
	return resize.Resize(IMAGE_SIZE, 0, img, resize.Bilinear)
}

func Handler(ctx context.Context, s3Event events.S3Event) error {

	serv, err := awsService.NewService()

	if err != nil {
		return err
	}

	for _, record := range s3Event.Records {
		srcKey := record.S3.Object.Key
		srcBucket := record.S3.Bucket.Name

		if !util.IsJpegExtension(srcKey) {
			continue
		}

		err := downloadResizeUpload(srcBucket, srcKey, &serv)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

// determine image size from environment variable
func getImageSize() uint {
	sizeString := os.Getenv("IMAGE_SIZE")
	i, err := strconv.Atoi(sizeString)
	if err != nil && i != 0 {
		return uint(i)
	}
	return IMAGE_SIZE
}

func getDestinationBucket() (string, error) {
	destBucket := os.Getenv("DESTINATION_BUCKET")

	if destBucket == "" {
		return "", fmt.Errorf("missing required environment varuable")
	}

	return destBucket, nil
}
