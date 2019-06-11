package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	mutexAllowList sync.Mutex
	allowPhoneList []string
	allowIpList    []string
)

func getContextWithTimeout(ms int) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), (time.Duration)(ms)*time.Millisecond)
	return ctx
}

func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	return dir
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	return lines
}

func isPhoneNumberAllow(phoneNumber string) bool {
	mutexAllowList.Lock()
	defer mutexAllowList.Unlock()
	for _, s := range allowPhoneList {
		if s == phoneNumber {
			return true
		}
	}
	return false
}

func compareCIDR(ip string, cidr string) bool {
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		dump(err)
		return false
	}
	return subnet.Contains(net.ParseIP(ip))
}

func isIpAllow(ip string) bool {
	mutexAllowList.Lock()
	defer mutexAllowList.Unlock()
	for _, s := range allowIpList {
		if ip == s || (strings.Index(s, "/") >= 0 && compareCIDR(ip, s)) {
			return true
		}
	}
	return false
}

func printError(info string, err error) {
	if info != "" {
		print(info+":", err)
	} else {
		print(err)
	}
}
