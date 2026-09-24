package scanner

import (
    "fmt"
    "net"
    "time"

    "github.com/vectalithlabs/vecscan/pkg/utils"
)

// UDPScan performs UDP port scanning
func (s *Scanner) UDPScan(ip net.IP, port int) PortInfo {
    start := time.Now()
    address := fmt.Sprintf("%s:%d", ip.String(), port)
    
    // Create UDP connection
    conn, err := net.DialTimeout("udp", address,
        time.Duration(s.config.Timeout)*time.Second)
    
    if err != nil {
        return PortInfo{
            Port:     port,
            Protocol: "udp",
            State:    StateFiltered,
            Reason:   "host-unreachable",
        }
    }
    defer conn.Close()
    
    // Send probe payload
    payload := utils.GetUDPProbe(port)
    conn.SetWriteDeadline(time.Now().Add(time.Duration(s.config.Timeout) * time.Second))
    _, err = conn.Write(payload)
    
    if err != nil {
        return PortInfo{
            Port:     port,
            Protocol: "udp",
            State:    StateFiltered,
            Reason:   "write-failed",
        }
    }
    
    // Try to read response
    buf := make([]byte, 1024)
    conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.Timeout) * time.Second))
    n, err := conn.Read(buf)
    
    if err != nil {
        // No response could mean open|filtered
        return PortInfo{
            Port:     port,
            Protocol: "udp",
            State:    StateOpenFiltered,
            Service:  utils.LookupService(port, "udp"),
            Reason:   "no-response",
        }
    }
    
    // Got response - port is open
    banner := ""
    if n > 0 {
        banner = string(buf[:n])
    }
    
    return PortInfo{
        Port:         port,
        Protocol:     "udp",
        State:        StateOpen,
        Service:      utils.LookupService(port, "udp"),
        Banner:       banner,
        Reason:       "udp-response",
        ResponseTime: time.Since(start),
    }
}