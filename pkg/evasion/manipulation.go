package evasion

import (
	"math/rand"
	"net"
	"time"
)

// PacketManipulator implements advanced packet manipulation
type PacketManipulator struct {
	decoyIPs []net.IP
}

// NewPacketManipulator creates a new manipulator
func NewPacketManipulator() *PacketManipulator {
	return &PacketManipulator{
		decoyIPs: generateDecoyIPs(),
	}
}

// GenerateDecoyIPs creates fake source IPs for noise generation
func generateDecoyIPs() []net.IP {
	var ips []net.IP
	// Generate decoys from RFC1918 ranges
	ranges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	
	for _, r := range ranges {
		_, ipnet, _ := net.ParseCIDR(r)
		if ipnet != nil {
			ip := make(net.IP, 4)
			copy(ip, ipnet.IP)
			// Randomize last octet
			ip[3] = byte(rand.Intn(254) + 1)
			ips = append(ips, ip)
		}
	}
	
	return ips
}

// RandomizeTiming returns random delay for timing obfuscation
func RandomizeTiming() time.Duration {
	// Random delay between 1-100ms
	return time.Duration(rand.Intn(100)+1) * time.Millisecond
}

// RandomizeOrder shuffles scan order to evade sequence detection
func RandomizeOrder(ports []int) []int {
	rand.Shuffle(len(ports), func(i, j int) {
		ports[i], ports[j] = ports[j], ports[i]
	})
	return ports
}

// GenerateTCPOptions creates variable TCP options for fingerprint evasion
func GenerateTCPOptions() []byte {
	// Common TCP option combinations
	options := [][]byte{
		// MSS + SACK + Timestamp
		{0x02, 0x04, 0x05, 0xb4, 0x01, 0x01, 0x04, 0x02},
		// MSS only
		{0x02, 0x04, 0x05, 0xb4},
		// Window Scale + MSS + Timestamp + SACK
		{0x03, 0x03, 0x07, 0x02, 0x04, 0x05, 0xb4, 0x08, 0x0a},
	}
	
	return options[rand.Intn(len(options))]
}

// CalculateJitter adds randomization to timing
func CalculateJitter(baseDelay time.Duration, jitterPercent float64) time.Duration {
	jitter := time.Duration(float64(baseDelay) * jitterPercent * (rand.Float64()*2 - 1))
	return baseDelay + jitter
}