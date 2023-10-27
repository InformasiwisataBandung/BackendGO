### Update Library GISWisataBandung

```sh
go get -u all
go mod tidy
git tag                                 #cek riwayat versi tag
git tag v1.1.0                          #set versi tag
git push origin --tags                  #push tag version ke repo
go list -m github.com/InformasiwisataBandung/BackendGO@v1.1.0   #publish ke PKG go Dev

go get github.com/InformasiwisataBandung/BackendGO@v1.1.0 #Jika ingin Menggunakan Package atau library
```