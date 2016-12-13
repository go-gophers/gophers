package net

import (
	"math/rand"
	"net"
	"sync"

	"github.com/go-gophers/gophers/config"
)

var (
	dns   = make(map[string][]string)
	dnsRW sync.RWMutex
)

// FlushDNSCache flushes DNS cache used by Dial.
func FlushDNSCache() {
	dnsRW.Lock()
	dns = make(map[string][]string)
	dnsRW.Unlock()
}

func lookupHost(host string) (string, error) {
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.String(), nil
	}

	dnsRW.RLock()
	addrs := dns[host]
	dnsRW.RUnlock()
	if addrs != nil {
		return addrs[rand.Intn(len(addrs))], nil
	}

	dnsRW.Lock()
	defer dnsRW.Unlock()
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	addrs = make([]string, 0, len(ips))
	for _, ip := range ips {
		if config.Default.DisableIPv6 {
			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}
		}
		addrs = append(addrs, ip.String())
	}

	dns[host] = addrs
	return addrs[rand.Intn(len(addrs))], nil
}
