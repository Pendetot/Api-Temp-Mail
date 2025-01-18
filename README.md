# 📧 Bot Email Temporary Telegram

Bot Telegram untuk bikin dan terima email sementara. Tinggal chat bot di Telegram, langsung dapet email baru!

## ✨ Fitur

- ⚡ Bikin email baru secepat kilat
- 📨 Terima email langsung di Telegram
- 📎 Support file lampiran
- 🔄 Ganti email kapan aja
- 👥 Bisa dipake rame-rame
- 🔍 Auto deteksi IP server
- 🛡️ Setup DNS otomatis

## 📱 Cara Pake

Gampang banget pakainya! Ada 4 perintah aja:
- `/start` - Mulai bot 
- `/newmail` - Minta email baru
- `/mymail` - Cek email yang lagi dipake
- `/help` - Bantuan

## 🔧 Apa Aja Yang Dibutuhin

Kalo mau pasang sendiri, siapin dulu:
- Node.js
- Domain yang udah terdaftar di Cloudflare
- Akun Cloudflare
- Bot Telegram (bikin di @BotFather)
- VPS/Server (yang support port 25)

## ⚙️ Cara Setting

1. Clone repository ini
```bash
git clone https://github.com/Pendetot/Tele-Temp-Mail
```

2. Install yang dibutuhin:
```bash
npm install
```

3. Jalanin setup:
```bash
npm start
```

4. Ikutin petunjuk setup untuk masukin:
- Token Bot Telegram
- Domain
- Port SMTP (default: 25)
- Email Cloudflare
- API Token Cloudflare
- Zone ID Cloudflare
- Account ID Cloudflare

## 📝 Cara Kerja

1. Setup Otomatis:
   - Deteksi IP server otomatis
   - Konfigurasi DNS di Cloudflare
   - Setup MX dan SPF records

2. Penggunaan:
   - Chat bot untuk minta email baru
   - Bot langsung bikinin email untukmu
   - Terima email langsung di Telegram
   - Support lampiran file
   - Ganti email kapan aja

## ⚠️ Penting Nih!

- Ini cuma buat email sementara ya
- Jangan dipake buat yang penting-penting
- Jangan buat simpen data rahasia
- Pastikan server support SMTP port 25
- Backup file .env kalo udah dibuat

## 🗂️ Struktur File

- `setup.js` - Untuk setup awal dan konfigurasi
- `index.js` - File utama bot dan server SMTP
- `ip.js` - Untuk deteksi IP server
- `cloudflare.js` - Manajemen DNS Cloudflare

## 🤝 Mau Bantuin Ngoding?

Kalo kamu programmer dan mau bantu ngembangkan:
1. Fork repo ini
2. Bikin branch baru
3. Coding deh!
4. Kirim Pull Request

## 💡 Butuh Bantuan?

Kalo ada masalah:
- Buka issue di GitHub
- Jelasin masalahnya
- Kita bakal bantuin secepatnya!

## 🙏 Thanks To

Project ini pake bantuan dari library:
- node-telegram-bot-api - Untuk bot Telegram
- smtp-server - Untuk terima email
- mailparser - Untuk parse email
- axios - Untuk API requests
- dotenv - Untuk environment variables
- crypto - Untuk generate random email

---
Dibuat dengan ❤️ buat yang butuh email dadakan!