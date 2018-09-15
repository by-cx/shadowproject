package volumes

import (
	"github.com/mholt/archiver"
	"github.com/minio/minio-go"
	"github.com/satori/go.uuid"
	"io"
	"os"
	shadowerrors "shadowproject/common/errors"
)

// Struct to access source code in S3 storage. It handles mounting and unmounting but it's just
// for the compatibility with other volume drivers. In this case the mount method will copies
// the source code from S3 bucket into the target.
type S3Volume struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	SSL       bool
}

func (s *S3Volume) GetMinioClient() *minio.Client {
	minioClient, err := minio.New(s.Endpoint, s.AccessKey, s.SecretKey, s.SSL)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "s3 backend connection error",
		})
	}

	return minioClient
}

// Downloads the object from S3 bucket/source, unzips it and copies content into the target.
func (s *S3Volume) Mount(source string, target string) {
	client := s.GetMinioClient()

	//Get the object
	object, err := client.GetObject(s.Bucket, source, minio.GetObjectOptions{})
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "get object from s3 error",
		})
	}

	// Copying the file into tmp
	archivePath := "/tmp/" + uuid.NewV4().String() + ".zip"

	file, err := os.OpenFile(archivePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "unarchiving object error",
		})
	}

	_, err = io.Copy(file, object)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "copying the file from the object error",
		})
	}
	file.Close()

	// Unzipping
	err = archiver.Zip.Open(archivePath, target)

	// Cleaning
	os.Remove(archivePath)

}

// Deletes the target.
func (s *S3Volume) Umount(target string) {
	err := os.RemoveAll(target)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "umounting error",
		})
	}
}
