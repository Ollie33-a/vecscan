package output

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/fatih/color"
    "github.com/vectalithlabs/vecscan/pkg/scanner"
)

// PrintSummary prints a formatted scan summary
func PrintSummary(results *scanner.Results, duration time.Duration, verbose bool) {
    fmt.Println()
    color.Cyan("═══════════════════════════════════════════════════════════")
    color.Cyan("                     SCAN SUMMARY")
    color.Cyan("═══════════════════════════════════════════════════════════")
    
    fmt.Printf("Target:      %s\n", results.Target)
    fmt.Printf("Scan Type:   %s\n", results.ScanType)
    fmt.Printf("Duration:    %s\n", duration.Round(time.Millisecond))
    fmt.Printf("Start Time:  %s\n", results.StartTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("End Time:    %s\n", results.EndTime.Format("2006-01-02 15:04:05"))
    
    fmt.Println()
    color.Cyan("PORT STATISTICS:")
    color.Green("  Open:     %d", results.Statistics.OpenCount)
    color.Red("  Closed:   %d", results.Statistics.ClosedCount)
    color.Yellow("  Filtered: %d", results.Statistics.FilteredCount)
    fmt.Printf("  Total:    %d\n", results.Statistics.TotalPorts)
    
    if len(results.OpenPorts) > 0 {
        fmt.Println()
        color.Green("OPEN PORTS:")
        fmt.Printf("%-8s %-10s %-12s %-20s %s\n", 
            "PORT", "PROTOCOL", "STATE", "SERVICE", "RESPONSE")
        fmt.Println("─────────────────────────────────────────────────────────")
        
        for _, port := range results.OpenPorts {
            fmt.Printf("%-8d %-10s %-12s %-20s %s\n",
                port.Port,
                port.Protocol,
                port.State,
                port.Service,
                port.ResponseTime.Round(time.Millisecond))
            
            if verbose && port.Banner != "" {
                fmt.Printf("  └─ Banner: %s\n", truncate(port.Banner, 60))
            }
        }
    }
    
    fmt.Println()
}

func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..."
}

// SaveJSON saves results to a JSON file
func SaveJSON(results *scanner.Results, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(results)
}