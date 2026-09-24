🔥 VecScan by Ollie
Advanced Network Reconnaissance & Port Scanning Tool
Go 1.21+ Linux Security Penetration Testing

██╗   ██╗███████╗ ██████╗███████╗ ██████╗ █████╗ ███╗   ██╗
██║   ██║██╔════╝██╔════╝██╔════╝██╔════╝██╔══██╗████╗  ██║
██║   ██║█████╗  ██║     ███████╗██║     ███████║██╔██╗ ██║
╚██╗ ██╔╝██╔══╝  ██║     ╚════██║██║     ██╔══██║██║╚██╗██║
 ╚████╔╝ ███████╗╚██████╗███████║╚██████╗██║  ██║██║ ╚████║
  ╚═══╝  ╚══════╝ ╚═════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝
                    BY OLLIE v1.0.0
📋 Table of Contents
Features
Installation
Usage
Scan Types
Firewall Evasion
Examples
Output Formats
GitHub Setup
✨ Features
🔍 Scanning Capabilities
Feature	Description	Status
TCP SYN Scan (-sS)	Stealth half-open scanning	✅ Ready
TCP Connect Scan (-sT)	Full TCP connection	✅ Ready
TCP ACK Scan (-sA)	Firewall rule detection	✅ Ready
TCP FIN Scan (-sF)	Stealth FIN probe	✅ Ready
TCP NULL Scan (-sN)	No flag probe	✅ Ready
TCP Xmas Scan (-sX)	FIN+PSH+URG probe	✅ Ready
UDP Scan (-sU)	UDP port scanning	✅ Ready
ICMP Scan	Host discovery	✅ Ready
🛡️ Firewall Evasion
Technique	Flag	Description
Packet Fragmentation	--fragment	Split packets to evade IDS/IPS
Custom MTU	--mtu	Set fragment size (default: 16)
Source Port Spoofing	--source-port	Spoof source port (e.g., 53, 80)
IP Spoofing	--spoof	Spoof source IP address
Bad Checksum	--badsum	Send packets with invalid checksum
Rate Limiting	--rate	Control packets per second
⚡ Performance
Concurrency: Configurable goroutine pools
Rate Limiting: Token bucket algorithm
Timeout Control: Per-port timeout configuration
Retry Logic: Automatic retry for filtered ports
📊 Output
Terminal: Beautiful colored output with tables
JSON: Structured export for automation
Graceful Shutdown: Partial reports on interrupt
🚀 Installation
Prerequisites
Go 1.21 or higher
Linux (Kali, Ubuntu, Debian)
Root privileges (for SYN/ACK/FIN/NULL/Xmas/UDP scans)
Method 1: Quick Install (Recommended)
# Clone the repository
git clone https://github.com/YOUR_USERNAME/vecscan.git
cd vecscan

# Run install script
chmod +x scripts/install.sh
sudo ./scripts/install.sh
Method 2: Manual Build
# Clone repository
git clone https://github.com/YOUR_USERNAME/vecscan.git
cd vecscan

# Download dependencies
go mod download

# Build
make build

# Install globally
sudo make install
Method 3: Docker
# Build image
docker build -t vecscan .

# Run
docker run --rm --network host vecscan 192.168.1.1
✅ Verify Installation:
vecscan --version
vecscan --help
📖 Usage
Basic Syntax
vecscan [target] [flags]
Global Flags
Flag	Short	Default	Description
--syn	-sS	true	TCP SYN scan
--connect		false	TCP Connect scan
--ack	-sA	false	TCP ACK scan
--fin	-sF	false	TCP FIN scan
--null	-sN	false	TCP NULL scan
--xmas	-sX	false	TCP Xmas scan
--udp	-sU	false	UDP scan
--ports	-p	1-1000	Port range
--rate		1000	Packets/second
--timeout	-t	3	Timeout seconds
--concurrency	-c	100	Goroutines
--output	-o	""	JSON output file
--verbose	-v	false	Verbose output
--fragment		false	Fragment packets
--mtu		16	Fragment MTU
--source-port		0	Source port
--spoof		""	Spoof source IP
--scan-filtered		false	Auto-scan filtered
🔬 Scan Types Explained
TCP SYN Scan (-sS) - Default
Half-open scanning that doesn't complete the TCP handshake. Fast and stealthy.

