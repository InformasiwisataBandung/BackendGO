package gisbdg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/whatsauth/wa"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Otorisasi(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response CredentialUser
	var auth User
	response.Status = false

	// Extract token from the request header
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Decode token values
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	// Create User struct with the decoded username
	auth.Username = tokenusername

	// Check if decoding results are valid
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Successful token decoding and user validation
	response.Message = "Berhasil decode token"
	response.Status = true
	response.Data.Username = tokenusername
	response.Data.Role = tokenrole

	return GCFReturnStruct(response)
}

func LoginHandler(token, privatekey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Decode user data from the request body
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, datauser) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the entered password is not valid
	if !IsPasswordValid(mconn, collname, datauser) {
		response.Message = "Password Salah"
		return GCFReturnStruct(response)
	}

	// Retrieve user details
	user := FindUser(mconn, collname, datauser)

	// Prepare and encode token
	tokenstring, tokenerr := Encode(user.Username, user.Role, os.Getenv(privatekey))
	if tokenerr != nil {
		response.Message = "Gagal encode token: " + tokenerr.Error()
		return GCFReturnStruct(response)
	}

	// Successful login
	response.Status = true
	response.Token = tokenstring
	response.Message = "Berhasil login"

	// Send a WhatsApp message notifying the user about the successful login
	var nama = user.Username
	var nohp = user.No_whatsapp
	dt := &wa.TextMessage{
		To:       nohp,
		IsGroup:  false,
		Messages: nama + " berhasil login\nPerlu diingat sesi login hanya berlaku 2 jam",
	}
	atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")

	return GCFReturnStruct(response)
}

func GCFPostHandlerSIGN(token, PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
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
func Registrasi(token, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Decode user data from the request body
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Check if the username already exists
	if usernameExists(MONGOCONNSTRINGENV, dbname, datauser) {
		response.Message = "Username telah dipakai"
		return GCFReturnStruct(response)
	}

	// Hash the user's password
	hash, hashErr := HashPassword(datauser.Password)
	if hashErr != nil {
		response.Message = "Gagal hash password: " + hashErr.Error()
		return GCFReturnStruct(response)
	}

	// Check if the 'No_whatsapp' field is empty
	if datauser.No_whatsapp == "" {
		response.Message = "Nomor We A wajib diisi"
		return GCFReturnStruct(response)
	}

	// Insert user data into the database
	InsertUserdata(mconn, collname, datauser.No_whatsapp, datauser.Username, hash, datauser.Role)
	response.Status = true
	response.Message = "Berhasil input data"

	// Prepare and send a WhatsApp message with registration details
	var username = datauser.Username
	var password = datauser.Password
	var nohp = datauser.No_whatsapp

	dt := &wa.TextMessage{
		To:       nohp,
		IsGroup:  false,
		Messages: "Registrasi Sukses buos, Username nya : " + username + "\nDengan Password yang dibuat adalah: " + password + "\nsimpan informasi berikut dengan baik",
	}

	// Make an API call to send WhatsApp message
	atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")

	return GCFReturnStruct(response)
}
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
func CreateWisata(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datawisata TempatWisata
	err := json.NewDecoder(r.Body).Decode(&datawisata)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	var auth User
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Decode token to get user details

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user has admin or author privileges
	if tokenrole != "admin" && tokenrole != "author" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	response.Status = true
	CreateWisataConn(mconn, collname, datawisata)
	response.Message = "Berhasil input data"
	return GCFReturnStruct(response)
}

// GET FIX
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
func UpdateWisata(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	// Koneksi
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var auth User
	var datawisata TempatWisata
	err := json.NewDecoder(r.Body).Decode(&datawisata)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Decode token to get user details
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user has admin or author privileges
	if tokenrole != "admin" && tokenrole != "author" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	response.Status = true
	UpdateWisataConn(mconn, collname, datawisata)
	response.Message = "Berhasil Update data Ya"
	return GCFReturnStruct(response)

}
func DeleteWisata(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var auth User
	var datawisata TempatWisata
	err := json.NewDecoder(r.Body).Decode(&datawisata)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Decode token to get user details
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user has admin or author privileges
	if tokenrole != "admin" && tokenrole != "author" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	response.Status = true
	DeleteWisataConn(mconn, collname, datawisata)
	response.Message = " Menghapus " + datawisata.Nama + "dari database"
	return GCFReturnStruct(response)
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
