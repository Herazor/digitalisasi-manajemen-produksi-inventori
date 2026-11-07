# Digitalisasi Manajemen Produksi & Inventori
### Omnichannel & Bill of Materials (BOM)

## Deskripsi Proyek

Aplikasi Manajemen Produksi dan Inventori merupakan sistem yang dirancang untuk membantu proses pengelolaan produksi, persediaan bahan baku, produk jadi, serta pemrosesan pesanan penjualan secara terintegrasi.

Sistem ini menggabungkan pengelolaan Bill of Materials (BOM), Work Order, Sales Order, dan inventori untuk meningkatkan efisiensi operasional serta menjaga akurasi persediaan.

## Fitur Utama

### 1. Manajemen Material
- Menambahkan data bahan baku.
- Melihat daftar bahan baku.
- Memperbarui informasi bahan baku.
- Menghapus data bahan baku.
- Mengelola kode, nama, satuan, dan deskripsi material.

### 2. Manajemen Produk
- Menambahkan produk jadi.
- Melihat daftar produk.
- Memperbarui informasi produk.
- Menghapus data produk.
- Memantau persediaan produk jadi.

### 3. Bill of Materials (BOM)
- Membuat BOM untuk setiap produk.
- Mendukung pengelolaan versi BOM.
- Menghitung kebutuhan bahan baku melalui BOM Explosion.
- Menyimpan informasi versi BOM yang digunakan dalam Work Order.

### 4. Manajemen Inventori
- Memantau persediaan bahan baku.
- Memantau persediaan produk jadi.
- Menampilkan jumlah stok On Hand, Reserved, dan Available.
- Mencatat penambahan stok bahan baku.
- Menggunakan transaksi database dan penguncian baris untuk menjaga konsistensi persediaan.

### 5. Manajemen Work Order
- Membuat perintah produksi.
- Menghitung kebutuhan bahan baku berdasarkan BOM.
- Melakukan reservasi bahan baku secara otomatis.
- Menyelesaikan proses produksi.
- Mengurangi stok bahan baku setelah produksi selesai.
- Menambahkan hasil produksi ke inventori produk jadi.

### 6. Manajemen Sales Order
- Membuat pesanan pelanggan.
- Melakukan reservasi stok produk jadi secara otomatis.
- Memproses pengiriman pesanan.
- Membatalkan pesanan dan mengembalikan stok yang direservasi.
- Memantau status pesanan.

## Teknologi yang Digunakan

| Komponen | Teknologi |
|---|---|
| Backend | Go |
| Web Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Frontend | React |
| Build Tool | Vite |
| Styling | CSS |
| Containerization | Docker |
| Version Control | Git dan GitHub |

## Arsitektur Sistem

Aplikasi menggunakan arsitektur frontend dan backend yang terpisah.

```text
+--------------------------+
|       React Frontend     |
|         Vite + CSS       |
+------------+-------------+
             |
             | HTTP REST API
             |
+------------v-------------+
|        Go Backend        |
|       Gin Framework      |
|         GORM ORM         |
+------------+-------------+
             |
             | SQL
             |
+------------v-------------+
|        PostgreSQL        |
|          Database        |
+--------------------------+
```

## Struktur Proyek

```text
study-case-2/
|
+-- backend/
|   +-- api.go
|   +-- main.go
|   +-- models.go
|   +-- go.mod
|   +-- go.sum
|
+-- frontend/
|   +-- public/
|   +-- src/
|   |   +-- assets/
|   |   +-- App.jsx
|   |   +-- App.css
|   |   +-- index.css
|   |   +-- main.jsx
|   +-- package.json
|   +-- vite.config.js
|
+-- README.md
```

## Persyaratan Sistem

Pastikan perangkat sudah memiliki:

- Go
- Node.js dan npm
- PostgreSQL
- Docker (opsional)
- Git

## Instalasi dan Konfigurasi

### 1. Clone Repository

```bash
git clone https://github.com/Herazor/digitalisasi-manajemen-produksi-inventori.git
cd study-case-2
```

### 2. Konfigurasi Database

Aplikasi menggunakan PostgreSQL sebagai database utama.

Buat database dengan nama:

```sql
CREATE DATABASE production_inventory;
```

