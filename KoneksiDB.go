package gisbdg

import (
	"context"
	"os"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}
//User
func IsPasswordValid(MONGOCONNSTRINGENV *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](MONGOCONNSTRINGENV, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}
func usernameExists(MONGOCONNSTRINGENV, dbname string, userdata User) bool {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname).Collection("Users")
	filter := bson.M{"username": userdata.Username}

	var user User
	err := mconn.FindOne(context.Background(), filter).Decode(&user)
	return err == nil
}
func InsertUserdata(MONGOCONNSTRINGENV *mongo.Database, collname, no_whatsapp, username, password, role string) (InsertedID interface{}) {
	req := new(User)
	req.No_whatsapp = no_whatsapp
	req.Username = username
	req.Password = password
	req.Role = role
	return atdb.InsertOneDoc(MONGOCONNSTRINGENV, collname, req)
}

func ReadUser(MONGOCONNSTRINGENV *mongo.Database, collname string)[]User{
	user:= atdb.GetAllDoc[[]User](MONGOCONNSTRINGENV, collname)
	return user
}

func EditUser(MONGOCONNSTRINGENV *mongo.Database, collname string, datauser User) interface {}{
	filter := bson.M{"username": datauser.Username}
	return atdb.ReplaceOneDoc(MONGOCONNSTRINGENV, collname, filter, datauser)
}

func HapusUser(MONGOCONNSTRINGENV *mongo.Database, collname string, userdata User)interface{}{
	filter := bson.M{"username": userdata.Username}
	return atdb.DeleteOneDoc(MONGOCONNSTRINGENV,collname,filter)
}
// Crud Connection Wisata

func CreateWisataConn(MONGOCONNSTRINGENV *mongo.Database, collname string, datawisata TempatWisata) interface{} {
	return atdb.InsertOneDoc(MONGOCONNSTRINGENV, collname, datawisata)
}

func UpdateWisataConn(MONGOCONNSTRINGENV *mongo.Database, collname string, datawisata TempatWisata) interface{} {
	filter := bson.M{"nama": datawisata.Nama}
	return atdb.ReplaceOneDoc(MONGOCONNSTRINGENV, collname, filter, datawisata)
}

func DeleteWisataConn(MONGOCONNSTRINGENV *mongo.Database, collname string, datawisata TempatWisata) interface{} {
	filter := bson.M{"nama": datawisata.Nama}
	return atdb.DeleteOneDoc(MONGOCONNSTRINGENV, collname, filter)
}
func FindUser(MONGOCONNSTRINGENV *mongo.Database, collname string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	return atdb.GetOneDoc[User](MONGOCONNSTRINGENV, collname, filter)
}

// Read All wisata
func GetAllWisata(MONGOCONNSTRINGENV *mongo.Database, collname string) []TempatWisata {
	tempat := atdb.GetAllDoc[[]TempatWisata](MONGOCONNSTRINGENV, collname)
	return tempat
}
func FindWisat(MONGOCONNSTRINGENV *mongo.Database, collname string, datawisata TempatWisata) TempatWisata {
	filter := bson.M{"nama": datawisata.Nama}
	return atdb.GetOneDoc[TempatWisata](MONGOCONNSTRINGENV, collname, filter)
}

func NamaWisataExist(MONGOCONNSTRINGENV, dbname string, datawisata TempatWisata)bool{
	mconn:= SetConnection(MONGOCONNSTRINGENV, dbname).Collection("TempatWisata")
	filter:= bson.M{"nama": datawisata.Nama}

	var wisata TempatWisata
	err := mconn.FindOne(context.Background(), filter).Decode(&wisata)
	return	err == nil
}