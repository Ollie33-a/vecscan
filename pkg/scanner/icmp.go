package scanner

import (
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"://github.com"
)

// Ping performs ICMP echo request
func (s *Scanner) Ping(ip net.IP) (bool, time.Duration, error) {
	start := time.Now()
	
	// Create ICMP connection
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false, 0, err
	}
	defer c.Close()
	
	// Create ICMP message
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte("VECSCAN"),
		},
	}
	
	data, err := msg.Marshal(nil)
	if err != nil {
		return false, 0, err
	}
	
	// Send packet
	_, err = c.WriteTo(data, &net.IPAddr{IP: ip})
	if err != nil {
		return false, 0, err
	}
	
	// Set timeout
	c.SetReadDeadline(time.Now().Add(time.Duration(s.config.Timeout) * time.Second))
	
	// Read reply
	reply := make([]byte, 1500)
	n, _, err := c.ReadFrom(reply)
	if err != nil {
		return false, 0, err
	}
	
	duration := time.Since(start)
	
	// Parse reply
	replyMsg, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		return false, 0, err
	}
	
	if replyMsg.Type == ipv4.ICMPTypeEchoReply {
		return true, duration, nil
	}
	
	return false, duration, nil
}

// HostDiscovery scans network for live hosts
func (s *Scanner) HostDiscovery(subnet *net.IPNet) ([]net.IP, error) {
	var hosts []net.IP
	
	for ip := subnet.IP.Mask(subnet.Mask); subnet.Contains(ip); incrementIP(ip) {
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		
		alive, _, err := s.Ping(ipCopy)
		if err == nil && alive {
			hosts = append(hosts, ipCopy)
		}
	}
	
	return hosts, nil
}