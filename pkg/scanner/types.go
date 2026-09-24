package scanner

import (
	"net"
	"time"

	"github.com/Ollie33-a/vecscan"
)

type PortState string

const (
	StateOpen         PortState = "open"
	StateClosed       PortState = "closed"
	StateFiltered     PortState = "filtered"
	StateUnfiltered   PortState = "unfiltered"
	StateOpenFiltered PortState = "open|filtered"
)

type PortInfo struct {
	Port         int           `json:"port"`
	Protocol     string        `json:"protocol"`
	State        PortState     `json:"state"`
	Service      string        `json:"service"`
	Version      string        `json:"version,omitempty"`
	Banner       string        `json:"banner,omitempty"`
	Reason       string        `json:"reason"`
	ResponseTime time.Duration `json:"response_time"`
}

type HostResult struct {
	IP        string     `json:"ip"`
	Hostname  string     `json:"hostname,omitempty"`
	Status    string     `json:"status"`
	OpenPorts []PortInfo `json:"open_ports"`
	Latency   time.Duration `json:"latency,omitempty"`
	OS        string     `json:"os,omitempty"`
}

type Results struct {
	Target        string       `json:"target"`
	ScanType      string       `json:"scan_type"`
	StartTime     time.Time    `json:"start_time"`
	EndTime       time.Time    `json:"end_time"`
	Hosts         []HostResult `json:"hosts"`
	OpenPorts     []PortInfo   `json:"open_ports"`
	FilteredPorts []PortInfo   `json:"filtered_ports,omitempty"`
	Statistics    Statistics   `json:"statistics"`
}

type Statistics struct {
	TotalPorts    int `json:"total_ports"`
	OpenCount     int `json:"open_count"`
	ClosedCount   int `json:"closed_count"`
	FilteredCount int `json:"filtered_count"`
	PacketsSent   int `json:"packets_sent"`
	PacketsRcvd   int `json:"packets_received"`
}

type Scanner struct {
	config  *config.Config
	stats   Statistics
	results Results
}