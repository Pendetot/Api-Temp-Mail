require('dotenv').config();
const TelegramBot = require('node-telegram-bot-api');
const SMTPServer = require('smtp-server').SMTPServer;
const simpleParser = require('mailparser').simpleParser;
const crypto = require('crypto');
const CloudflareManager = require('./cloudflare');

// Validasi environment variables
const requiredEnvVars = [
    'BOT_TOKEN',
    'DOMAIN',
    'SMTP_PORT',
    'CLOUDFLARE_EMAIL',
    'CLOUDFLARE_API_TOKEN',
    'CLOUDFLARE_ZONE_ID',
    'CLOUDFLARE_ACCOUNT_ID',
    'SERVER_IP'
];

for (const envVar of requiredEnvVars) {
    if (!process.env[envVar]) {
        console.error(`Error: Environment variable ${envVar} tidak ditemukan`);
        process.exit(1);
    }
}

const emailPengguna = new Map();
const emailKeChat = new Map();
const cloudflare = new CloudflareManager();

async function setupDomain() {
    try {
        console.log('Memulai setup domain di Cloudflare...');
        await cloudflare.setupDNSRecords(process.env.DOMAIN, process.env.SERVER_IP);
        
        console.log('Menunggu propagasi DNS...');
        let configured = false;
        let attempts = 0;
        
        while (!configured && attempts < 10) {
            const status = await cloudflare.checkDNSPropagation(process.env.DOMAIN);
            if (status.mxConfigured && status.spfConfigured) {
                configured = true;
                console.log('DNS telah terkonfigurasi dengan benar');
            } else {
                attempts++;
                console.log(`Menunggu propagasi DNS... Percobaan ${attempts}/10`);
                await new Promise(resolve => setTimeout(resolve, 30000));
            }
        }
        
        if (!configured) {
            console.log('Peringatan: DNS belum sepenuhnya terpropagasi');
        }
    } catch (error) {
        console.error('Error dalam setup domain:', error);
        throw error;
    }
}

const botOptions = {
    polling: {
        interval: 300,
        autoStart: true,
        params: {
            timeout: 10
        }
    },
    request: {
        proxy: undefined,
        timeout: 60000
    }
};

const bot = new TelegramBot(process.env.BOT_TOKEN, botOptions);

let pollingErrorCount = 0;
const MAX_POLLING_ERRORS = 5;

bot.on('polling_error', async (error) => {
    console.error('Bot polling error:', error.code || error.message);
    pollingErrorCount++;

    if (pollingErrorCount >= MAX_POLLING_ERRORS) {
        console.log('Terlalu banyak error polling, mencoba restart bot...');
        try {
            await bot.stopPolling();
            await new Promise(resolve => setTimeout(resolve, 10000));
            await bot.startPolling();
            pollingErrorCount = 0;
            console.log('Bot berhasil direstart');
        } catch (restartError) {
            console.error('Gagal restart bot:', restartError);
            process.exit(1);
        }
    } else {
        await new Promise(resolve => setTimeout(resolve, 5000));
    }
});