sudo vecscan 192.168.1.1 -sS -p 1-1000
TCP Connect Scan (-sT)
Full TCP connection. Works without root but is noisier and slower.

vecscan 192.168.1.1 --connect -p 80,443,8080
TCP ACK Scan (-sA)
Used to map firewall rulesets. Determines if ports are filtered.

sudo vecscan 192.168.1.1 -sA -p 1-65535
Stealth Scans (FIN/NULL/Xmas)
Bypasses some firewalls and IDS systems that don't handle unusual flag combinations.

# FIN scan
sudo vecscan 192.168.1.1 -sF -p 1-1000

# NULL scan
sudo vecscan 192.168.1.1 -sN -p 22,80,443

# Xmas scan
sudo vecscan 192.168.1.1 -sX -p 1-500
UDP Scan (-sU)
Scans UDP ports. Requires protocol-specific probes for accuracy.

sudo vecscan 192.168.1.1 -sU -p 53,67,68,123,161
🛡️ Firewall Evasion Techniques
1. Packet Fragmentation
Split packets into small fragments to bypass packet inspection.

# Fragment with default MTU (16 bytes)
sudo vecscan 192.168.1.1 --fragment -p 1-1000

# Custom MTU size
sudo vecscan 192.168.1.1 --fragment --mtu 8 -p 1-1000
2. Source Port Manipulation
Spoof common service ports to bypass "allow" rules.

# Spoof DNS source port
sudo vecscan 192.168.1.1 --source-port 53 -p 1-1000

# Spoof HTTP source port
sudo vecscan 192.168.1.1 --source-port 80 -p 1-1000
3. Combined Evasion
# Maximum evasion mode
sudo vecscan 192.168.1.1 \
  --fragment --mtu 8 \
  --source-port 53 \
  --badsum \
  -sF \
  -p 1-65535 \
  --rate 100
⚠️ Warning: Some evasion techniques require root privileges and may trigger security alerts. Use responsibly and only on authorized systems.
💡 Usage Examples
Example 1: Basic Network Scan
# Scan top 1000 ports on a single host
vecscan 192.168.1.1

# Scan specific ports
vecscan 192.168.1.1 -p 22,80,443,3306,8080
Example 2: Full Port Scan
# Scan all 65535 ports
sudo vecscan 192.168.1.1 --all-ports

# Or specify range
sudo vecscan 192.168.1.1 -p 1-65535
Example 3: Network Range Scan
# CIDR notation
sudo vecscan 192.168.1.0/24 -p 22,80,443

# Multiple specific hosts
sudo vecscan 192.168.1.1,192.168.1.10,192.168.1.20 -p 1-1000
Example 4: Stealth Scan with Evasion
# Slow, stealthy scan with fragmentation
sudo vecscan 192.168.1.1 \
  -sF \
  --fragment \
  --source-port 53 \
  -p 1-1000 \
  --rate 100 \
  -t 5
Example 5: UDP Service Discovery
# Scan common UDP services
sudo vecscan 192.168.1.1 -sU -p 53,67,68,69,123,161,162,500,514,1900
Example 6: Export Results
# JSON output for further processing
sudo vecscan 192.168.1.1 -p 1-65535 -o scan-results.json

# Pretty print JSON
cat scan-results.json | jq
Example 7: Filtered Port Handling
# Scan with automatic filtered port scanning
sudo vecscan 192.168.1.1 --scan-filtered -p 1-1000

