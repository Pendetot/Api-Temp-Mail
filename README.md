# 📧 Bot Email Temporary Telegram

```js
console.log(`
=========================
   BOT EMAIL TEMPORARY
   Bikin Email Dadakan
     Lewat Telegram
=========================
`);
```

Hai! 👋 

Kenalin nih bot yang bisa bantuin kamu bikin email sementara lewat Telegram. Gampang banget pakainya, tinggal chat bot di Telegram, langsung dapet email baru deh! 

## ✨ Apa aja yang bisa dilakuin?

```js
// Fitur-fitur keren bot ini
let fitur = [
  "⚡ Bikin email baru dalam hitungan detik",
  "📨 Terima email langsung di Telegram",
  "📎 Bisa terima file lampiran",
  "🔄 Email baru? Tinggal minta lagi!",
  "👥 Bisa dipake bareng-bareng"
];
```

## 📱 Cara Pakainya

Command bot yang bisa kamu pake:
- `/start` - Mulai bot
- `/newmail` - Minta email baru
- `/mymail` - Liat email yang lagi dipake
- `/help` - Bantuan

## 🔧 Mau Pasang Sendiri?

Yang kamu butuhin:
1. Node.js
2. Domain
3. Akun Cloudflare
4. Bot Telegram
5. Server

## ⚙️ Cara Setting

1. Bikin file `.env`, isinya:
```env
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

```js
// Gampang kok cara kerjanya!
let caraPake = [
  "1. Kamu minta email baru ke bot",
  "2. Bot bikinin email untukmu",
  "3. Ada yang kirim email? Langsung masuk ke Telegram",
  "4. Mau email baru? Tinggal minta lagi!"
];
```

## ⚠️ Penting!

```js
console.log(`
INGAT YA!
- Ini cuma buat email sementara
- Jangan buat yang penting-penting
- Jangan buat data rahasia
`);
```

## 🤝 Mau Ikut Ngoding?

Kamu programmer? Mau bantu ngembangkan bot ini? Ini caranya:
1. Fork repo ini
2. Bikin branch baru
3. Coding deh
4. Kirim Pull Request

## 💡 Butuh Bantuan?

```js
let bantuan = {
  cara1: "Buka issue di GitHub",
  cara2: "Jelasin masalahnya",
  janji: "Kita bakal bantuin secepetnya!"
};
```

## 🙏 Thanks To

Bot ini pake bantuan dari:
- node-telegram-bot-api
- smtp-server
- mailparser
- axios
- dotenv

---

Dibuat dengan ❤️ buat semua yang butuh email temporary!