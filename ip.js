const os = require('os');
const dns = require('dns');
const { promisify } = require('util');

const dnsLookup = promisify(dns.lookup);

async function getServerIP() {
    try {
        const networkInterfaces = os.networkInterfaces();
        
        for (const [interfaceName, interfaces] of Object.entries(networkInterfaces)) {
            if (interfaceName.includes('lo')) continue;
            
            for (const interface of interfaces) {
                if (interface.family === 'IPv4' && !interface.internal) {
                    try {
                        const { address } = await dnsLookup(os.hostname());
                        if (address === interface.address) {
                            return address;
                        }
                    } catch (dnsError) {
                        return interface.address;
                    }
                }
            }
        }
        
        throw new Error('Tidak dapat menemukan IP yang sesuai');
    } catch (error) {
        throw error;
    }
}

module.exports = getServerIP;