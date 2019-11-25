package aws

import (
	"fmt"
	m "github.com/ajitchahal/terraform-s3/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

func getS3Session(awsS3 m.AwsS3) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsS3.Region)},
	)
	ifErrorExit("error connecting to s3", err)
	return sess
}
func downloadFromBucket(awsS3 m.AwsS3) {
	file, err := os.Create(awsS3.FileName)
	if err != nil {
		exitErrorf("Unable to create file %q, %v", err)
	}

	defer file.Close()
	downloader := s3manager.NewDownloader(getS3Session(awsS3))

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(awsS3.Bucket),
			Key:    aws.String(awsS3.FileName),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", awsS3.FileName, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}
func uploadToBucket(awsS3 m.AwsS3) {
	file, err := os.Open(awsS3.FileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()
	uploader := s3manager.NewUploader(getS3Session(awsS3))

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsS3.Bucket),
		Key:    aws.String(awsS3.FileName),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", awsS3.FileName, awsS3.Bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", awsS3.FileName, awsS3.Bucket)
}
func listBucketItems(awsS3 m.AwsS3) {
	// Create S3 service client
	svc := s3.New(getS3Session(awsS3))
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(awsS3.Bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", awsS3.Bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}
func listBuckets(awsS3 m.AwsS3) {
	// Create S3 service client
	svc := s3.New(getS3Session(awsS3))

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