# Or let it prompt you
sudo vecscan 192.168.1.1 -p 1-1000
# [?] Scan filtered ports with evasion techniques? [Y/n]: Y
Example 8: High Performance Scan
# Fast scan with high concurrency
sudo vecscan 192.168.1.1 \
  -p 1-65535 \
  --rate 10000 \
  -c 500 \
  -t 2
📊 Output Formats
Terminal Output
═══════════════════════════════════════════════════════════
                     SCAN SUMMARY
═══════════════════════════════════════════════════════════
Target:      192.168.1.1
Scan Type:   TCP SYN
Duration:    15.234s
Start Time:  2024-01-15 14:30:22
End Time:    2024-01-15 14:30:37

PORT STATISTICS:
  Open:     12
  Closed:   988
  Filtered: 0
  Total:    1000

OPEN PORTS:
PORT     PROTOCOL   STATE        SERVICE              RESPONSE
─────────────────────────────────────────────────────────
22       tcp        open         ssh                  2.341ms
80       tcp        open         http                 1.892ms
443      tcp        open         https                2.156ms
3306     tcp        open         mysql                3.421ms
8080     tcp        open         http-proxy           2.098ms
JSON Output Structure
{
  "target": "192.168.1.1",
  "scan_type": "TCP SYN",
  "start_time": "2024-01-15T14:30:22Z",
  "end_time": "2024-01-15T14:30:37Z",
  "hosts": [
    {
      "ip": "192.168.1.1",
      "hostname": "router.local",
      "status": "up",
      "open_ports": [
        {
          "port": 22,
          "protocol": "tcp",
          "state": "open",
          "service": "ssh",
          "reason": "syn-ack",
          "response_time": "2.341ms"
        }
      ]
    }
  ],
  "statistics": {
    "total_ports": 1000,
    "open_count": 12,
    "closed_count": 988,
    "filtered_count": 0
  }
}
🐙 GitHub Repository Setup
Step 1: Create Repository
# Create new repository on GitHub
# Name: vecscan
# Description: Advanced Network Scanner by Vectalith Labs
Step 2: Initialize and Push
# In your vecscan directory
git init
git add .
git commit -m "Initial commit: VecScan v1.0.0"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/vecscan.git
git push -u origin main
Step 3: Clone on Kali Linux
# On your Kali machine
sudo apt update
sudo apt install golang git make

# Clone repository
cd /opt
sudo git clone https://github.com/YOUR_USERNAME/vecscan.git

# Build and install
cd vecscan
make build
sudo make install

# Verify
vecscan --version
Step 4: Create Global Command
The make install command already installs to /usr/local/bin. To make it work everywhere, ensure /usr/local/bin is in your PATH:

echo $PATH | grep /usr/local/bin

# If not present, add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/usr/local/bin
💡 Pro Tip: Create a symlink for development:
sudo ln -sf /opt/vecscan/build/vecscan /usr/local/bin/vecscan
This allows you to rebuild without re-installing.
🔧 Troubleshooting
Issue	Solution
"permission denied"	Run with sudo for raw socket scans
"command not found"	Check PATH or reinstall with make install
"no such host"	Check DNS resolution or use IP address
"too many open files"	Increase ulimit: ulimit -n 65535
All ports filtered	Target may be down or blocking all probes
📜 License
MIT License - See LICENSE file for details.

⚠️ Legal Disclaimer: This tool is for authorized security testing only. Unauthorized scanning of networks you do not own or have explicit permission to test is illegal. The authors assume no liability for misuse or damage caused by this program.
🤝 Contributing
Fork the repository
Create feature branch: git checkout -b feature/amazing-feature
Commit changes: git commit -m 'Add amazing feature'
Push to branch: git push origin feature/amazing-feature
Open Pull Request
📧 Contact
Vectalith Labs

GitHub: @Ollie33-a
Issues: GitHub Issues
Built with 🔥 by Ollie
Advanced Network Reconnaissance