Buat file `.env` di dalam folder `backend` dengan konfigurasi berikut:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=admin
DB_PASSWORD=production_demo
DB_NAME=production_inventory
DB_SSLMODE=disable
```

Sesuaikan konfigurasi dengan pengaturan PostgreSQL di perangkat masing-masing.

### 3. Menjalankan Backend

Buka terminal dan jalankan:

```bash
cd backend
go mod download
go run .
```

Backend akan berjalan pada:

```text
http://localhost:9000
```

### 4. Menjalankan Frontend

Buka terminal baru dan jalankan:

```bash
cd frontend
npm install
npm run dev
```

Frontend dapat diakses melalui:

```text
http://localhost:5173
```

## Dokumentasi API

Alamat dasar API:

```text
http://localhost:9000/api
```

### Material

| Metode | Endpoint | Keterangan |
|---|---|---|
| GET | `/materials` | Menampilkan daftar material |
| POST | `/materials` | Menambahkan material |
| PUT | `/materials/:id` | Memperbarui material |
| DELETE | `/materials/:id` | Menghapus material |

### Product

| Metode | Endpoint | Keterangan |
|---|---|---|
| GET | `/products` | Menampilkan daftar produk |
| POST | `/products` | Menambahkan produk |
| PUT | `/products/:id` | Memperbarui produk |
| DELETE | `/products/:id` | Menghapus produk |

### BOM

| Metode | Endpoint | Keterangan |
|---|---|---|
| POST | `/products/:id/boms` | Membuat BOM produk |
| POST | `/products/:id/bom-explosion` | Menghitung kebutuhan material |

### Inventory

| Metode | Endpoint | Keterangan |
|---|---|---|
| GET | `/inventory/materials` | Menampilkan inventori material |
| GET | `/inventory/products` | Menampilkan inventori produk |
| POST | `/materials/:id/stock` | Menambahkan stok material |

### Work Order

| Metode | Endpoint | Keterangan |
|---|---|---|
| GET | `/work-orders` | Menampilkan Work Order |
| POST | `/work-orders` | Membuat Work Order |
| POST | `/work-orders/:id/complete` | Menyelesaikan Work Order |

### Sales Order

| Metode | Endpoint | Keterangan |
|---|---|---|
| GET | `/sales-orders` | Menampilkan Sales Order |
| POST | `/sales-orders` | Membuat Sales Order |
| POST | `/sales-orders/:id/ship` | Memproses pengiriman pesanan |
| POST | `/sales-orders/:id/cancel` | Membatalkan pesanan |

## Logika Manajemen Inventori

Sistem menggunakan tiga komponen utama untuk memantau persediaan:

- **On Hand:** Jumlah stok fisik yang tersedia di gudang.
- **Reserved:** Jumlah stok yang sudah dialokasikan untuk pesanan.
- **Available:** Jumlah stok yang masih dapat digunakan atau dialokasikan.

Rumus perhitungan:

```text
Available = On Hand - Reserved
```

Proses reservasi dan pengurangan stok menggunakan transaksi database untuk menjaga konsistensi data.

## Alur Proses Sistem

### Alur Produksi (Work Order)

```text
Membuat Work Order
        |
        v
Menghitung Kebutuhan Material
Berdasarkan BOM
        |
        v
Reservasi Bahan Baku
        |
        v
Menyelesaikan Produksi
        |
        v
Mengurangi Stok Bahan Baku
        |
        v
Menambahkan Stok Produk Jadi
```

### Alur Penjualan (Sales Order)

```text
Membuat Sales Order
        |
        v
Reservasi Stok Produk Jadi
        |
        v
    +---+---+
    |       |
    v       v
   Ship   Cancel
    |       |
    v       v
Mengurangi Mengembalikan
   Stok    Reservasi
    |       |
    v       v
 SHIPPED CANCELLED
```

## Pengujian

Pengujian dilakukan melalui permintaan API backend dan interaksi langsung melalui frontend.

Skenario yang telah diuji meliputi:

- Pengelolaan material dan produk.
- Pembuatan BOM dan BOM Explosion.
- Penambahan stok material.
- Pembuatan dan penyelesaian Work Order.
- Reservasi bahan baku.
- Pembuatan Sales Order.
- Reservasi stok produk jadi.
- Proses pengiriman pesanan.
- Pembatalan pesanan.
- Pembaruan inventori setelah pengiriman dan pembatalan.

## Pengembangan Selanjutnya

Beberapa pengembangan yang dapat dilakukan:

- Notifikasi perubahan inventori secara real-time.
- Integrasi dengan berbagai kanal penjualan.
- Riwayat pergerakan stok.
- Laporan dan analisis produksi.
- Pengaturan hak akses pengguna.
- Pengujian otomatis.

## Repository GitHub

https://github.com/Herazor/digitalisasi-manajemen-produksi-inventori

---

**Nama Proyek:** Digitalisasi Manajemen Produksi & Inventori

**Fokus:** Manajemen Produksi, Manajemen Inventori, Bill of Materials (BOM), dan Pemrosesan Pesanan Omnichannel.