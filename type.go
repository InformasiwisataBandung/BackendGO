package gisbdg

import "time"

type GeometryPolygon struct {
	Coordinates [][][]float64 `json:"coordinates" bson:"coordinates"`
	Type        string        `json:"type" bson:"type"`
}

type GeometryLineString struct {
	Coordinates [][]float64 `json:"coordinates" bson:"coordinates"`
	Type        string      `json:"type" bson:"type"`
}

type GeometryPoint struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Type        string    `json:"type" bson:"type"`
}

type GeoJsonLineString struct {
	Type       string             `json:"type" bson:"type"`
	Properties Properties         `json:"properties" bson:"properties"`
	Geometry   GeometryLineString `json:"geometry" bson:"geometry"`
}

type GeoJsonPolygon struct {
	Type       string          `json:"type" bson:"type"`
	Properties Properties      `json:"properties" bson:"properties"`
	Geometry   GeometryPolygon `json:"geometry" bson:"geometry"`
}

type Geometry struct {
	Coordinates interface{} `json:"coordinates" bson:"coordinates"`
	Type        string      `json:"type" bson:"type"`
}
type GeoJson struct {
	Type       string     `json:"type" bson:"type"`
	Properties Properties `json:"properties" bson:"properties"`
	Geometry   Geometry   `json:"geometry" bson:"geometry"`
}

type Properties struct {
	Name string `json:"name" bson:"name"`
}

type User struct {
	No_whatsapp string `json:"no_whatsapp,omitempty" bson:"no_whatsapp"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password,omitempty" bson:"password"`
	Role        string `json:"role,omitempty" bson:"role,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}
type TempatWisata struct {
	Nama      string  `json:"nama"`
	Jenis     string  `json:"jenis"`
	Deskripsi string  `json:"deskripsi"`
	Lokasi    Lokasi  `json:"lokasi"`
	Alamat    string  `json:"alamat"`
	Gambar    string  `json:"gambar"`
	Rating    float64 `json:"rating"`
}

type Lokasi struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
type CredentialUser struct {
	Status  bool   `json:"status" bson:"status"`
	Data    User   `json:"data,omitempty" bson:"data,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}
type Payload struct {
	Username string    `json:"username"`
	Role string    `json:"role"`
	Exp  time.Time `json:"exp"`
	Iat  time.Time `json:"iat"`
	Nbf  time.Time `json:"nbf"`
}
type BeriPesan struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
}
