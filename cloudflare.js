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
    }

    async verifyCredentials() {
        try {
            const response = await axios.get(
                `${this.baseUrl}/zones/${this.zoneId}`,
                { headers: this.headers }
            );
            
            if (response.data.success) {
                return true;
            } else {
                return false;
            }
        } catch (error) {
            return false;
        }
    }

    async setupDNSRecords(domain, serverIp) {
        try {
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

            return true;
        } catch (error) {
            throw error;
        }
    }

    async createDNSRecord(record) {
        try {
            const response = await axios.post(
                `${this.baseUrl}/zones/${this.zoneId}/dns_records`,
                record,
                { headers: this.headers }
            );
            return response.data;
        } catch (error) {
            throw error;
        }
    }

    async deleteDNSRecord(recordId) {
        try {
            await axios.delete(
                `${this.baseUrl}/zones/${this.zoneId}/dns_records/${recordId}`,
                { headers: this.headers }
            );
        } catch (error) {}
    }

    async getDNSRecords() {
        try {
            const response = await axios.get(
                `${this.baseUrl}/zones/${this.zoneId}/dns_records`,
                { headers: this.headers }
            );
            return response.data.result || [];
        } catch (error) {
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
            return {
                mxConfigured: false,
                spfConfigured: false
            };
        }
    }
}

module.exports = CloudflareManager;