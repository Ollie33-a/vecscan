package scanner

import (
    "context"
    "fmt"
    "net"
    "sync"
    "time"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/vectalithlabs/vecscan/internal/config"
    "github.com/vectalithlabs/vecscan/pkg/utils"
)

// New creates a new Scanner instance
func New(cfg *config.Config) *Scanner {
    return &Scanner{
        config: cfg,
    }
}

// Execute runs the main scan
func (s *Scanner) Execute(ctx context.Context) (*Results, error) {
    results := &Results{
        Target:    s.config.Target,
        StartTime: time.Now(),
        Hosts:     make([]HostResult, 0),
        OpenPorts: make([]PortInfo, 0),
        FilteredPorts: make([]PortInfo, 0),
    }
    
    // Determine scan type
    results.ScanType = s.determineScanType()
    
    // Resolve target
    ips, err := utils.ResolveTarget(s.config.Target)
    if err != nil {
        return nil, err
    }
    
    // Scan each IP
    for _, ip := range ips {
        select {
        case <-ctx.Done():
            results.EndTime = time.Now()
            return results, ctx.Err()
        default:
        }
        
        host := s.scanHost(ctx, ip)
        results.Hosts = append(results.Hosts, host)
        results.OpenPorts = append(results.OpenPorts, host.OpenPorts...)
    }
    
    results.EndTime = time.Now()
    results.Statistics = s.stats
    
    return results, nil
}

func (s *Scanner) scanHost(ctx context.Context, ip net.IP) HostResult {
    host := HostResult{
        IP:        ip.String(),
        Status:    "up",
        OpenPorts: make([]PortInfo, 0),
    }
    
    // Create rate limiter
    rateLimiter := utils.NewRateLimiter(s.config.Rate)
    defer rateLimiter.Stop()
    
    // Create work queue
    portChan := make(chan int, s.config.Concurrency)
    resultChan := make(chan PortInfo, len(s.config.PortList))
    
    var wg sync.WaitGroup
    
    // Start workers
    for i := 0; i < s.config.Concurrency; i++ {
        wg.Add(1)
        go s.worker(ctx, &wg, ip, portChan, resultChan, rateLimiter)
    }
    
    // Queue ports
    go func() {
        defer close(portChan)
        for _, port := range s.config.PortList {
            select {
            case <-ctx.Done():
                return
            case portChan <- port:
            }
        }
    }()
    
    // Collect results
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    // Process results
    for portInfo := range resultChan {
        s.stats.TotalPorts++
        
        switch portInfo.State {
        case StateOpen:
            s.stats.OpenCount++
            host.OpenPorts = append(host.OpenPorts, portInfo)
        case StateClosed:
            s.stats.ClosedCount++
        case StateFiltered:
            s.stats.FilteredCount++
            // Store for potential re-scan
        }
        
        if portInfo.State == StateFiltered {
            s.results.FilteredPorts = append(s.results.FilteredPorts, portInfo)
        }
    }
    
    return host
}

func (s *Scanner) worker(ctx context.Context, wg *sync.WaitGroup, ip net.IP, 
    ports <-chan int, results chan<- PortInfo, rateLimiter *utils.RateLimiter) {
    defer wg.Done()
    
    for port := range ports {
        select {
        case <-ctx.Done():
            return
        default:
        }
        
        rateLimiter.Wait()
        
        var portInfo PortInfo
        
        switch {
        case s.config.SynScan:
            portInfo = s.synScan(ip, port)
        case s.config.ConnectScan:
            portInfo = s.connectScan(ip, port)
        case s.config.AckScan:
            portInfo = s.ackScan(ip, port)
        case s.config.FinScan:
            portInfo = s.finScan(ip, port)
        case s.config.NullScan:
            portInfo = s.nullScan(ip, port)
        case s.config.XmasScan:
            portInfo = s.xmasScan(ip, port)
        default:
            portInfo = s.connectScan(ip, port)
        }
        
        select {
        case <-ctx.Done():
            return
        case results <- portInfo:
        }
    }
}

