package awsservice

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/vidhill/terraform-lambda-play/util"
)

type awsService struct {
	Session *session.Session
}

func (serv *awsService) LoadFile(bucket, key string) (io.ReadCloser, string, error) {
	downloadPath := makeTmpPath(key)

	file, err := os.Create(downloadPath)
	if err != nil {
		return io.NopCloser(strings.NewReader("")), downloadPath, err
	}

	downloader := s3manager.NewDownloader(serv.Session)

	s := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	numBytes, err := downloader.Download(file, &s)

	if err != nil {
		return file, downloadPath, fmt.Errorf("unable to download item %q, %v", key, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	// "rewind" the reader so that it can be read from again by the jpeg decoder
	_, err = file.Seek(0, io.SeekStart)

	return file, downloadPath, err
}

func (serv *awsService) WriteFile(filePath, bucket, key string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return serv.WriteReaderContent(file, bucket, key)
}

func (serv *awsService) WriteReaderContent(r io.Reader, bucket, key string) error {

	uploader := s3manager.NewUploader(serv.Session)

	input := s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   r,
	}

	_, err := uploader.Upload(&input)
	return err
}

func (serv *awsService) RemoveFiles(paths ...string) error {
	return util.RemoveFiles(paths...)
}

func NewService() (awsService, error) {
	region := os.Getenv("AWS_REGION")

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return awsService{}, err
	}

	return awsService{
		Session: s,
	}, nil
}

func makeTmpPath(s string) string {
	return path.Join(os.TempDir(), s)
}
