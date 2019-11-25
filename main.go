package main

import (
	"flag"
	"fmt"

	//aws "github.com/ajitchahal/terraform-s3/aws"
	m "github.com/ajitchahal/terraform-s3/model"
	"github.com/ajitchahal/terraform-s3/tf"
	//"github.com/ajitchahal/terraform-s3/tf"

	"log"

	"github.com/octago/sflags/gen/gflag"
)

var cfg = &m.Config{
	S3: m.AwsS3{
		FileName: "terraform.tfstate",
		Region:   "us-east-1",
	},
}

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
	parseCmdLineArgs() //go run main.go --help

	

	// switch cfg.S3.Operation {
	// case "up":
	// 	aws.UploadToBucket(&cfg.S3)
	// case "dwn":
	// 	aws.DownloadFromBucket(&cfg.S3)
	// case "li":
	// 	aws.ListBucketItems(&cfg.S3)
	// case "l":
	// 	aws.ListBuckets(&cfg.S3)
	// }

	// //log.Println(string(buff))
	 tf.ReplacePasswordsWithPlaceHolders()

	// tf.ReplacePlaceHoldersWithPasswords(cfg.S3, "abcf")
}
