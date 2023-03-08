package localfile

import (
	"fmt"
	"io"
	"os"
)

type service struct {
	outputPath string
}

func (s service) LoadFile(bucket, key string) (io.ReadCloser, string, error) {
	f, err := os.Open(key)
	fmt.Println("opening", key)

	return f, key, err
}

func (s service) WriteReaderContent(r io.Reader, bucket, key string) error {

	f, err := os.Create(s.outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)

	return err
}

func (s service) CleanUp() error {
	return os.Remove(s.outputPath)
}

func (s service) WriteFile(filePath, bucket, key string) error {
	return nil
}

func (s service) RemoveFiles(...string) error {
	return nil
}

func NewLocalServ(s string) service {
	return service{
		outputPath: s,
	}
}
