# FindYourFit Backend

Backend API untuk aplikasi FindYourFit pada tugas UAS Aplikasi Mobile Lanjutan. Backend ini dipakai oleh aplikasi E-Commerce untuk autentikasi, produk, cart, order, dan status pembayaran.

## Repository Terkait

- Aplikasi FindYourFit: [WhoIsR/FashionApp_FindYourFit](https://github.com/WhoIsR/FashionApp_FindYourFit.git)

## Fitur

- Autentikasi pengguna.
- Data produk fashion.
- Keranjang belanja.
- Checkout dan order.
- Riwayat transaksi/order.
- Status pembayaran dari integrasi Kashi E Money.

## Menjalankan Backend

```bash
go mod download
go run .
```

Perintah `go mod download` digunakan untuk mengambil dependency Go. Perintah `go run .` digunakan untuk menjalankan backend secara lokal.

## Catatan Konfigurasi

Konfigurasi environment seperti database, Firebase, dan token disimpan di file `.env`. File `.env` tidak ikut repo karena berisi credential lokal.
