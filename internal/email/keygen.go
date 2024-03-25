package email

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"reflect"
)

func getPrivateKey() (keyString string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	keyBlock := pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}
	keyString = string(pem.EncodeToMemory(&keyBlock))

	return keyString, err
}

func getPublicKey(privateString string) (publicKeyString string, dkim string, err error) {
	privateBlock, rest := pem.Decode([]byte(privateString))
	if len(rest) > 0 {
		log.Fatal(len(rest))
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		log.Fatal("ParsePKCS8PrivateKey:", err)
	}

	if reflect.TypeOf(privateKey).String() != "*rsa.PrivateKey" {
		return "", "", fmt.Errorf("pkey is not *rsa.PrivateKey")
	}

	publicKey := privateKey.(*rsa.PrivateKey).Public()
	if reflect.TypeOf(publicKey).String() != "*rsa.PublicKey" {
		return "", "", fmt.Errorf("not rsa")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyBlock := pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}
	publicKeyString = string(pem.EncodeToMemory(&publicKeyBlock))

	dkim = "v=DKIM1;k=rsa;p=" + base64.StdEncoding.EncodeToString(publicKeyBytes)

	return publicKeyString, dkim, err
}

func writeToFile(data, filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func GenerateKeys() (err error) {
	privString, err := getPrivateKey()
	if err != nil {
		log.Fatal("generatePrivateKEY:", err)
	}

	pubString, dkim, err := getPublicKey(privString)
	if err != nil {
		return err
	}

	err = writeToFile(privString, "dkim.key")
	if err != nil {
		return err
	}

	err = writeToFile(pubString, "dkim.pub")
	if err != nil {
		return err
	}

	err = writeToFile(dkim, "dkim.txt")
	if err != nil {
		return err
	}

	// fmt.Println(privString)
	// fmt.Println(pubString)
	// fmt.Println(dkim)

	return nil
}