async function sendTelegramMessage(chatId, message, options = {}) {
    let retries = 3;
    while (retries > 0) {
        try {
            const defaultOptions = {
                disable_web_page_preview: true
            };
            const mergedOptions = { ...defaultOptions, ...options };
            await bot.sendMessage(chatId, message, mergedOptions);
            return;
        } catch (error) {
            console.error(`Error mengirim pesan (${retries} percobaan tersisa):`, error.code || error.message);
            retries--;
            if (retries > 0) {
                if (error.code === 'ETELEGRAM' && options.parse_mode) {
                    delete options.parse_mode;
                    message = message.replace(/[\*\_\[\]\(\)\~\`\>\#\+\-\=\|\{\}\.\\]/g, '');
                }
                await new Promise(resolve => setTimeout(resolve, 1000));
            } else {
                throw error;
            }
        }
    }
}

function buatEmail() {
    const acak = crypto.randomBytes(8).toString('hex');
    return `${acak}@${process.env.DOMAIN}`;
}

const serverSMTP = new SMTPServer({
    secure: false,
    authOptional: true,
    disabledCommands: ['AUTH'],
    size: 25 * 1024 * 1024,
    onConnect(session, callback) {
        try {
            console.log('Koneksi SMTP baru dari:', session.remoteAddress);
            callback();
        } catch (error) {
            console.error('Error pada koneksi SMTP:', error);
            callback();
        }
    },
    onData(stream, session, callback) {
        let chunks = [];

        stream.on('data', (chunk) => {
            chunks.push(chunk);
        });

        stream.on('end', async () => {
            try {
                const buffer = Buffer.concat(chunks);
                const parsed = await simpleParser(buffer);
                const emailTujuan = parsed.to.text.toLowerCase();
                const chatId = emailKeChat.get(emailTujuan);
                
                if (chatId) {
                    const formatSender = (from) => {
                        if (from.value && from.value[0]) {
                            return from.value[0].name || from.value[0].address || 'Tidak diketahui';
                        }
                        return from.text || 'Tidak diketahui';
                    };

                    const pengirim = formatSender(parsed.from);
                    const subjek = parsed.subject || '(Tidak ada subjek)';
                    const isiPesan = parsed.text || '(Tidak ada isi pesan)';

                    let pesanTeks = `📬 Email Baru Diterima!\n\n`;
                    pesanTeks += `👤 Dari: ${pengirim}\n`;
                    pesanTeks += `📌 Subjek: ${subjek}\n`;
                    pesanTeks += `🕐 Waktu: ${parsed.date.toLocaleString('id-ID')}\n\n`;
                    pesanTeks += `📝 Isi Pesan:\n${isiPesan}`;

                    try {
                        await sendTelegramMessage(chatId, pesanTeks);

                        if (parsed.attachments && parsed.attachments.length > 0) {
                            await sendTelegramMessage(chatId, 
                                `📎 Wah, ada ${parsed.attachments.length} file terlampir nih! Tunggu sebentar ya...`
                            );
                            
                            for (const lampiran of parsed.attachments) {
                                try {
                                    await bot.sendDocument(chatId, lampiran.content, {
                                        filename: lampiran.filename,
                                        caption: `📎 File: ${lampiran.filename}`
                                    });
                                } catch (error) {
                                    console.error('Error mengirim lampiran:', error);
                                    await sendTelegramMessage(chatId, 
                                        `❌ Maaf, gagal mengirim file "${lampiran.filename}". File mungkin terlalu besar atau tidak didukung.`
                                    );
                                }
                            }
                        }
                    } catch (error) {
                        console.error('Error mengirim pesan:', error);
                    }
                } else {
                    console.log('Email diterima untuk alamat yang tidak terdaftar:', emailTujuan);
                }
            } catch (error) {
                console.error('Error memproses email:', error);
            }
            callback();
        });

        stream.on('error', (error) => {
            console.error('Error pada stream email:', error);
            callback(error);
        });
    },
    onError(error) {
        console.error('SMTP Error:', error);
    }
});

bot.onText(/\/start/, async (msg) => {
    const chatId = msg.chat.id;
    const nama = msg.from.first_name || 'teman';
    
    try {
        await sendTelegramMessage(chatId, 
            `Hai ${nama}! 👋\n\n` +
            `Selamat datang di layanan Email Temporary! Aku bakal bantu kamu buat bikin email sementara yang bisa kamu pakai dimana aja.\n\n` +
            `Ini nih perintah yang bisa kamu gunakan:\n\n` +
            `📧 /newmail - Bikin email baru\n` +
            `📨 /mymail - Lihat email kamu sekarang\n` +
            `❓ /help - Kalau kamu butuh bantuan\n\n` +
            `Mau langsung bikin email? Ketik /newmail aja! 😊`
        );
    } catch (error) {
        console.error('Error di command start:', error);
    }
});

bot.onText(/\/newmail/, async (msg) => {
    const chatId = msg.chat.id;
    const nama = msg.from.first_name || 'teman';
    
    try {
        const oldEmail = emailPengguna.get(chatId);
        if (oldEmail) {
            emailKeChat.delete(oldEmail);
        }
        
        const email = buatEmail();
        emailPengguna.set(chatId, email);
        emailKeChat.set(email, chatId);
        
        await sendTelegramMessage(chatId, 
            `✨ Siap ${nama}! Aku udah bikin email baru buat kamu:\n\n` +
            `📧 ${email}\n\n` +
            `Email ini bisa langsung kamu pakai. Tenang aja, nanti kalau ada email masuk, aku langsung kabarin kamu disini ya! 😉\n\n` +
            `Oh iya, kalau butuh email baru lagi, tinggal ketik /newmail aja ya!`
        );
    } catch (error) {
        console.error('Error membuat email baru:', error);
        await sendTelegramMessage(chatId,
            `Maaf ${nama}, ada masalah saat membuat email baru. Coba lagi dalam beberapa saat ya!`
        );
    }
});

bot.onText(/\/mymail/, async (msg) => {
    const chatId = msg.chat.id;
    const nama = msg.from.first_name || 'teman';
    
    try {
        const email = emailPengguna.get(chatId);
        
        if (email) {
            await sendTelegramMessage(chatId, 
                `Hai ${nama}! 👋\n\n` +
                `Ini email kamu yang aktif sekarang:\n\n` +
                `📧 ${email}\n\n` +
                `Email ini masih aktif dan siap dipakai ya! 😊`
            );
        } else {
            await sendTelegramMessage(chatId, 
                `Hai ${nama}! Sepertinya kamu belum punya email nih... 🤔\n\n` +
                `Mau bikin email baru? Gampang kok!\n` +
                `Tinggal ketik /newmail aja ya! 😊`
            );
        }
    } catch (error) {
        console.error('Error mengecek email:', error);
        await sendTelegramMessage(chatId,
            `Maaf ${nama}, ada masalah saat mengecek email kamu. Coba lagi dalam beberapa saat ya!`
        );
    }
});

bot.onText(/\/help/, async (msg) => {
    const chatId = msg.chat.id;
    const nama = msg.from.first_name || 'teman';
    
    try {
        await sendTelegramMessage(chatId,
            `Hai ${nama}! 👋\n\n` +
            `Tenang, aku disini buat bantu kamu! 😊\n\n` +
            `🎯 Ini nih yang bisa kamu lakukan:\n\n` +
            `📧 /newmail - Bikin email baru\n` +
            `📨 /mymail - Lihat email kamu sekarang\n` +
            `❓ /help - Buat lihat bantuan ini lagi\n\n` +
            `📝 Cara pakainya gampang banget:\n` +
            `1. Ketik /newmail buat bikin email\n` +
            `2. Pakai email itu dimana aja yang kamu mau\n` +
            `3. Nanti kalau ada email masuk, aku langsung kabarin kamu disini\n` +
            `4. Mau ganti email? Ketik /newmail lagi aja!\n\n` +
            `Ada yang masih bingung? Jangan ragu buat tanya ke aku ya! 😊`
        );
    } catch (error) {
        console.error('Error mengirim pesan bantuan:', error);
        await sendTelegramMessage(chatId,
            `Maaf ${nama}, ada masalah saat menampilkan bantuan. Coba lagi dalam beberapa saat ya!`
        );
    }
});

async function startServer() {
    try {
        console.log('Memulai server...');
        await setupDomain();
        
        serverSMTP.listen(process.env.SMTP_PORT, () => {
            console.log(`Server SMTP berjalan di port ${process.env.SMTP_PORT}`);
        });
        
        console.log('Bot sedang berjalan...');
    } catch (error) {
        console.error('Error saat memulai server:', error);
        process.exit(1);
    }
}

process.on('SIGTERM', () => {
    console.log('Menerima SIGTERM. Membersihkan...');
    cleanup();
});

process.on('SIGINT', () => {
    console.log('Menerima SIGINT. Membersihkan...');
    cleanup();
});

async function cleanup() {
    console.log('Membersihkan dan menutup koneksi...');
    try {
        await bot.stopPolling();
        console.log('Bot polling dihentikan');
        
        await new Promise((resolve) => {
            serverSMTP.close(() => {
                console.log('Server SMTP ditutup');
                resolve();
            });
        });
        
        console.log('Cleanup selesai');
        process.exit(0);
    } catch (error) {
        console.error('Error saat cleanup:', error);
        process.exit(1);
    }
}

process.on('uncaughtException', (error) => {
    console.error('Uncaught Exception:', error);
    cleanup();
});

process.on('unhandledRejection', (error) => {
    console.error('Unhandled Rejection:', error);
    cleanup();
});

startServer();