package tf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	m "github.com/ajitchahal/terraform-s3/model"
)

//ParseConfig parses secrets-config.json e.g. {"secretnames": ["A","sdsdsds","e","d"]}
func ParseConfig() m.TfSecrets {
	file, err := ioutil.ReadFile("secrets-config.json")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	secretsConfig := m.TfSecrets{}
	fmt.Println(string(file))
	err = json.Unmarshal([]byte(file), &secretsConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	return secretsConfig
}

//ReplacePasswordsWithPlaceHolders replaces secret with place holder such as "secret_master_password":"abcdefgh3242" with "secret_master_password":"secret_secret_master_password_end"
func ReplacePasswordsWithPlaceHolders() {
	buff, err := ioutil.ReadFile("terraform.tfstate")
	if err != nil {
		log.Fatal(err)
	}
	fileContent := string(buff)

	secretsToReplace := ParseConfig()
	for i := 0; i < len(secretsToReplace.SecretNames); i++ {

		secretToReplace := secretsToReplace.SecretNames[i]
		regEx := fmt.Sprintf(`(?m)\b%s":( |\t)*"\w+",$`, secretToReplace)

		fmt.Println(regEx)
		re := regexp.MustCompile(regEx)

		template := "%s\": \"%s\","
		templatedStr := fmt.Sprintf(template, secretToReplace, "secret_"+secretToReplace+"_end")
		fileContent = re.ReplaceAllString(fileContent, templatedStr)
	}
	ioutil.WriteFile("terraform.tfstate", []byte(fileContent), 0644)
}

//ReplacePlaceHoldersWithPasswords replaces placeholder with provicded secert such as "secret_master_password":"secret_secret_master_password_end" with "secret_master_password":"abcdefgh3242"
func ReplacePlaceHoldersWithPasswords(secret string) {
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
