package main

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/nfnt/resize"
)

// const (
// 	bucket = "vidhill-my-tf-test-bucket"
// 	key    = "3375475.jpeg"
// )

type awsService struct {
	session *session.Session
}

func main() {

	// err := downloadResizeUpload(bucket, w)
	// if err != nil {
	// 	panic(err)
	// }
	lambda.Start(Handler)
}

func downloadResizeUpload(bucket, key string) error {
	serv, err := NewAWSService()

	if err != nil {
		return err
	}

	p := filepath.Base(key)

	dlPath, err := serv.DownloadImages3(bucket, p)
	if err != nil {
		return err
	}

	newKey := "resized-" + p

	rePath, err := resizeCopyImg(dlPath, newKey)
	if err != nil {
		return err
	}

	fmt.Println("rePath", rePath)

	if err := serv.UploadImages3(rePath, bucket+"-resized", newKey); err != nil {
		return err
	}

	// clean up
	return removeFiles(dlPath, rePath)
}

func resizeImg(img image.Image) image.Image {
	return resize.Resize(200, 0, img, resize.Bilinear)
}

func resizeCopyImg(src, destFileName string) (string, error) {
	sFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer sFile.Close()

	// decode jpeg into image.Image
	img, err := jpeg.Decode(sFile)
	if err != nil {
		return "", err
	}

	resized := resizeImg(img)

	resizedPath := makeTmpPath(destFileName)

	dFile, err := os.Create(resizedPath)
	if err != nil {
		return "", err
	}
	defer dFile.Close()

	// write new image to file
	return resizedPath, jpeg.Encode(dFile, resized, nil)

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

func (serv *awsService) DownloadImages3(bucket, key string) (string, error) {
	downloadPath := makeTmpPath(key)

	file, err := os.Create(downloadPath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(serv.session)

	s := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	numBytes, err := downloader.Download(file, &s)

	if err != nil {
		return downloadPath, fmt.Errorf("unable to download item %q, %v", key, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return downloadPath, nil
}

func (serv *awsService) UploadImages3(filePath, bucket, key string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(serv.session)

	input := s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	}

	_, err = uploader.Upload(&input)
	return err
}

func isJpegExtension(p string) bool {
	switch path.Ext(p) {
	case ".jpeg", ".jpg":
		return true
	default:
		return false
	}
}

func makeTmpPath(s string) string {
	return path.Join(os.TempDir(), s)
}

func removeFiles(paths ...string) error {
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

func NewAWSService() (awsService, error) {
	region := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return awsService{}, err
	}

	return awsService{
		session: sess,
	}, nil
}
