package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/fatih/color"
    "github.com/spf13/cobra"
    "github.com/vectalithlabs/vecscan/internal/config"
    "github.com/vectalithlabs/vecscan/pkg/output"
    "github.com/vectalithlabs/vecscan/pkg/scanner"
    "github.com/vectalithlabs/vecscan/pkg/utils"
)

var (
    version = "1.0.0"
    cfg     *config.Config
)

var rootCmd = &cobra.Command{
    Use:   "vecscan [target]",
    Short: "Advanced Network Scanner by Vectalith Labs",
    Long: `
██╗   ██╗███████╗ ██████╗███████╗ ██████╗ █████╗ ███╗   ██╗
██║   ██║██╔════╝██╔════╝██╔════╝██╔════╝██╔══██╗████╗  ██║
██║   ██║█████╗  ██║     ███████╗██║     ███████║██╔██╗ ██║
╚██╗ ██╔╝██╔══╝  ██║     ╚════██║██║     ██╔══██║██║╚██╗██║
 ╚████╔╝ ███████╗╚██████╗███████║╚██████╗██║  ██║██║ ╚████║
  ╚═══╝  ╚══════╝ ╚═════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝
                    BY VECTALITH LABS v` + version + `

VecScan is an elite network reconnaissance tool featuring:
• TCP SYN/Connect/ACK/Window/Maimon scanning
• Stealth scans (FIN/NULL/Xmas)
• UDP scanning with protocol detection
• Advanced firewall evasion (fragmentation, spoofing)
• Rate limiting and concurrency control
• JSON and terminal output formats
`,
    Args: cobra.ExactArgs(1),
    Run:  runScan,
}

