package gisbdg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GCFHandler(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datagedung := GetAllBangunanLineString(mconn, collectionname)
	return GCFReturnStruct(datagedung)
}

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
func CreateWisata(MONGOCONNSTRING, dbname, collectionname string, tempat TempatWisata) error {
	clientOptions := options.Client().ApplyURI(MONGOCONNSTRING)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(dbname).Collection(collectionname)

	_, err = collection.InsertOne(context.TODO(), tempat)
	if err != nil {
		return err
	}

	return nil
}
func ReadWisata(MONGOCONNSTRING, dbname, collectionname string) ([]TempatWisata, error) {
	clientOptions := options.Client().ApplyURI(MONGOCONNSTRING)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(dbname).Collection(collectionname)

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tempatList []TempatWisata

	for cursor.Next(context.TODO()) {
		var tempat TempatWisata
		err := cursor.Decode(&tempat)
		if err != nil {
			return nil, err
		}
		tempatList = append(tempatList, tempat)
	}

	return tempatList, nil
}
func UpdateWisata(MONGOCONNSTRING, dbname, collectionname string, filter bson.D, update bson.D) error {
	clientOptions := options.Client().ApplyURI(MONGOCONNSTRING)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(dbname).Collection(collectionname)

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func DeleteWisata(MONGOCONNSTRING, dbname, collectionname string, filter bson.D) error {
	clientOptions := options.Client().ApplyURI(MONGOCONNSTRING)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(dbname).Collection(collectionname)

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

// Geocoding (untuk menemukan lokasi dari konten yang sudah dibuat)
func Geocoding(MONGOCONNSTRINGENV, dbname, collectionname string, query string) ([]TempatWisata, error) {
	clientOptions := options.Client().ApplyURI(MONGOCONNSTRINGENV)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(dbname).Collection(collectionname)

	var filter bson.M

	// Jika query merupakan koordinat, maka cari berdasarkan koordinat
	if isCoordinates(query) {
		var coordinates [2]float64
		_, err := fmt.Sscanf(query, "[%f,%f]", &coordinates[0], &coordinates[1])
		if err != nil {
			return nil, err
		}

		// Buat filter untuk pencarian berdasarkan koordinat
		filter = bson.M{"lokasi.coordinates": coordinates}
	} else {
		// Jika query adalah nama, cari berdasarkan nama
		filter = bson.M{"nama": query}
	}

	var tempatList []TempatWisata
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var tempat TempatWisata
		if err := cursor.Decode(&tempat); err != nil {
			return nil, err
		}
		tempatList = append(tempatList, tempat)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tempatList, nil
}

// Fungsi untuk mengecek apakah input adalah koordinat atau bukan
func isCoordinates(input string) bool {
	var coordinates [2]float64
	_, err := fmt.Sscanf(input, "[%f,%f]", &coordinates[0], &coordinates[1])
	return err == nil
}
