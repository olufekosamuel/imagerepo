package helpers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadFileToS3 saves a file to aws bucket and returns the url to // the file and an error if there's any
func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader, directory string) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := directory + fileHeader.Filename

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	MyBucket := os.Getenv("BUCKET_NAME")
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(tempFileName),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tempFileName, err
}

// DeleteFileFromS3 deletes a file from aws bucket and returns an error if there's any or true if deleted
func DeleteFileFromS3(s *session.Session, file string) (bool, error) {

	// create a unique file name for the file
	tempFileName := file

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	MyBucket := os.Getenv("BUCKET_NAME")
	_, err := s3.New(s).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(tempFileName),
	})
	if err != nil {
		fmt.Println(err)
		return true, err
	}

	return true, err
}
