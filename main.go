package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func main() {
	flag.Parse()
	user := flag.Arg(0)

	if flag.NArg() == 2 {
		passcode := flag.Arg(1)
		data, err := ioutil.ReadFile(fname(user))
		if err != nil {
			log.Fatal(err)
		}
		key, err := otp.NewKeyFromURL(string(data))
		if err != nil {
			log.Fatal(err)
		}
		valid := totp.Validate(passcode, key.Secret())
		if valid {
			fmt.Println("code is valid")
		} else {
			log.Fatal("code is not valid")
		}
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "rwcr.net",
		AccountName: user,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Convert TOTP key into a QR code encoded as a PNG image.
	var buf bytes.Buffer
	img, err := key.Image(600, 600)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(&buf, img)
	ioutil.WriteFile("qr.png", buf.Bytes(), 0644)

	// Now Validate that the user's successfully added the passcode.
	passcode := promptForPasscode()
	valid := totp.Validate(passcode, key.Secret())

	if valid {
		fmt.Println("passcode correct - secret registered")
		ioutil.WriteFile(fname(user), []byte(key.String()), 0644)
	}
}

func fname(user string) string {
	return "secret-" + user + ".blog"
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}