func init() {
    cfg = config.New()
    
    // Scan types
    rootCmd.Flags().BoolVarP(&cfg.SynScan, "syn", "sS", true, "TCP SYN scan (default)")
    rootCmd.Flags().BoolVar(&cfg.ConnectScan, "connect", false, "TCP Connect scan")
    rootCmd.Flags().BoolVarP(&cfg.AckScan, "ack", "sA", false, "TCP ACK scan (firewall evasion)")
    rootCmd.Flags().BoolVarP(&cfg.FinScan, "fin", "sF", false, "TCP FIN scan (stealth)")
    rootCmd.Flags().BoolVarP(&cfg.NullScan, "null", "sN", false, "TCP NULL scan (stealth)")
    rootCmd.Flags().BoolVarP(&cfg.XmasScan, "xmas", "sX", false, "TCP Xmas scan (stealth)")
    rootCmd.Flags().BoolVarP(&cfg.UdpScan, "udp", "sU", false, "UDP scan")
    
    // Port configuration
    rootCmd.Flags().StringVarP(&cfg.Ports, "ports", "p", "1-1000", "Port range (e.g., 80,443,1-1000,8080-8090)")
    rootCmd.Flags().BoolVar(&cfg.TopPorts, "top-ports", false, "Scan top 1000 common ports")
    rootCmd.Flags().BoolVar(&cfg.AllPorts, "all-ports", false, "Scan all 65535 ports")
    
    // Evasion techniques
    rootCmd.Flags().BoolVar(&cfg.Fragment, "fragment", false, "Fragment packets (firewall evasion)")
    rootCmd.Flags().IntVar(&cfg.FragmentMTU, "mtu", 16, "Fragment MTU size")
    rootCmd.Flags().IntVar(&cfg.SourcePort, "source-port", 0, "Source port for scans")
    rootCmd.Flags().StringVar(&cfg.SpoofIP, "spoof", "", "Spoof source IP (requires root)")
    rootCmd.Flags().BoolVar(&cfg.BadSum, "badsum", false, "Send packets with bogus checksum")
    
    // Performance
    rootCmd.Flags().IntVar(&cfg.Rate, "rate", 1000, "Packets per second rate limit")
    rootCmd.Flags().IntVarP(&cfg.Timeout, "timeout", "t", 3, "Timeout in seconds")
    rootCmd.Flags().IntVar(&cfg.Retries, "retries", 2, "Number of retries")
    rootCmd.Flags().IntVarP(&cfg.Concurrency, "concurrency", "c", 100, "Concurrent goroutines")
    
    // Output
    rootCmd.Flags().StringVarP(&cfg.Output, "output", "o", "", "Output file (JSON format)")
    rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")
    rootCmd.Flags().BoolVar(&cfg.NoColor, "no-color", false, "Disable colored output")
    
    // Host discovery
    rootCmd.Flags().BoolVar(&cfg.Ping, "ping", true, "Perform host discovery")
    rootCmd.Flags().BoolVar(&cfg.SkipHost, "skip-host", false, "Skip host discovery (treat all hosts as up)")
    
    // Filtered port handling
    rootCmd.Flags().BoolVar(&cfg.ScanFiltered, "scan-filtered", false, "Automatically scan filtered ports")
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func runScan(cmd *cobra.Command, args []string>) {
    cfg.Target = args[0]
    
    if cfg.NoColor {
        color.NoColor = true
    }
    
    // Create cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Setup signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Handle shutdown in goroutine
    go func() {
        <-sigChan
        color.Yellow("\n\n[!] Interrupt received, initiating graceful shutdown...")
        cancel()
    }()
    
    // Print banner
    printBanner()
    
    // Validate and parse configuration
    if err := cfg.Validate(); err != nil {
        color.Red("[!] Configuration error: %v", err)
        os.Exit(1)
    }
    
    // Parse port ranges
    ports, err := utils.ParsePortRange(cfg.Ports)
    if err != nil {
        color.Red("[!] Port range error: %v", err)
        os.Exit(1)
    }
    cfg.PortList = ports
    
    // Create scanner
    scan := scanner.New(cfg)
    
    // Run scan with context
    startTime := time.Now()
    results, err := scan.Execute(ctx)
    if err != nil && err != context.Canceled {
        color.Red("[!] Scan error: %v", err)
    }
    
    duration := time.Since(startTime)
    
    // Generate reports
    if cfg.Output != "" {
        if err := output.SaveJSON(results, cfg.Output); err != nil {
            color.Red("[!] Failed to save JSON: %v", err)
        } else {
            color.Green("[+] Results saved to: %s", cfg.Output)
        }
    }
    
    // Print summary
    output.PrintSummary(results, duration, cfg.Verbose)
    
    // Handle filtered ports if any
    if len(results.FilteredPorts) > 0 && !cfg.ScanFiltered && ctx.Err() != context.Canceled {
        handleFilteredPorts(ctx, scan, results)
    }
}

func printBanner() {
    banner := `
╔═══════════════════════════════════════════════════════════╗
║  ██╗   ██╗███████╗ ██████╗███████╗ ██████╗ █████╗ ███╗   ██╗  ║
║  ██║   ██║██╔════╝██╔════╝██╔════╝██╔════╝██╔══██╗████╗  ██║  ║
║  ██║   ██║█████╗  ██║     ███████╗██║     ███████║██╔██╗ ██║  ║
║  ╚██╗ ██╔╝██╔══╝  ██║     ╚════██║██║     ██╔══██║██║╚██╗██║  ║
║   ╚████╔╝ ███████╗╚██████╗███████║╚██████╗██║  ██║██║ ╚████║  ║
║    ╚═══╝  ╚══════╝ ╚═════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝  ║
║                    BY VECTALITH LABS                      ║
╚═══════════════════════════════════════════════════════════╝
`
    color.Cyan(banner)
}

func handleFilteredPorts(ctx context.Context, scan *scanner.Scanner, results *scanner.Results>) {
    color.Yellow("\n[!] Found %d filtered ports", len(results.FilteredPorts))
    color.Yellow("    These ports may be protected by a firewall")
    
    if !cfg.ScanFiltered {
        fmt.Print("\n[?] Scan filtered ports with evasion techniques? [Y/n]: ")
        var response string
        fmt.Scanln(&response)
        
        if strings.ToLower(response) == "y" || response == "" {
            color.Cyan("\n[*] Initiating filtered port scan with evasion...")
            cfg.ScanFiltered = true
            cfg.EvasionMode = true
            
            filteredResults, err := scan.ScanFilteredPorts(ctx, results.FilteredPorts)
            if err != nil {
                color.Red("[!] Filtered scan error: %v", err)
                return
            }
            
            // Merge results
            results.OpenPorts = append(results.OpenPorts, filteredResults.OpenPorts...)
            results.FilteredPorts = filteredResults.FilteredPorts
            
            color.Green("[+] Filtered scan complete. Found %d additional open ports", 
                len(filteredResults.OpenPorts))
        }
    }
}