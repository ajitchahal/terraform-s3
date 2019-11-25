package aws

import (
	m "github.com/ajitchahal/terraform-s3/model"
	"testing"
)

var cfg = &m.Config{
	S3: m.AwsS3{
		FileName: "terraform.tfstate",
		Region:   "us-east-1",
		Bucket:   "kiwis-resources-stage-us",
	},
}

func TestDownload(t *testing.T) {
	DownloadFromBucket(&cfg.S3)
}
func TestUpload(t *testing.T) {
	UploadToBucket(&cfg.S3)
}
