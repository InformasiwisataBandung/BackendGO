### Update Library GISWisataBandung

```sh
go get -u all
go mod tidy
git tag                                 #cek riwayat versi tag
git tag v1.0.0                          #set versi tag
git push origin v1.0.0     #push tag version ke repo
go list -m github.com/urlgithubanda@v1.0.0   #publish ke PKG go Dev

go get github.com/InformasiwisataBandung/BackendGO@v1.1.0 #Jika ingin Menggunakan Package atau library
```

#### Update enkripsi Password Mongo DB, Deploy function Signup, Memasukan Token jika user berhasil login kedalam Cookies
#### •Update 21-10-2023

```
{
  "username": "ucup",
  "password": "$2a$10$r.Z8w/WHkd7uHcE6ZGlqCOcsNQEQOdXyrYYcDMMY9V4/HLOmXloCq"
}
```
#### Update 23-10-2023
-API clooud functions Signup
```
https://asia-southeast2-bustling-walker-340203.cloudfunctions.net/function-Signup
```
input
```
{
  "username": "username",
  "password": "Password"
}
```
Send Post
```
{
    "message": "Pendaftaran berhasil"
}
```
-Memasukan Token ke Cookies
```
cookie := http.Cookie{
		Name:     "token",     // Nama cookie
		Value:    tokenString, // Token sebagai nilai cookie
		HttpOnly: true,        // Hanya bisa diakses melalui HTTP
		Path:     "/",         // Path di mana cookie berlaku (misalnya, seluruh situs)
		MaxAge:   3600,        // Durasi cookie (dalam detik), sesuaikan sesuai kebutuhan
		// Secure: true, // Jika situs dijalankan melalui HTTPS
	}
```
#### Update 4-5 November 2023

API CreateWisata
```
https://us-central1-bustling-walker-340203.cloudfunctions.net/function-6CreateWisata
```
Input
```
{
  "nama": "Nama Tempat Wisata",
  "jenis": "Jenis Tempat",
  "deskripsi": "Deskripsi Tempat Wisata",
  "lokasi": {
    "latitude": 123.456789,
    "longitude": 987.654321
  },
  "alamat": "Alamat Tempat Wisata",
  "gambar": "URL_Gambar",
  "rating": 4.5
}

```
SendPost CreateWisata
```
{
    "message": "Data Create successfully"
}
```
API ReadWisata
```
https://us-central1-bustling-walker-340203.cloudfunctions.net/function-7ReadWisata
```

Send GET Wisata
```
{
    "data": [
        {
            "nama": "Nama Tempat Wisata",
            "jenis": "Jenis Tempat",
            "deskripsi": "Deskripsi Tempat Wisata",
            "lokasi": {
                "type": "",
                "coordinates": null
            },
            "alamat": "Alamat Tempat Wisata",
            "gambar": "URL_Gambar",
            "rating": 4.5
        },
        {
            "nama": "Gedung Sate",
            "jenis": "Traveling",
            "deskripsi": "Tempat Bersejarah",
            "lokasi": {
                "type": "",
                "coordinates": null
            },
            "alamat": "Alamat Tempat Wisata",
            "gambar": "URL_Gambar",
            "rating": 4.5
        }
    ]
}
```
API UpdateWisata
```
https://us-central1-bustling-walker-340203.cloudfunctions.net/function-8UpdateWisata
```
Input
```
{
  "filter": {
    "Nama": "Tempat A"
  },
  "update": {
    "$set": {
      "Deskripsi": "Deskripsi baru untuk Tempat A"
    }
  }
}
```
SendPost UpdateWisata
```
{
    "message": "Data updated successfully"
}
```