# Instalasi Go di Linux
Unduh arsip installer dari https://golang.org/dl/

Pilih installer untuk Linux yang sesuai dengan jenis bit komputer anda. Proses download bisa dilakukan lewat CLI, menggunakan wget atau curl.

$ wget https://storage.googleapis.com/golang/go1...

Buka terminal, extract arsip tersebut ke /usr/local.

$ tar -C /usr/local -xzf go1...

Tambahkan path binary Go ke PATH environment variable.

$ echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc

$ source ~/.bashrc

Selanjutnya, eksekusi perintah berikut untuk mengetes apakah Go sudah terinstal dengan benar

go version

Jika output adalah sama dengan versi Go yang ter-install, menandakan proses instalasi berhasil.

# Running
cd [folder-project]

ganti config di yaml rubah jadi dev.yaml / prod.yaml

go run main.go

# Install module di dalam go
go get [nama-module] -> go get github.com/xuanbo/eureka-client