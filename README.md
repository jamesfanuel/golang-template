﻿# Struktur Folder
1. Folder Repository : Folder untuk membuat fungsi - fungsi yang masuk ke dalam database
2. Folder Service : Folder untuk membuat business logic dari aplikasi
3. Folder Controller : Folder untuk menjalankan fungsi dari router
4. Folder App : Berisikan config-config yang dijalankan ketika running aplikasi
5. Folder Docs : Hasil generate dari go swagger (OpenAPI) untuk dokumentasi API
6. Folder Exception : Untuk memberikan return error yang terjadi ketika ada panic yang dijalankan
7. Folder Model : Folder untuk membuat pemodelan data baik request maupun response juga untuk mendefinisikan tabel di dalam database
    - Folder ini dibagi menjadi domain (Tabel dalam database) dan web (Hasil pemodelan data dari komunikasi antara client dan server)
8. Folder Test : Berisikan unit test yang dapat dijalankan dengan running di setiap fungsinya atau mengetikan perintah : go test (harus masuk ke dalam folder test)
9. Folder Middleware : Berisikan fungsi untuk melakukan autentifikasi sebelum masuk ke dalam routing

# Di dalam template ini terdapat library sebagai berikut
1. Gorm : Object Relational Model populer dari Go , untuk melakukan pemodelan fungsi-fungsi database (Sama seperti hibernate di Spring)
2. Validator : Library untuk melakukan validasi data
3. Testify : Library untuk melakukan unit test
4. Google Wire : Untuk membuat Dependency Injection (Sama seperti Autowired di Spring)
5. Logrus : Untuk melakukan logging di setiap fungsi yang dijalankan
6. Go Swag : Untuk mengenerate swagger
7. Xuanbo/Eureka Client : Integrasi terhadap eureka
8. Httprouter : Manajemen routing
Seluruh fungsi ini harus diinstal terlebih dahulu untuk dijalankan

# Perintah - perintah dasar
1. Wire : Melakukan generate autowire , hasilnya akan menggenerate file berimana wire_gen.go. Apabila ingin melihat apa saja yang diwire bisa dicek ke injector.go
2. go run . : Running aplikasi umum
3. go run build : Build aplikasi
4. go env -w GOOS=windows | go env -w GOOS=linux : Melakukan perpindahan os ke windows / linux. (go run di windows apabila menggunakan windows dan go build ke linux karena servernya menggunakan linux)

# Langkah - langkah dasar
1. Pastikan pembuatan aplikasi dimulai dari membuat model-modelnya terlebih dahulu
2. Dilanjutkan dengan pembuatan repository (Secara umum hanya ada CRUD , tapi ada fungsi-fungsi khusus yang bisa dibuat juga apabila dibutuhkan)
3. Dilanjutkan dengan pembuatan service
4. Dilanjutkan dengan pembuatan controller
5. Modifikasi database
6. Modifikasi routing di dalam file app/router.go
7. Pastikan http portnya dirubah agar tidak bentrok ketika masuk ke eureka
8. Modifikasi annotation untuk kebutuhan generate swagger
9. Pastikan nama - nama repository, service, dan controller mengikuti nama servicenya untuk memudahkan dalam pemahaman alur aplikasi. Juga diwajibkan mengikuti standar penamaan baku dari aplikasi
