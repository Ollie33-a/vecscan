package utils

import (
    "fmt"
    "net"
    "strconv"
    "strings"
)

// Common services database
var commonServices = map[int]string{
    21:    "ftp",
    22:    "ssh",
    23:    "telnet",
    25:    "smtp",
    53:    "domain",
    80:    "http",
    110:   "pop3",
    111:   "rpcbind",
    135:   "msrpc",
    139:   "netbios-ssn",
    143:   "imap",
    443:   "https",
    445:   "microsoft-ds",
    993:   "imaps",
    995:   "pop3s",
    1723:  "pptp",
    3306:  "mysql",
    3389:  "ms-wbt-server",
    5432:  "postgresql",
    5900:  "vnc",
    8080:  "http-proxy",
    8443:  "https-alt",
    9200:  "elasticsearch",
    27017: "mongodb",
}

// UDP services
var udpServices = map[int]string{
    53:    "dns",
    67:    "dhcp",
    68:    "dhcp",
    69:    "tftp",
    123:   "ntp",
    161:   "snmp",
    162:   "snmptrap",
    500:   "isakmp",
    514:   "syslog",
    520:   "rip",
    1900:  "upnp",
    5353:  "mdns",
}

// ParsePortRange parses port range strings like "1-1000,8080,9000-9100"
func ParsePortRange(portStr string) ([]int, error) {
    var ports []int
    
    // Handle special cases
    if portStr == "-" || portStr == "all" {
        for i := 1; i <= 65535; i++ {
            ports = append(ports, i)
        }
        return ports, nil
    }
    
    // Split by comma
    ranges := strings.Split(portStr, ",")
    
    for _, r := range ranges {
        r = strings.TrimSpace(r)
        
        // Check if it's a range
        if strings.Contains(r, "-") {
            parts := strings.Split(r, "-")
            if len(parts) != 2 {
                return nil, fmt.Errorf("invalid port range: %s", r)
            }
            
            start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
            if err != nil {
                return nil, err
            }
            
            end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
            if err != nil {
                return nil, err
            }
            
            if start < 1 || end > 65535 || start > end {
                return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
            }
            
            for i := start; i <= end; i++ {
                ports = append(ports, i)
            }
        } else {
            // Single port
            port, err := strconv.Atoi(r)
            if err != nil {
                return nil, err
            }
            
            if port < 1 || port > 65535 {
                return nil, fmt.Errorf("invalid port: %d", port)
            }
            
            ports = append(ports, port)
        }
    }
    
    return ports, nil
}

// ResolveTarget resolves a target string to IP addresses
func ResolveTarget(target string) ([]net.IP, error) {
    var ips []net.IP
    
    // Check if it's CIDR notation
    _, ipnet, err := net.ParseCIDR(target)
    if err == nil {
        // Generate all IPs in CIDR range
        for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
            ipCopy := make(net.IP, len(ip))
            copy(ipCopy, ip)
            ips = append(ips, ipCopy)
        }
        return ips, nil
    }
    
    // Check if it's a single IP
    ip := net.ParseIP(target)
    if ip != nil {
        return []net.IP{ip}, nil
    }
    
    // Try DNS resolution
    addrs, err := net.LookupIP(target)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve %s: %v", target, err)
    }
    
    for _, addr := range addrs {
        if addr.To4() != nil {
            ips = append(ips, addr)
        }
    }
    
    return ips, nil
}

func incrementIP(ip net.IP) {
    for j := len(ip) - 1; j >= 0; j-- {
        ip[j]++
        if ip[j] > 0 {
            break
        }
    }
}

// LookupService returns the service name for a port
func LookupService(port int, proto string) string {
    if proto == "udp" {
        if svc, ok := udpServices[port]; ok {
            return svc
        }
    }
    
    if svc, ok := commonServices[port]; ok {
        return svc
    }
    
    return "unknown"
}

// IsCommonService checks if port is a common service
func IsCommonService(port int) bool {
    _, ok := commonServices[port]
    return ok
}

// GetUDPProbe returns a probe payload for UDP scanning
func GetUDPProbe(port int) []byte {
    // Protocol-specific probes for better accuracy
    probes := map[int][]byte{
        53:   {0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x13, 0x12},
        123:  {0xe3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        161:  {0x30, 0x26, 0x02, 0x01, 0x00, 0x04, 0x06, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63},
        1900: []byte("M-SEARCH * HTTP/1.1\r\nHost: 239.255.255.250:1900\r\n\r\n"),
    }
    
    if probe, ok := probes[port]; ok {
        return probe
    }
    
    // Generic probe
    return []byte{0x00, 0x00, 0x00, 0x00}
}