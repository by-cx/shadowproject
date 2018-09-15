package volumes

import (
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

// Test the S3 volume driver
func TestS3Volume_Mount(t *testing.T) {
	s3driver := S3Volume{
		Endpoint:  "localhost:9000",
		Bucket:    "shadowtest",
		AccessKey: "3D2U2V66A3CP0CB088Z3",
		SecretKey: "Uipi4szPTGhyjoTFsmtXJrIf9cbqnfLRPQH6e8Ho",
		SSL:       false,
	}

	// Create the destination directory
	os.MkdirAll("/tmp/shadowtest/mnt", 0755)

	// Set the S3 environment to do the test
	client := s3driver.GetMinioClient()
	client.MakeBucket("shadowtest", "")

	file, err := os.OpenFile("../../contrib/s3_test_archive.zip", os.O_RDONLY, 0)
	assert.Nil(t, err)
	stat, err := file.Stat()
	assert.Nil(t, err)

	_, err = client.PutObject("shadowtest", "s3_test_archive.zip", file, stat.Size(), minio.PutObjectOptions{})
	assert.Nil(t, err)

	// Test mounting
	s3driver.Mount("s3_test_archive.zip", "/tmp/shadowtest/mnt")

	content, err := ioutil.ReadFile("/tmp/shadowtest/mnt/README.md")
	assert.Nil(t, err)

	assert.Equal(t, []byte("The S3 driver is working.\n"), content)

	// Test umounting
	s3driver.Umount("/tmp/shadowtest/mnt")

	_, err = ioutil.ReadFile("/tmp/shadowtest/mnt/README.md")
	assert.NotNil(t, err)
}
