package gisbdg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

/* Bagian Awal */
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
		Messages: nama + " berhasil login\n Nikmati Web Wisata di kota bandung ",
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

	var username = datauser.Username
	var password = datauser.Password
	var nohp = datauser.No_whatsapp

	dt := &wa.TextMessage{
		To:       nohp,
		IsGroup:  false,
		Messages: "Registrasi Sukses buos, Username nya : " + username + "\nDengan Password yang dibuat adalah: " + password + "\nsimpan informasi berikut dengan baik",
	}

	atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")

	return GCFReturnStruct(response)
}

// Bagian Akhir Signin Singnup & otorisasi

// User Edit Read Delete

func ReadsatuUser(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	//koneksi
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var auth User
	var userdata User
	err := json.NewDecoder(r.Body).Decode(&userdata)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header token tidak ditemukan"
		return GCFReturnStruct(response)
	}
	// Decode token untuk GET username dan role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	if usernameExists(MONGOCONNSTRINGENV, dbname, userdata) {
		// fetch wisata dari database
		user := FindUser(mconn, collname, userdata)
		return GCFReturnStruct(user)
	} else {
		response.Message = "User tidak ditemukan"
		return GCFReturnStruct(response)
	}
}

func ReadUserHandler(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}

	// Get all users if the user is an admin
	datauser := GetAllUser(mconn, collname)
	return GCFReturnStruct(datauser)
}

func UpdateUser(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Decode user data from the request body
	var auth User
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header token tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ada adalam database"
		return GCFReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}

	// Check if the username parameter is provided
	if datauser.Username == "" {
		response.Message = "Parameter dari function ini adalah username"
		return GCFReturnStruct(response)
	}

	// Check if the user to be edited exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, datauser) {
		response.Message = "Akun yang ingin diedit tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Hash the user's password if provided
	if datauser.Password != "" {
		hash, hashErr := HashPassword(datauser.Password)
		if hashErr != nil {
			response.Message = "Gagal Hash Password: " + hashErr.Error()
			return GCFReturnStruct(response)
		}
		datauser.Password = hash
	} else {
		// Retrieve user details
		user := FindUser(mconn, collname, datauser)
		datauser.Password = user.Password
	}

	// Perform user update
	EditUser(mconn, collname, datauser)

	response.Status = true
	response.Message = "Berhasil update " + datauser.Username + " dari database"
	return GCFReturnStruct(response)
}

func DeleteUser(publickey, MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Decode user data from the request body
	var auth User
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
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

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}

	// Check if the username parameter is provided
	if datauser.Username == "" {
		response.Message = "Parameter dari function ini adalah username"
		return GCFReturnStruct(response)
	}

	// Check if the user to be deleted exists
	if !usernameExists(MONGOCONNSTRINGENV, dbname, datauser) {
		response.Message = "Akun yang ingin dihapus tidak ditemukan"
		return GCFReturnStruct(response)
	}

	// Perform user deletion
	HapusUser(mconn, collname, datauser)

	response.Status = true
	response.Message = "Berhasil hapus " + datauser.Username + " dari database"
	return GCFReturnStruct(response)
}

// Akhir EDIT UPDATE DELETE USER

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// WISATA
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
		response.Message = "Header token tidak ditemukan"
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

	// Check if the user has admin or user privileges
	if tokenrole != "admin" && tokenrole != "user" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	response.Status = true
	CreateWisataConn(mconn, collname, datawisata)
	response.Message = "Berhasil input data"
	return GCFReturnStruct(response)
}

// GET FIX
func ReadWisata(MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	//koneksi
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	//ngambil semua tempat wisata
	datawisata := GetAllWisata(mconn, collname)
	return GCFReturnStruct(datawisata)
}
func ReadOnWisata(MONGOCONNSTRINGENV, dbname, collname string, r *http.Request) string {
	var response BeriPesan
	response.Status = false

	//koneksi
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datawisata TempatWisata
	err := json.NewDecoder(r.Body).Decode(&datawisata)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	if datawisata.Nama == "" {
		response.Message = "Isi dengan Field Nama"
		return GCFReturnStruct(response)
	}

	if NamaWisataExist(MONGOCONNSTRINGENV, dbname, datawisata) {
		// fetch wisata dari database
		wisata := FindWisat(mconn, collname, datawisata)
		return GCFReturnStruct(wisata)
	} else {
		response.Message = "Belum Mendapatkan Informasi Wisata"
	}
	return GCFReturnStruct(response)

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

	// Check if the user has admin or user privileges
	if tokenrole != "admin" && tokenrole != "user" {
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

	// Check if the user has admin or user privileges
	if tokenrole != "admin" && tokenrole != "user" {
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
func Geocode(address, apiKey string) (string, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", address, apiKey)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result map[string]interface{}
	if err := decodeJSON(response.Body, &result); err != nil {
		return "", err
	}

	geometry := result["results"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"].(map[string]interface{})
	lat := fmt.Sprintf("%v", geometry["lat"])
	lng := fmt.Sprintf("%v", geometry["lng"])

	return fmt.Sprintf("Latitude: %s, Longitude: %s", lat, lng), nil
}

func decodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
//geocoding handler
func GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing 'address' parameter", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		http.Error(w, "Missing Google Maps API key", http.StatusInternalServerError)
		return
	}

	result, err := Geocode(address, apiKey)
	if err != nil {
		http.Error(w, "Geocoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"result": "%s"}`, result)
}