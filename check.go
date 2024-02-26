package main

import (
	"net"
	"time"
)


type ServiceStatus struct {
	WebsiteName    string
	URLChecked     string
	ResponseTime   string
	LastDown       string
	CurrentStatus  string
	LocalAddress   string
	RemoteAddress  string
	ServiceHistory []ServiceHistoryEntry
}

type ServiceHistoryEntry struct {
	Date     string
	Time     string
	PingTime string
}


func Check(domain, port string) ServiceStatus {
	address := domain + ":" + port
	timeout := 5 * time.Second
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	responseTime := time.Since(startTime).String()
	var lastDown string
	var currentStatus string
	var localAddress string
	var remoteAddress string
	var serviceHistory []ServiceHistoryEntry

	if err != nil {
		currentStatus = "DOWN"
		lastDown = time.Now().String()
	} else {
		currentStatus = "UP"
		lastDown = "N/A"
		localAddress = conn.LocalAddr().String()
		remoteAddress = conn.RemoteAddr().String()
	}
	
	if conn != nil {
		conn.Close()
	}

	serviceHistory = append(serviceHistory, ServiceHistoryEntry{
		Date:     time.Now().Format("02.Jan.2006"),
		Time:     time.Now().Format("15:04"),
		PingTime: responseTime,
	})

	return ServiceStatus{
		WebsiteName:    domain,
		URLChecked:     domain,
		ResponseTime:   responseTime,
		LastDown:       lastDown,
		CurrentStatus:  currentStatus,
		LocalAddress:   localAddress,
		RemoteAddress:  remoteAddress,
		ServiceHistory: serviceHistory,
	}
}
