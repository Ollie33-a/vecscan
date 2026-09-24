package output

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/vectalithlabs/vecscan/pkg/scanner"
)

// GeneratePartialReport creates a report when scan is interrupted
func GeneratePartialReport(results *scanner.Results, duration time.Duration) {
	fmt.Println()
	color.Yellow("═══════════════════════════════════════════════════════════")
	color.Yellow("                  PARTIAL SCAN REPORT")
	color.Yellow("═══════════════════════════════════════════════════════════")
	
	fmt.Printf("Target:      %s\n", results.Target)
	fmt.Printf("Scan Type:   %s\n", results.ScanType)
	fmt.Printf("Duration:    %s\n", duration.Round(time.Millisecond))
	fmt.Printf("Status:      INTERRUPTED\n")
	
	fmt.Println()
	color.Cyan("PROGRESS:")
	fmt.Printf("  Ports Scanned: %d\n", results.Statistics.TotalPorts)
	fmt.Printf("  Open Found:    %d\n", results.Statistics.OpenCount)
	fmt.Printf("  Closed:        %d\n", results.Statistics.ClosedCount)
	fmt.Printf("  Filtered:      %d\n", results.Statistics.FilteredCount)
	
	if len(results.OpenPorts) > 0 {
		fmt.Println()
		color.Green("OPEN PORTS FOUND SO FAR:")
		for _, port := range results.OpenPorts {
			fmt.Printf("  %d/%s - %s\n", port.Port, port.Protocol, port.Service)
		}
	}
	
	fmt.Println()
	color.Yellow("[!] Scan was interrupted. Results may be incomplete.")
}

// SavePartialJSON saves partial results
func SavePartialJSON(results *scanner.Results, filename string) error {
	results.EndTime = time.Now()
	return SaveJSON(results, filename)
}

// PrintHostDiscovery shows discovered hosts
func PrintHostDiscovery(hosts []string) {
	if len(hosts) == 0 {
		color.Yellow("[!] No hosts discovered")
		return
	}
	
	color.Green("\n[+] Discovered %d host(s):", len(hosts))
	for _, host := range hosts {
		fmt.Printf("    %s\n", host)
	}
}