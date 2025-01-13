# 📧 Bot Email Temporary Telegram

<img src="/api/placeholder/800/400" alt="Banner Bot Email Temporary"/>

Halo! Kenalin nih bot keren yang bisa bikin email temporary (email sementara) lewat Telegram. Jadi kamu bisa bikin email dadakan dan terima email langsung di chat Telegram kamu. Mantep kan? 😎

## ✨ Apa aja sih yang bisa dilakuin?

<img src="/api/placeholder/400/300" alt="Fitur Bot"/>

- Bikin email temporary secepat kilat ⚡
- Terima email langsung di Telegram, gak perlu buka-buka app lain 📨
- Support file attachment, jadi bisa terima file juga 📎
- Setup DNS otomatis pakai Cloudflare, jadi gak ribet 🛠️
- Ada sistem anti-error, jadi lebih reliable 💪
- Bisa dipake bareng-bareng sama banyak user 👥

## 🔧 Apa aja yang dibutuhin?

- Node.js v14 ke atas
- Domain yang udah didaftarin di Cloudflare
- Punya akun Cloudflare yang bisa akses API
- Token Bot Telegram
- Server yang punya IP public

## ⚙️ Cara Setting

<img src="/api/placeholder/600/300" alt="Setup Guide"/>

1. Bikin file `.env` di folder utama, isinya kayak gini:

```env
BOT_TOKEN=token_bot_telegram_kamu
DOMAIN=domain.kamu.com
SMTP_PORT=25
CLOUDFLARE_EMAIL=email_cloudflare_kamu
CLOUDFLARE_API_TOKEN=token_api_cloudflare_kamu
CLOUDFLARE_ZONE_ID=id_zone_cloudflare_kamu
CLOUDFLARE_ACCOUNT_ID=id_akun_cloudflare_kamu
SERVER_IP=ip_server_kamu
```

2. Install yang dibutuhin:

```bash
npm install
```

## 🚀 Cara Pakainya

1. Pastiin dulu semua setting di `.env` udah bener yaa
2. Jalanin botnya:

```bash
npm start
```

## 📱 Command Bot yang Bisa Dipake

<img src="/api/placeholder/400/300" alt="Bot Commands"/>

- `/start` - Mulai bot dan liat pesan selamat datang
- `/newmail` - Bikin email baru
- `/mymail` - Cek email yang lagi dipake sekarang
- `/help` - Kalo butuh bantuan

## 🛡️ Fitur Keamanan

Bot ini udah dilengkapin sama:
- Cek kredensial otomatis
- Batasan ukuran email (25MB)
- Sistem anti-error
- Auto-cleaning pas di-shutdown
- Sistem retry kalo gagal kirim pesan

## 🔍 Gimana Cara Kerjanya?

<img src="/api/placeholder/600/400" alt="How It Works"/>

1. Bot bikin email random pas kamu minta
2. Server SMTP nerima email yang masuk
3. Email diproses terus dikirim ke chat Telegram kamu
4. File attachment dikirim sebagai dokumen di Telegram
5. DNS diatur otomatis pake API Cloudflare

## ⚠️ Kalo Error Gimana?

Tenang! Bot ini udah dibikin anti-error:
- Auto-retry kalo gagal kirim pesan
- Auto-recovery kalo polling gagal
- Cleaning yang aman kalo error parah
- Ada log lengkap buat debugging

## 📄 Lisensi

Project ini pake Lisensi MIT - cek [LICENSE](LICENSE) buat lengkapnya.

## 🙏 Thanks To

Bot ini pake library keren ini:
- node-telegram-bot-api
- smtp-server
- mailparser
- axios
- dotenv

## ⚠️ Penting Nih!

Bot ini dibuat buat keperluan temporary aja ya, jadi jangan dipake buat email yang penting-penting atau rahasia. Pake yang bijak ya! 😉

## 💬 Butuh Bantuan?

<img src="/api/placeholder/400/300" alt="Support"/>

Kalo ada masalah atau mau nanya-nanya, langsung aja bikin issue di repo ini ya! Kita bakal bantuin sebisa mungkin 🙌