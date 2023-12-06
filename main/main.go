package main

import (
	"fmt"
	"net/http"

	gisbdg "github.com/InformasiwisataBandung/BackendGO"
)

func main() {
	http.HandleFunc("/", HelloHTTP)
	http.ListenAndServe(":8080", nil)
}
func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Token")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var tempat gisbdg.TempatWisata // Membuat variabel tempat
	err := gisbdg.CreateWisata("publickey", "MONGOSTRING", "InformasiWisataBandung", "TempatWisata", tempat, r)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error()) // Mengonversi error ke string dan menampilkannya
		return
	}
	fmt.Fprintf(w, "Success") // Jika tidak ada error
}
