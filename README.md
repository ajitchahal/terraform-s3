# terraform-s3
A cli to store and retrieve tfstate to/from S3

##
comand line flags

[blog about cmd flags](https://golangexample.com/generate-flags-by-parsing-structures/)

[git repo of the above project](https://github.com/octago/sflags/)

[get cmd line vars from env variables](https://github.com/jamiealquiza/envy)

# Note: still not a proper solution of

get variable from ENV_variable if set, then from command line args anf parse to a structure.

## Replacing password

option 1
```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`\bsecret_\w+":"\w+",$`)
	fmt.Println(re.ReplaceAllString(`"secret_master_password":"abcdefgh3242",`, "secret_master_password\":\"pass\""))
}

		re := regexp.MustCompile(`(?m)\bsecret_\w+":( |\t)*"\w+",$`)

```
option 2
```go
var or = `
 "iam_roles": null,
            "id": "abntest-service",
            "kms_key_id": "arn:aws:kms:us-east-1:dd:key/db",
            "master_password": "secret_abc1234_end",
            "master_username": "masteruser",
            "port": 5432,`

	re := regexp.MustCompile(`\bsecret_\w+_end`)
	fmt.Println(re.ReplaceAllString(or, "pass"))
```