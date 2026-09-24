package evasion

import (
    "encoding/binary"
    "net"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
)

// Fragmenter implements IP packet fragmentation for evasion
type Fragmenter struct {
    MTU int
}

// NewFragmenter creates a new fragmenter with specified MTU
func NewFragmenter(mtu int) *Fragmenter {
    if mtu < 8 {
        mtu = 8
    }
    return &Fragmenter{MTU: mtu}
}

// FragmentPacket fragments a TCP packet into smaller pieces
func (f *Fragmenter) FragmentTCPPacket(srcIP, dstIP net.IP, srcPort, dstPort int, 
    payload []byte) ([][]byte, error) {
    
    // Create TCP header
    tcpLayer := &layers.TCP{
        SrcPort: layers.TCPPort(srcPort),
        DstPort: layers.TCPPort(dstPort),
        Seq:     0,
        ACK:     0,
        SYN:     true,
        Window:  64240,
    }
    
    ipLayer := &layers.IPv4{
        Version:  4,
        IHL:      5,
        TTL:      64,
        Protocol: layers.IPProtocolTCP,
        SrcIP:    srcIP,
        DstIP:    dstIP,
    }
    
    tcpLayer.SetNetworkLayerForChecksum(ipLayer)
    
    // Serialize TCP layer
    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}
    
    err := gopacket.SerializeLayers(buf, opts, tcpLayer, gopacket.Payload(payload))
    if err != nil {
        return nil, err
    }
    
    packetData := buf.Bytes()
    
    // Calculate fragment size (must be multiple of 8)
    fragmentSize := (f.MTU / 8) * 8
    if fragmentSize < 8 {
        fragmentSize = 8
    }
    
    var fragments [][]byte
    offset := 0
    
    for offset < len(packetData) {
        end := offset + fragmentSize
        if end > len(packetData) {
            end = len(packetData)
        }
        
        fragData := packetData[offset:end]
        fragOffset := uint16(offset / 8)
        moreFragments := end < len(packetData)
        
        // Create IP fragment
        fragIP := &layers.IPv4{
            Version:        4,
            IHL:            5,
            TTL:            64,
            Protocol:       layers.IPProtocolTCP,
            SrcIP:          srcIP,
            DstIP:          dstIP,
            FragOffset:     fragOffset,
            Flags:          layers.IPv4MoreFragments,
            Length:         uint16(len(fragData) + 20),
            Id:             0x1234,
        }
        
        if !moreFragments {
            fragIP.Flags &^= layers.IPv4MoreFragments
        }
        
        // Serialize fragment
        fragBuf := gopacket.NewSerializeBuffer()
        err := gopacket.SerializeLayers(fragBuf, 
            gopacket.SerializeOptions{FixLengths: true}, 
            fragIP, 
            gopacket.Payload(fragData))
        
        if err != nil {
            return nil, err
        }
        
        fragments = append(fragments, fragBuf.Bytes())
        offset = end
    }
    
    return fragments, nil
}

// SendFragments sends fragmented packets with delay
func (f *Fragmenter) SendFragments(fragments [][]byte, sender func([]byte) error) error {
    for i, frag := range fragments {
        if err := sender(frag); err != nil {
            return err
        }
        
        // Add small delay between fragments to evade detection
        if i < len(fragments)-1 {
            time.Sleep(time.Microsecond * 100)
        }
    }
    return nil
}