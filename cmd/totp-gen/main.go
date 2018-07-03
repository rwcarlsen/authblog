package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

var envname = flag.String("env", "", "name of environment variable to reed")
var secret = flag.String("secret", "", "secret to generate totp code from")

func main() {
	log.SetFlags(0)
	flag.Parse()
	if *envname != "" {
		*secret = os.Getenv(*envname)
	}

	if *secret == "" {
		log.Fatal("secret cannot be empty")
	}

	code, err := totp.GenerateCode(*secret, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code)
}
