package gisbdg

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

func GenerateKey() (privatekey, publickey string) {
	secretKey := paseto.NewV4AsymmetricSecretKey() // don't share this!!!
	privatekey = secretKey.ExportHex()             // DO share this one
	publickey = secretKey.Public().ExportHex()
	return privatekey, publickey
}

func Encode(name, username, role, privatekey string) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.SetString("username", username)
	token.SetString("role", role)
	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(privatekey)
	return token.V4Sign(key, nil), err
}

func Decode(publickey, tokenstr string) (payload Payload, err error) {
	var token *paseto.Token
	var pubKey paseto.V4AsymmetricPublicKey

	// Pastikan bahwa kunci publik dalam format heksadesimal yang benar
	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publickey)
	if err != nil {
		return payload, fmt.Errorf("failed to create public key: %s", err)
	}

	parser := paseto.NewParser()

	// Pastikan bahwa token memiliki format yang benar
	token, err = parser.ParseV4Public(pubKey, tokenstr, nil)
	if err != nil {
		return payload, fmt.Errorf("failed to parse token: %s", err)
	} else {
		// Handle token claims
		json.Unmarshal(token.ClaimsJSON(), &payload)
	}

	return payload, nil
}

func DecodeGetUsername(publickey string, tokenstring string) string {
	payload, err := Decode(publickey, tokenstring)
	if err != nil {
		fmt.Println("Decode DecodeGetId : ", err)
	}
	return payload.User
}

func DecodeGetRole(publickey string, tokenstring string) string {
	payload, err := Decode(publickey, tokenstring)
	if err != nil {
		fmt.Println("Decode DecodeGetId : ", err)
	}
	return payload.Role
}