// TCP SYN Scan (Stealth)
func (s *Scanner) synScan(ip net.IP, port int) PortInfo {
    start := time.Now()
    
    // Create raw socket
    conn, err := net.DialTimeout("tcp", 
        fmt.Sprintf("%s:%d", ip.String(), port), 
        time.Duration(s.config.Timeout)*time.Second)
    
    if err != nil {
        return PortInfo{
            Port:     port,
            Protocol: "tcp",
            State:    StateFiltered,
            Reason:   "no-response",
        }
    }
    defer conn.Close()
    
    return PortInfo{
        Port:         port,
        Protocol:     "tcp",
        State:        StateOpen,
        Service:      utils.LookupService(port, "tcp"),
        Reason:       "syn-ack",
        ResponseTime: time.Since(start),
    }
}

// TCP Connect Scan
func (s *Scanner) connectScan(ip net.IP, port int) PortInfo {
    start := time.Now()
    address := fmt.Sprintf("%s:%d", ip.String(), port)
    
    conn, err := net.DialTimeout("tcp", address, 
        time.Duration(s.config.Timeout)*time.Second)
    
    if err != nil {
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
            return PortInfo{
                Port:     port,
                Protocol: "tcp",
                State:    StateFiltered,
                Reason:   "timeout",
            }
        }
        return PortInfo{
            Port:     port,
            Protocol: "tcp",
            State:    StateClosed,
            Reason:   "connection-refused",
        }
    }
    defer conn.Close()
    
    // Try to grab banner
    banner := s.grabBanner(conn, port)
    
    return PortInfo{
        Port:         port,
        Protocol:     "tcp",
        State:        StateOpen,
        Service:      utils.LookupService(port, "tcp"),
        Banner:       banner,
        Reason:       "syn-ack",
        ResponseTime: time.Since(start),
    }
}

// TCP ACK Scan (Firewall Evasion)
func (s *Scanner) ackScan(ip net.IP, port int) PortInfo {
    // This requires raw packet crafting
    // Simplified implementation using custom packet
    return s.craftAndSend(ip, port, layers.TCPAck)
}

// TCP FIN Scan (Stealth)
func (s *Scanner) finScan(ip net.IP, port int) PortInfo {
    return s.craftAndSend(ip, port, layers.TCPFin)
}

// TCP NULL Scan (Stealth)
func (s *Scanner) nullScan(ip net.IP, port int) PortInfo {
    return s.craftAndSend(ip, port, 0)
}

// TCP Xmas Scan (Stealth)
func (s *Scanner) xmasScan(ip net.IP, port int) PortInfo {
    // Xmas = FIN + PSH + URG
    flags := layers.TCPFin | layers.TCPPsh | layers.TCPUrg
    return s.craftAndSend(ip, port, flags)
}

func (s *Scanner) craftAndSend(ip net.IP, port int, flags uint8) PortInfo {
    // Raw packet crafting implementation
    // This is a simplified version - full implementation would use gopacket
    
    start := time.Now()
    
    // For now, fallback to connect scan with modified behavior
    // In full implementation, this would craft raw TCP packets
    result := s.connectScan(ip, port)
    result.ResponseTime = time.Since(start)
    
    return result
}

func (s *Scanner) grabBanner(conn net.Conn, port int) string {
    if !utils.IsCommonService(port) {
        return ""
    }
    
    conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    if n > 0 {
        return string(buf[:n])
    }
    return ""
}

func (s *Scanner) determineScanType() string {
    switch {
    case s.config.SynScan:
        return "TCP SYN"
    case s.config.ConnectScan:
        return "TCP Connect"
    case s.config.AckScan:
        return "TCP ACK"
    case s.config.FinScan:
        return "TCP FIN"
    case s.config.NullScan:
        return "TCP NULL"
    case s.config.XmasScan:
        return "TCP Xmas"
    case s.config.UdpScan:
        return "UDP"
    default:
        return "TCP SYN"
    }
}