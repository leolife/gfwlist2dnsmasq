package cmd

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
)

func isIPv4(address string) bool {
	ip := net.ParseIP(address)
	return ip.To4() != nil && strings.Contains(address, ".")
}

func isDomain(domain string) bool {
	if len(domain) > 255 {
		return false
	}
	if strings.HasSuffix(domain, ".") {
		domain = strings.TrimSuffix(domain, ".")
	}
	return true
}

func getHostname(str string) (string, error) {
	if str == "" {
		return "", errors.New("str is empty")
	}

	strTemp := str
	if !strings.Contains(str, "://") {
		strTemp = "http://" + str
	}

	u, err := url.Parse(strTemp)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(u.Host, ".") || u.Host == "" {
		return "", errors.New("Parse err")
	}
	return u.Host, nil
}

// TLDList is the the list of top-level domain
var (
	TLDList map[string]bool
	TLDOnce sync.Once
)

func getTLDList() map[string]bool {
	TLDOnce.Do(func() {
		r, err := os.Open("../resources/tld.txt")
		if err != nil {
			return
		}
		defer r.Close()

		scanner := bufio.NewScanner(r)
		tld := make(map[string]bool)
		for scanner.Scan() {
			tld[scanner.Text()] = true
		}
		TLDList = tld
	})

	return TLDList
}

func reduceDomain(str string) string {
	parts := strings.Split(str, ".")
	tld := getTLDList()

	var lastMix string
	for i := range parts {
		mix := strings.Join(parts[len(parts)-i-1:], ".")
		if i == 0 {
			if _, ok := tld[mix]; !ok {
				break
			}
		}
		lastMix = mix
		if _, ok := tld[mix]; ok {
			continue
		} else {
			break
		}
	}

	return lastMix
}

// ParseList parse content
func ParseList(content io.Reader) (map[string]bool, error) {
	domains := make(map[string]bool)

	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ".*") {
			continue
		} else if strings.Contains(line, "*") {
			line = strings.Replace(line, "*", "/", -1)
		}

		if strings.HasPrefix(line, "||") {
			line = strings.TrimPrefix(line, "||")
		} else if strings.HasPrefix(line, "|") {
			line = strings.TrimPrefix(line, "|")
		} else if strings.HasPrefix(line, ".") {
			line = strings.TrimPrefix(line, ".")
		}

		if strings.HasPrefix(line, "!") {
			continue
		} else if strings.HasPrefix(line, "[") {
			continue
		} else if strings.HasPrefix(line, "@") {
			continue
		}

		hostname, err := getHostname(line)
		if err != nil {
			continue
		}
		if isIPv4(hostname) {
			continue
		}

		hostname = reduceDomain(hostname)
		if hostname == "" {
			continue
		}

		domains[hostname] = true
	}
	return domains, nil
}
