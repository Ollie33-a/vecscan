package evasion

import (
    "math/rand"
    "net"
    "time"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
)

// StealthScanner implements advanced evasion techniques
type StealthScanner struct {
    sourcePort int
    spoofIP    net.IP
    useBadSum  bool
}

// NewStealthScanner creates a new stealth scanner
func NewStealthScanner(sourcePort int, spoofIP string, badSum bool) *StealthScanner {
    var spoof net.IP
    if spoofIP != "" {
        spoof = net.ParseIP(spoofIP)
    }
    
    return &StealthScanner{
        sourcePort: sourcePort,
        spoofIP:    spoof,
        useBadSum:  badSum,
    }
}

// CraftSYNPacket creates a stealth SYN packet
func (s *StealthScanner) CraftSYNPacket(srcIP, dstIP net.IP, srcPort, dstPort int) ([]byte, error) {
    // IP Layer
    ipLayer := &layers.IPv4{
        Version:  4,
        IHL:      5,
        TTL:      randomTTL(),
        Protocol: layers.IPProtocolTCP,
        SrcIP:    s.getSourceIP(srcIP),
        DstIP:    dstIP,
    }
    
    // TCP Layer
    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(s.getSourcePort(srcPort)),
        DstPort: layers.TCPPort(dstPort),
        Seq:     randomSeq(),
        Window:  randomWindow(),
        SYN:     true,
    }
    
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    
    if s.useBadSum {
        tcpLayer.Checksum = 0xDEAD // Bogus checksum
    }
    
    // Serialize
    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{
        ComputeChecksums: !s.useBadSum,
        FixLengths:       true,
    }
    
    err := gopacket.SerializeLayers(buf, opts, ipLayer, tcpLayer)
    if err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

// CraftFINPacket creates a stealth FIN packet
func (s *StealthScanner) CraftFINPacket(srcIP, dstIP net.IP, srcPort, dstPort int) ([]byte, error) {
    ipLayer := &layers.IPv4{
        Version:  4,
        IHL:      5,
        TTL:      randomTTL(),
        Protocol: layers.IPProtocolTCP,
        SrcIP:    s.getSourceIP(srcIP),
        DstIP:    dstIP,
    }
    
    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(s.getSourcePort(srcPort)),
        DstPort: layers.TCPPort(dstPort),
        Seq:     randomSeq(),
        Window:  0,
        FIN:     true,
    }
    
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    
    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{
        ComputeChecksums: true,
        FixLengths:       true,
    }
    
    err := gopacket.SerializeLayers(buf, opts, ipLayer, tcpLayer)
    return buf.Bytes(), err
}

// CraftXmasPacket creates a Christmas tree packet (FIN+PSH+URG)
func (s *StealthScanner) CraftXmasPacket(srcIP, dstIP net.IP, srcPort, dstPort int) ([]byte, error) {
    ipLayer := &layers.IPv4{
        Version:  4,
        IHL:      5,
        TTL:      randomTTL(),
        Protocol: layers.IPProtocolTCP,
        SrcIP:    s.getSourceIP(srcIP),
        DstIP:    dstIP,
    }
    
    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(s.getSourcePort(srcPort)),
        DstPort: layers.TCPPort(dstPort),
        Seq:     randomSeq(),
        Window:  0,
        FIN:     true,
        PSH:     true,
        URG:     true,
    }
    
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    
    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{
        ComputeChecksums: true,
        FixLengths:       true,
    }
    
    err := gopacket.SerializeLayers(buf, opts, ipLayer, tcpLayer)
    return buf.Bytes(), err
}

// CraftNULLPacket creates a NULL packet (no flags)
func (s *StealthScanner) CraftNULLPacket(srcIP, dstIP net.IP, srcPort, dstPort int) ([]byte, error) {
    ipLayer := &layers.IPv4{
        Version:  4,
        IHL:      5,
        TTL:      randomTTL(),
        Protocol: layers.IPProtocolTCP,
        SrcIP:    s.getSourceIP(srcIP),
        DstIP:    dstIP,
    }
    
    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(s.getSourcePort(srcPort)),
        DstPort: layers.TCPPort(dstPort),
        Seq:     randomSeq(),
        Window:  0,
        // No flags set - NULL packet
    }
    
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    
    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{
        ComputeChecksums: true,
        FixLengths:       true,
    }
    
    err := gopacket.SerializeLayers(buf, opts, ipLayer, tcpLayer)
    return buf.Bytes(), err
}

func (s *StealthScanner) getSourcePort(defaultPort int) int {
    if s.sourcePort > 0 {
        return s.sourcePort
    }
    return defaultPort
}

func (s *StealthScanner) getSourceIP(defaultIP net.IP) net.IP {
    if s.spoofIP != nil {
        return s.spoofIP
    }
    return defaultIP
}

func randomTTL() uint8 {
    ttlValues := []uint8{64, 128, 255, 53, 128}
    return ttlValues[rand.Intn(len(ttlValues))]
}

func randomSeq() uint32 {
    return uint32(rand.Int31())
}

func randomWindow() uint16 {
    windows := []uint16{1024, 2048, 4096, 8192, 14600, 64240}
    return windows[rand.Intn(len(windows))]
}

func init() {
    rand.Seed(time.Now().UnixNano())
}