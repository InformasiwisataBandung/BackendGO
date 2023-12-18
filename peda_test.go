package gisbdg

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

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
	fmt.Println("ini pvtkey",privateKey)
	fmt.Println("ini pbckey",publicKey)
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

var pvtKey = "b6dee19a1238c4a97ddd9d937a4952e3c2cccadc21a66efdd6ceaaae6374d8d73f2c1fdbe947962ca245abca60431becf618ca8d5803ab6c71e0c96df3740411"
var pbcKey = "3f2c1fdbe947962ca245abca60431becf618ca8d5803ab6c71e0c96df3740411"
var token = "v4.public.eyJleHAiOiIyMDIzLTEyLTExVDAwOjA2OjUyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMi0xMFQyMjowNjo1MiswNzowMCIsImlkIjoic2FsbWFuIiwibmJmIjoiMjAyMy0xMi0xMFQyMjowNjo1MiswNzowMCJ9A3aJitjwFa4asInrwmkvazquVgUh0sGji_sYxs5EtNEh9uHF4_3-UPtcZ-iHomtMWWOXAvPiGPKoJJ644prqDg"
