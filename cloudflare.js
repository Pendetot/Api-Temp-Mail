const axios = require('axios');

class CloudflareManager {
    constructor() {
        this.email = process.env.CLOUDFLARE_EMAIL;
        this.zoneId = process.env.CLOUDFLARE_ZONE_ID;
        this.accountId = process.env.CLOUDFLARE_ACCOUNT_ID;
        this.baseUrl = 'https://api.cloudflare.com/client/v4';
        
        this.headers = {
            'X-Auth-Email': process.env.CLOUDFLARE_EMAIL,
            'X-Auth-Key': process.env.CLOUDFLARE_API_TOKEN,
            'Content-Type': 'application/json'
        };

        console.log('Konfigurasi Cloudflare:');
        console.log('Email:', process.env.CLOUDFLARE_EMAIL);
        console.log('Zone ID:', process.env.CLOUDFLARE_ZONE_ID);
        console.log('Account ID:', process.env.CLOUDFLARE_ACCOUNT_ID);
    }

    async verifyCredentials() {
        try {
            console.log('Memverifikasi kredensial Cloudflare...');
            const response = await axios.get(
                `${this.baseUrl}/zones/${this.zoneId}`,
                { headers: this.headers }
            );
            
            if (response.data.success) {
                console.log('Verifikasi kredensial berhasil');
                return true;
            } else {
                console.error('Verifikasi gagal:', response.data);
                return false;
            }
        } catch (error) {
            console.error('Error verifikasi:', error.response?.data || error.message);
            return false;
        }
    }

    async setupDNSRecords(domain, serverIp) {
        try {
            console.log(`Setup DNS records untuk domain: ${domain}`);
            console.log(`IP Server: ${serverIp}`);

            const isVerified = await this.verifyCredentials();
            if (!isVerified) {
                throw new Error('Gagal verifikasi kredensial Cloudflare');
            }

            const existingRecords = await this.getDNSRecords();
            for (const record of existingRecords) {
                if (['MX', 'A', 'TXT'].includes(record.type)) {
                    await this.deleteDNSRecord(record.id);
                }
            }

            await this.createDNSRecord({
                type: 'MX',
                name: '@',
                content: domain,
                priority: 10,
                proxied: false
            });

            await this.createDNSRecord({
                type: 'A',
                name: '@',
                content: serverIp,
                proxied: false
            });

            await this.createDNSRecord({
                type: 'TXT',
                name: '@',
                content: `v=spf1 ip4:${serverIp} ~all`,
                proxied: false
            });

            console.log('Setup DNS records berhasil');
            return true;
        } catch (error) {
            console.error('Error setup DNS records:', error.message);
            if (error.response) {
                console.error('Response Cloudflare:', error.response.data);
            }
            throw error;
        }
    }

    async createDNSRecord(record) {
        try {
            console.log(`Membuat record ${record.type}...`);
            const response = await axios.post(
                `${this.baseUrl}/zones/${this.zoneId}/dns_records`,
                record,
                { headers: this.headers }
            );
            console.log(`Record ${record.type} berhasil dibuat`);
            return response.data;
        } catch (error) {
            console.error(`Error membuat record ${record.type}:`, error.response?.data || error.message);
            throw error;
        }
    }

    async deleteDNSRecord(recordId) {
        try {
            console.log(`Menghapus record ID: ${recordId}`);
            await axios.delete(
                `${this.baseUrl}/zones/${this.zoneId}/dns_records/${recordId}`,
                { headers: this.headers }
            );
            console.log('Record berhasil dihapus');
        } catch (error) {
            console.error('Error menghapus record:', error.response?.data || error.message);
        }
    }

    async getDNSRecords() {
        try {
            console.log('Mengambil daftar DNS records...');
            const response = await axios.get(
                `${this.baseUrl}/zones/${this.zoneId}/dns_records`,
                { headers: this.headers }
            );
            return response.data.result || [];
        } catch (error) {
            console.error('Error mengambil DNS records:', error.response?.data || error.message);
            throw error;
        }
    }

    async checkDNSPropagation(domain) {
        try {
            const records = await this.getDNSRecords();
            const mxRecord = records.find(r => r.type === 'MX');
            const spfRecord = records.find(r => r.type === 'TXT' && r.content.includes('v=spf1'));
            
            return {
                mxConfigured: !!mxRecord,
                spfConfigured: !!spfRecord
            };
        } catch (error) {
            console.error('Error memeriksa propagasi DNS:', error.message);
            return {
                mxConfigured: false,
                spfConfigured: false
            };
        }
    }
}

module.exports = CloudflareManager;