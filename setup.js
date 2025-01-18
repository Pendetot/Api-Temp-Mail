const fs = require('fs');
const readline = require('readline');
const path = require('path');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const question = (query) => new Promise((resolve) => rl.question(query, resolve));

async function setup() {
    const envPath = path.join(__dirname, '.env');
    
    if (fs.existsSync(envPath)) {
        console.log('File .env sudah ada, melanjutkan program...');
        rl.close();
        require('./index.js');
        return;
    }

    console.log('Selamat datang di setup Telegram Temp Mail!');
    console.log('Silakan masukkan informasi yang diperlukan:\n');

    const configs = {
        BOT_TOKEN: await question('Masukkan Token Bot Telegram: '),
        DOMAIN: await question('Masukkan Domain: '),
        SMTP_PORT: await question('Masukkan Port SMTP (default: 25): ') || '25',
        CLOUDFLARE_EMAIL: await question('Masukkan Email Cloudflare: '),
        CLOUDFLARE_API_TOKEN: await question('Masukkan API Token Cloudflare: '),
        CLOUDFLARE_ZONE_ID: await question('Masukkan Zone ID Cloudflare: '),
        CLOUDFLARE_ACCOUNT_ID: await question('Masukkan Account ID Cloudflare: ')
    };

    let envContent = '';
    for (const [key, value] of Object.entries(configs)) {
        envContent += `${key}=${value}\n`;
    }

    try {
        fs.writeFileSync(envPath, envContent);
        console.log('\nFile .env berhasil dibuat!');
        console.log('Memulai program...\n');
        rl.close();
        require('./index.js');
    } catch (error) {
        console.error('Gagal membuat file .env:', error.message);
        process.exit(1);
    }
}

setup();