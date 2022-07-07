package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"time"

	"9m/consts"
	"9m/model"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mitchellh/mapstructure"
)

var Secret string = "secret"
var Alg jwa.SignatureAlgorithm = jwa.RS384
var PrivKey *rsa.PrivateKey
var PubKey jwk.Key
var KeySET jwk.Set
var RealKey jwk.Key

func GenerateRsaKeys() error {
	var err error

	PrivKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
		return err
	}

	PubKey, err = jwk.New(PrivKey.PublicKey)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func SaveRsaKeys() error {
	var err error

	privASN := x509.MarshalPKCS1PrivateKey(PrivKey)

	privBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privASN,
	})

	err = ioutil.WriteFile(JwtPrivateKeyFileName, privBytes, 0644)
	if err != nil {
		return err
	}

	pubASN, err := x509.MarshalPKIXPublicKey(&PrivKey.PublicKey)
	if err != nil {
		return err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN,
	})

	err = ioutil.WriteFile(JwtPublicKeyFileName, pubBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadRsaKeys() error {
	var err error

	privBytes, err := ioutil.ReadFile(JwtPrivateKeyFileName)
	if err != nil {
		return err
	}

	privASN, _ := pem.Decode(privBytes)

	PrivKey, err = x509.ParsePKCS1PrivateKey(privASN.Bytes)
	if err != nil {
		return err
	}

	pubBytes, err := ioutil.ReadFile(JwtPublicKeyFileName)
	if err != nil {
		return err
	}

	pubASN, _ := pem.Decode(pubBytes)

	publicKey, err := x509.ParsePKIXPublicKey(pubASN.Bytes)
	if err != nil {
		return err
	}

	PubKey, err = jwk.New(publicKey)
	if err != nil {
		return err
	}

	return nil
}

// GenerateKeySet - Generate key set
func GenerateKeySet() error {
	var err error

	PubKey.Set(jwk.AlgorithmKey, Alg)
	PubKey.Set(jwk.KeyIDKey, Secret)

	bogusKey := jwk.NewSymmetricKey()
	bogusKey.Set(jwk.AlgorithmKey, jwa.NoSignature)
	bogusKey.Set(jwk.KeyIDKey, "otherkey")

	KeySET = jwk.NewSet()
	KeySET.Add(PubKey)
	KeySET.Add(bogusKey)

	RealKey, err = jwk.New(PrivKey)
	if err != nil {
		log.Printf("failed to create JWK: %s\n", err)
		return err
	}

	RealKey.Set(jwk.KeyIDKey, Secret)
	RealKey.Set(jwk.AlgorithmKey, Alg)

	return nil
}

// GenerateToken - Generate token
func GenerateToken(authinfo model.AuthInfo) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer(consts.ProgramName).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(time.Duration(authinfo.Duration.Int64)*time.Second)).
		Subject("auth token").
		Claim("token", authinfo).
		Build()

	if err != nil {
		log.Println("GenerateToken token", err)
		return "", err
	}

	signed, err := jwt.Sign(token, Alg, RealKey)
	if err != nil {
		log.Println("GenerateToken sign", RealKey, "/", err)
		return "", err
	}

	result := string(signed)

	return result, err
}

// ParseToken - Parse token
func ParseToken(payloadSTR string) (jwt.Token, model.AuthInfo, error) {
	payload := []byte(payloadSTR)

	token, err := jwt.Parse(
		payload,
		jwt.WithKeySet(KeySET),
	)
	if err != nil {
		// log.Printf("ParseToken parse payload: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	var authinfo model.AuthInfo

	cfg := &mapstructure.DecoderConfig{
		Result:     &authinfo,
		DecodeHook: ConvertToNullTypeHookFunc,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		// log.Printf("ParseToken decoder set: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	claim, valid := token.Get("token")
	if !valid {
		// log.Printf("ParseToken token to claim: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	err = decoder.Decode(claim)
	if err != nil {
		// log.Printf("ParseToken decode claim to struct: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	now := time.Now()
	if token.Expiration().Before(now) {
		// log.Printf("ParseToken token expired: %s\n", err)
		return nil, model.AuthInfo{}, errors.New("token expired")
	}

	return token, authinfo, err
}
