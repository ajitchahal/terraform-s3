package main

import (
	"flag"
	"fmt"
	m "github.com/ajitchahal/terraform-s3/model"
	"github.com/ajitchahal/terraform-s3/tf"
	"github.com/octago/sflags/gen/gflag"
	"log"
)

var cfg = &m.Config{
	S3: m.AwsS3{
		//kiwis-resources-stage-us
		FileName: "terraform.tfstate",
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
	// parseCmdLineArgs() //go run main.go --help

	// switch cfg.S3.Operation {
	// case "up":
	// 	downloadFromBucket()
	// case "down":
	// 	downloadFromBucket()
	// }

	//log.Println(string(buff))
	tf.ReplacePasswordsWithPlaceHolders(cfg.S3)

	tf.ReplacePlaceHoldersWithPasswords(cfg.S3, "abcf")
}
