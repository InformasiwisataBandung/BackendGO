package gisbdg

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdateGetData(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "InformasiWisataBandung")
	datagedung := GetAllBangunanLineString(mconn, "InformasiWisataBandung")
	fmt.Println(datagedung)
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("salman", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "InformasiWisataBandung")
	var userdata User
	userdata.Username = "salman"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "Users", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "petapedia")
	var userdata User
	userdata.Username = "salman"
	userdata.Password = "secret"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestEnCode(t *testing.T) {
	anu, err := watoken.Encode("ini testing", pvtKey)
	fmt.Println(err)
	fmt.Println(anu)
}
func TestDcode(t *testing.T) {
	anu := watoken.DecodeGetId(pbcKey, token)
	fmt.Println(anu)
}

var pvtKey = "fbd4a28176db98361b4fb8936e2a1cb499bfe6a3760a3a0726fba735c6edac75513f6d70886d4abd6c475403e025afb4a053cc988cb8ba31ef062847e5e8b4d6"
var pbcKey = "513f6d70886d4abd6c475403e025afb4a053cc988cb8ba31ef062847e5e8b4d6"
var token = "v4.public.eyJleHAiOiIyMDIzLTEyLTA0VDEyOjEzOjMzKzA3OjAwIiwiaWF0IjoiMjAyMy0xMi0wNFQxMDoxMzozMyswNzowMCIsImlkIjoiaW5pIHRlc3RpbmciLCJuYmYiOiIyMDIzLTEyLTA0VDEwOjEzOjMzKzA3OjAwIn38C8vCUHizxSzELIDp4svHcQntBi8CsyrxDJ8j_pTpCsBHZBBSku8mioJ0qK5Gn2Z58aHXxW7x5B_dUHyZGG0L"
