# 📧 Bot Email Temporary Telegram

Bot Telegram untuk bikin dan terima email sementara. Tinggal chat bot di Telegram, langsung dapet email baru!

## ✨ Fitur

- ⚡ Bikin email baru secepat kilat
- 📨 Terima email langsung di Telegram
- 📎 Support file lampiran
- 🔄 Ganti email kapan aja
- 👥 Bisa dipake rame-rame

## 📱 Cara Pake

Gampang banget pakainya! Ada 4 perintah aja:

- `/start` - Mulai bot 
- `/newmail` - Minta email baru
- `/mymail` - Cek email yang lagi dipake
- `/help` - Bantuan

## 🔧 Apa Aja Yang Dibutuhin

Kalo mau pasang sendiri, siapin dulu:
- Node.js
- Domain
- Akun Cloudflare
- Bot Telegram
- Server

## ⚙️ Cara Setting

1. Bikin file `.env` isinya:
```
BOT_TOKEN=token_bot_telegram_kamu
DOMAIN=domain.kamu.com
SMTP_PORT=25
CLOUDFLARE_EMAIL=email_kamu
CLOUDFLARE_API_TOKEN=token_kamu
CLOUDFLARE_ZONE_ID=zone_id_kamu
CLOUDFLARE_ACCOUNT_ID=account_id_kamu
SERVER_IP=ip_server_kamu
```

2. Install yang dibutuhin:
```bash
npm install
```

3. Jalanin botnya:
```bash
npm start
```

## 📝 Cara Kerja

1. Kamu chat bot minta email baru
2. Bot langsung bikinin email untukmu
3. Kalo ada yang kirim email, langsung masuk ke Telegram
4. Mau ganti email? Tinggal minta lagi!

## ⚠️ Penting Nih!

- Ini cuma buat email sementara ya
- Jangan dipake buat yang penting-penting
- Jangan buat simpen data rahasia

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
- node-telegram-bot-api
- smtp-server
- mailparser
- axios
- dotenv

---

Dibuat dengan ❤️ buat yang butuh email dadakan!