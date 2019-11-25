package tf

import (
	"fmt"
	m "github.com/ajitchahal/terraform-s3/model"
	"io/ioutil"
	"log"
	"regexp"
)

//ReplacePasswordsWithPlaceHolders replaces secret with place holder such as "secret_master_password":"abcdefgh3242" with "secret_master_password":"secret_secret_master_password_end"
func ReplacePasswordsWithPlaceHolders(s3 m.AwsS3) {
	buff, err := ioutil.ReadFile("terraform.tfstate")
	if err != nil {
		log.Fatal(err)
	}
	fileContent := string(buff)

	secretToReplace := "master_password"
	regEx := fmt.Sprintf(`(?m)\b%s":( |\t)*"\w+",$`, secretToReplace)

	fmt.Println(regEx)
	re := regexp.MustCompile(regEx)

	template := "%s\": \"%s\","
	templatedStr := fmt.Sprintf(template, secretToReplace, "secret_"+secretToReplace+"_end")
	fileContent = re.ReplaceAllString(fileContent, templatedStr)

	ioutil.WriteFile("terraform.tfstate", []byte(fileContent), 0644)
}

//ReplacePlaceHoldersWithPasswords replaces placeholder with provicded secert such as "secret_master_password":"secret_secret_master_password_end" with "secret_master_password":"abcdefgh3242"
func ReplacePlaceHoldersWithPasswords(s3 m.AwsS3, secret string) {
	buff, err := ioutil.ReadFile("terraform.tfstate")
	if err != nil {
		log.Fatal(err)
	}
	fileContent := string(buff)

	secretToReplace := "master_password"
	regEx := fmt.Sprintf(`\bsecret_%s_end`, secretToReplace)

	fmt.Println(regEx)
	re := regexp.MustCompile(regEx)
	fileContent = re.ReplaceAllString(fileContent, secret)
	ioutil.WriteFile("terraform.tfstate", []byte(fileContent), 0644)
}
