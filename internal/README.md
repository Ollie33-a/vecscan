# 🔥 VecScan by Vectalith Labs

<p align="center">
  <b>Advanced Network Reconnaissance & Port Scanning Tool</b><br>
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go">
  <img src="https://img.shields.io/badge/Linux-Supported-FCC624?style=flat&logo=linux">
  <img src="https://img.shields.io/badge/License-MIT-green.svg">
  <img src="https://img.shields.io/badge/Security-Pentesting-red">
</p>

<pre align="center">
██╗   ██╗███████╗ ██████╗███████╗ ██████╗ █████╗ ███╗   ██╗
██║   ██║██╔════╝██╔════╝██╔════╝██╔════╝██╔══██╗████╗  ██║
██║   ██║█████╗  ██║     ███████╗██║     ███████║██╔██╗ ██║
╚██╗ ██╔╝██╔══╝  ██║     ╚════██║██║     ██╔══██║██║╚██╗██║
 ╚████╔╝ ███████╗╚██████╗███████║╚██████╗██║  ██║██║ ╚████║
  ╚═══╝  ╚══════╝ ╚═════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝
                    BY VECTALITH LABS v1.0.0
</pre>

---

## 📋 Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Scan Types](#scan-types)
5. [Firewall Evasion](#firewall-evasion)
6. [Examples](#examples)
7. [Output Formats](#output-formats)
8. [GitHub Setup](#github-setup)

---

## ✨ Features

### 🔍 Scanning Capabilities

| Feature                | Description                | Status   |
| ---------------------- | -------------------------- | -------- |
| TCP SYN Scan (-sS)     | Stealth half-open scanning | ✅ Ready |
| TCP Connect Scan (-sT) | Full TCP connection        | ✅ Ready |
| TCP ACK Scan (-sA)     | Firewall rule detection    | ✅ Ready |
| TCP FIN Scan (-sF)     | Stealth FIN probe          | ✅ Ready |
| TCP NULL Scan (-sN)    | No flag probe              | ✅ Ready |
| TCP Xmas Scan (-sX)    | FIN+PSH+URG probe          | ✅ Ready |
| UDP Scan (-sU)         | UDP port scanning          | ✅ Ready |
| ICMP Scan              | Host discovery             | ✅ Ready |

### 🛡️ Firewall Evasion

| Technique            | Flag            | Description                        |
| -------------------- | --------------- | ---------------------------------- |
| Packet Fragmentation | `--fragment`    | Split packets to evade IDS/IPS     |
| Custom MTU           | `--mtu`         | Set fragment size (default: 16)    |
| Source Port Spoofing | `--source-port` | Spoof source port (e.g., 53, 80)   |
| IP Spoofing          | `--spoof`       | Spoof source IP address            |
| Bad Checksum         | `--badsum`      | Send packets with invalid checksum |
| Rate Limiting        | `--rate`        | Control packets per second         |

### ⚡ Performance

- **Concurrency**: Configurable goroutine pools
- **Rate Limiting**: Token bucket algorithm
- **Timeout Control**: Per-port timeout configuration
- **Retry Logic**: Automatic retry for filtered ports

### 📊 Output

- **Terminal**: Beautiful colored output with tables
- **JSON**: Structured export for automation
- **Graceful Shutdown**: Partial reports on interrupt

---

## 🚀 Installation

### Prerequisites

- Go 1.21 or higher
- Linux (Kali, Ubuntu, Debian)
- Root privileges (for SYN/ACK/FIN/NULL/Xmas/UDP scans)

### Method 1: Quick Install (Recommended)

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/vecscan.git
cd vecscan

# Run install script
chmod +x scripts/install.sh
sudo ./scripts/install.sh
```
