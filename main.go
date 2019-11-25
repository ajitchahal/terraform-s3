package main

import (
	"flag"
	"time"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/octago/sflags/gen/gflag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type awsS3 struct {
	Bucket    string `desc:"S3 bucket name"`
	FileName  string `desc:"File to upload or download"`
	Operation string `desc:"up|dwn upload or download file" flag:"op"`
	Region    string `desc:"AWS region the bucket is in e.g us-east-1"`
}
type config struct {
	S3 awsS3
}

var cfg = &config{}

func parseCmdLineArgs() {
	err := gflag.ParseToDef(cfg) //parseEnv
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	flag.Parse()
	log.Println(cfg)
	fmt.Println("tail:", flag.Args())
}
func main() {
	// parseCmdLineArgs() //go run main.go --help

	// switch cfg.S3.Operation {
	// case "up":
	// 	downloadFromBucket()
	// case "down":
	// 	downloadFromBucket()
	// }
	buff, err := ioutil.ReadFile("terraform.tfstate")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(string(buff))
	fileContent := string(buff)
	replacePasswordsWithPlaceHolders(fileContent)
	time.Sleep(2 * time.Second)
	buff, err = ioutil.ReadFile("terraform.tfstate")
	if err != nil {
		log.Fatal(err)
	}
	fileContent = string(buff)

	replacePlaceHoldersWithPasswords(fileContent, "ajit")
}
func replacePasswordsWithPlaceHolders(fileContent string) {
	secretToReplace := "master_password"
	regEx := fmt.Sprintf(`(?m)\b%s":( |\t)*"\w+",$`, secretToReplace)

	fmt.Println(regEx)
	re := regexp.MustCompile(regEx)

	template := "%s\": \"%s\","
	templatedStr := fmt.Sprintf(template, secretToReplace, "secret_" + secretToReplace + "_end")
	fileContent = re.ReplaceAllString(fileContent, templatedStr)

	//fmt.Println(fileContent)
	ioutil.WriteFile("terraform.tfstate", []byte(fileContent), 0644)
}
func replacePlaceHoldersWithPasswords(fileContent string, secret string) {
	secretToReplace := "master_password"
	regEx := fmt.Sprintf(`\bsecret_%s_end`, secretToReplace)

	fmt.Println(regEx)
	re := regexp.MustCompile(regEx)
	fileContent = re.ReplaceAllString(fileContent, secret)
	ioutil.WriteFile("terraform.tfstate", []byte(fileContent), 0644)
}
func downloadFromBucket() {
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
		Key:    aws.String(filename),
		Body:   file,
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
