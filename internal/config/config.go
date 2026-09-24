package config

import (
    "fmt"
    "os"
)

// Config holds all scanner configuration
type Config struct {
    // Target
    Target string
    
    // Scan types
    SynScan       bool
    ConnectScan   bool
    AckScan       bool
    FinScan       bool
    NullScan      bool
    XmasScan      bool
    UdpScan       bool
    
    // Port configuration
    Ports     string
    PortList  []int
    TopPorts  bool
    AllPorts  bool
    
    // Evasion
    Fragment    bool
    FragmentMTU int
    SourcePort  int
    SpoofIP     string
    BadSum      bool
    EvasionMode bool
    
    // Performance
    Rate        int
    Timeout     int
    Retries     int
    Concurrency int
    
    // Output
    Output  string
    Verbose bool
    NoColor bool
    
    // Host discovery
    Ping      bool
    SkipHost  bool
    
    // Filtered handling
    ScanFiltered bool
}

// New creates a new Config with defaults
func New() *Config {
    return &Config{
        SynScan:     true,
        Ports:       "1-1000",
        Rate:        1000,
        Timeout:     3,
        Retries:     2,
        Concurrency: 100,
        FragmentMTU: 16,
        Ping:        true,
    }
}

// Validate checks configuration validity
func (c *Config) Validate() error {
    if c.Target == "" {
        return fmt.Errorf("target is required")
    }
    
    if os.Getuid() != 0 && (c.SynScan || c.AckScan || c.FinScan || c.NullScan || c.XmasScan || c.UdpScan) {
        fmt.Println("[!] Warning: Some scans require root privileges. Falling back to connect scan.")
        c.SynScan = false
        c.ConnectScan = true
    }
    
    if c.Rate < 1 {
        return fmt.Errorf("rate must be at least 1")
    }
    
    if c.Concurrency < 1 {
        return fmt.Errorf("concurrency must be at least 1")
    }
    
    return nil
}