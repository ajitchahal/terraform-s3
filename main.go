package main

import (
	"flag"
	"github.com/ajitchahal/terraform-s3/aws"
	m "github.com/ajitchahal/terraform-s3/model"
	"github.com/ajitchahal/terraform-s3/tf"
	"github.com/octago/sflags/gen/gflag"
	"log"
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
	log.Println("tail:", flag.Args())
}

//var secretNamesMap = make(map[string]string)



func main() {
	parseCmdLineArgs() //go run main.go --help
	switch cfg.S3.Operation {
	case "dwn":
		aws.DownloadFromBucket(&cfg.S3)
		tf.ReplacePlaceHoldersWithPasswords()
	case "up":
		tf.ReplacePasswordsWithPlaceHolders()
		aws.UploadToBucket(&cfg.S3)
	case "li":
		aws.ListBucketItems(&cfg.S3)
	case "l":
		aws.ListBuckets(&cfg.S3)
	}
}
