# 📧 REST API Temp Mail

REST API untuk layanan email sementara yang dibangun dengan Golang. API ini memungkinkan Anda membuat email sementara, menerima pesan, dan mengelola email melalui HTTP endpoints.

## ✨ Fitur

- 🚀 **REST API** - Interface HTTP yang mudah digunakan
- ⚡ **Performa Tinggi** - Dibangun dengan Golang untuk kecepatan optimal
- 📨 **Terima Email Real-time** - Server SMTP terintegrasi
- 📎 **Support Lampiran** - Mendukung file attachment
- 🔄 **Auto Cleanup** - Email otomatis terhapus setelah expired
- 👥 **Multi-user** - Mendukung banyak pengguna secara bersamaan
- 🔍 **Auto IP Detection** - Deteksi IP server otomatis
- 🛡️ **Setup DNS Otomatis** - Konfigurasi Cloudflare otomatis
- 🐳 **Docker Ready** - Mudah deploy dengan Docker

## 📱 API Endpoints

### 1. Health Check
```http
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "message": "Layanan REST API Email Sementara berjalan normal",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0"
}
```

### 2. Buat Email Sementara
```http
POST /api/v1/email
Content-Type: application/json

{
  "domain": "example.com",
  "ttl": 60
}
```

**Response:**
```json
{
  "email": {
    "id": "uuid-email-id",
    "address": "random123@example.com",
    "domain": "example.com",
    "created_at": "2024-01-01T12:00:00Z",
    "expires_at": "2024-01-01T13:00:00Z"
  },
  "message": "Email sementara berhasil dibuat"
}
```

### 3. Ambil Detail Email
```http
GET /api/v1/email/{id}
```

**Response:**
```json
{
  "id": "uuid-email-id",
  "address": "random123@example.com",
  "domain": "example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "expires_at": "2024-01-01T13:00:00Z",
  "messages": []
}
```

### 4. Ambil Pesan Email
```http
GET /api/v1/email/{id}/messages?page=1&limit=10
```

**Response:**
```json
{
  "messages": [
    {
      "id": "uuid-message-id",
      "email_id": "uuid-email-id",
      "from": "sender@example.com",
      "to": "random123@example.com",
      "subject": "Test Email",
      "body": "Isi pesan email",
      "html": "<p>Isi pesan email</p>",
      "attachments": [],
      "headers": {},
      "received_at": "2024-01-01T12:30:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

### 5. Hapus Email
```http
DELETE /api/v1/email/{id}
```

**Response:**
```json
{
  "message": "Email berhasil dihapus"
}
```

## 🔧 Kebutuhan Sistem

- **Go 1.21+** - Bahasa pemrograman
- **Domain** - Terdaftar di Cloudflare
- **Akun Cloudflare** - Untuk manajemen DNS
- **VPS/Server** - Yang mendukung port 25 (SMTP)

## ⚙️ Instalasi & Setup

### 1. Clone Repository
```bash
git clone https://github.com/Pendetot/tele-temp-mail.git
cd tele-temp-mail
```

### 2. Konfigurasi Environment
```bash
cp .env.example .env
nano .env
```

Isi file `.env` dengan konfigurasi Anda:
```env
DOMAIN=your-domain.com
SMTP_PORT=25
PORT=8080
CLOUDFLARE_EMAIL=your-email@example.com
CLOUDFLARE_API_TOKEN=your-api-token
CLOUDFLARE_ZONE_ID=your-zone-id
CLOUDFLARE_ACCOUNT_ID=your-account-id
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Build & Run
```bash
# Build aplikasi
go build -o temp-mail-api ./cmd/api

# Jalankan aplikasi
./temp-mail-api
```

## 🐳 Deploy dengan Docker

### 1. Build & Run dengan Docker Compose
```bash
docker-compose up -d
```

### 2. Build Manual
```bash
# Build image
docker build -t temp-mail-api .

# Run container
docker run -d \
  --name temp-mail-api \
  -p 8080:8080 \
  -p 25:25 \
  --env-file .env \
  temp-mail-api
```

## 📝 Cara Kerja

1. **Setup Otomatis:**
   - Deteksi IP server otomatis
   - Konfigurasi DNS di Cloudflare (MX, A, SPF records)
   - Setup server SMTP dan HTTP API

2. **Penggunaan API:**
   - Kirim POST request untuk membuat email baru
   - Gunakan GET request untuk mengambil pesan
   - Email otomatis expired sesuai TTL yang ditentukan

3. **Penerimaan Email:**
   - Server SMTP menerima email masuk
   - Email diparsing dan disimpan dalam memori
   - Dapat diakses melalui API endpoints

## 🗂️ Struktur Project

```
├── cmd/api/              # Entry point aplikasi
│   └── main.go
├── internal/             # Kode internal aplikasi
│   ├── config/          # Konfigurasi
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Business logic
├── pkg/                 # Package yang dapat digunakan ulang
│   ├── cloudflare/      # Cloudflare API client
│   └── utils/           # Utility functions
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── go.mod              # Go modules
└── README.md           # Dokumentasi
```

## 🔒 Keamanan

- ⚠️ **Hanya untuk email sementara** - Jangan gunakan untuk data penting
- 🔐 **Tidak ada enkripsi** - Email disimpan dalam plain text
- ⏰ **Auto-expire** - Email otomatis terhapus setelah TTL habis
- 🚫 **Tidak ada autentikasi** - API terbuka untuk semua

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Logs
```bash
# Lihat logs Docker
docker-compose logs -f temp-mail-api

# Lihat logs aplikasi langsung
./temp-mail-api
```

## 🛠️ Development

### Menjalankan dalam Mode Development
```bash
# Install air untuk hot reload
go install github.com/cosmtrek/air@latest

# Jalankan dengan hot reload
air
```

### Testing
```bash
# Run tests
go test ./...

# Test dengan coverage
go test -cover ./...
```

## 🤝 Kontribusi

Kami menerima kontribusi dari developer lain:

1. Fork repository ini
2. Buat branch fitur baru (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📋 TODO

- [ ] Implementasi database untuk persistensi
- [ ] Autentikasi API dengan JWT
- [ ] Rate limiting
- [ ] Webhook notifications
- [ ] Web interface
- [ ] Email templates
- [ ] Metrics dan monitoring

## 💡 Troubleshooting

### Port 25 Blocked
Jika port 25 diblokir oleh provider:
```bash
# Gunakan port alternatif
SMTP_PORT=587
```

### DNS Tidak Terpropagasi
```bash
# Cek DNS records
dig MX your-domain.com
dig TXT your-domain.com
```

### Memory Usage Tinggi
```bash
# Kurangi TTL default email
# Edit di internal/services/email.go
ttl = 30 // 30 menit instead of 60
```

## 📞 Support

Jika mengalami masalah:
- 🐛 **Bug Reports**: Buka issue di GitHub
- 💬 **Diskusi**: Gunakan GitHub Discussions
-  **Instagram**: Dm Aol_Ra

## 📄 Lisensi

Project ini menggunakan lisensi MIT. Lihat file `LICENSE` untuk detail lengkap.

## 🙏 Acknowledgments

Terima kasih kepada:
- **Gin Framework** - HTTP web framework
- **go-smtp** - SMTP server implementation
- **Cloudflare** - DNS management
- **Docker** - Containerization

---

**Dibuat dengan ❤️ menggunakan Golang untuk kebutuhan email sementara yang cepat dan reliable!**
