package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

func main() {
	//listBucketItems()
	//uploadToBucket()
	downloadFromBucket()
}
func downloadFromBucket(){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	ifErrorExit("error connecting to s3", err)
	bucket := "kiwis-resources-stage-us"

	filename := "README.md"
	file, err := os.Create(filename)
    if err != nil {
        exitErrorf("Unable to open file %q, %v", err)
    }

    defer file.Close()
	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", filename, err)
	}
	
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}
func uploadToBucket() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	ifErrorExit("error connecting to s3", err)
	bucket := "kiwis-resources-stage-us"

	filename := "README.md"

	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(filename),
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}
	
	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}
func listBucketItems() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	ifErrorExit("error connecting to s3", err)
	bucket := "kiwis-resources-stage-us"

	// Create S3 service client
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}
func listBuckets() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	ifErrorExit("error connecting to s3", err)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}
	ifErrorExit("error listing s3 buckets", err)

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		fmt.Println("....")
	}

	log.Println("----done----")
}

func ifErrorExit(msg string, err error) {
	if err != nil {
		exitErrorf(msg, err)
	}
}
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
