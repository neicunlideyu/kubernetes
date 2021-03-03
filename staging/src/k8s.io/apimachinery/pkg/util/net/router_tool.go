package net

import (
	"errors"
	"k8s.io/klog"
	"net"
)

var (
	errNotSupportFamily   = errors.New("not support network family")
	errCannotGetDefaultIP = errors.New("cannot get default ip address")
)

func getDefaultIPAddressByDial(addressFamilies AddressFamilyPreference) (net.IP, error) {
	var (
		v4Targets = []string{"8.8.8.8:53", "10.8.8.8:53"}
		v6Targets = []string{"[240C::6644]:53", "[fdbd:dc00::10:8:8:8]:53"}
	)
	for _, family := range addressFamilies {
		var network string
		var targets []string
		switch family {
		case familyIPv4:
			network = "udp4"
			targets = v4Targets
		case familyIPv6:
			network = "udp6"
			targets = v6Targets
		default:
			return nil, errNotSupportFamily
		}

		for _, target := range targets {
			klog.V(4).Infof("Looking for default routes with IPv%d addresses via dial %s", uint(family), target)
			ip, err := getLocalIPAddressByDial(network, target)
			if err == nil {
				klog.V(4).Infof("Found active IP %v ", ip)
				return ip, nil
			}
		}
	}
	klog.V(4).Infof("No active IP found by dial")
	return nil, errCannotGetDefaultIP

}

func getLocalIPAddressByDial(network string, testAddr string) (net.IP, error) {

	dnsUDPAddr, err := net.ResolveUDPAddr(network, testAddr)
	if err != nil {
		return nil, err
	}
	c, err := net.DialUDP(network, nil, dnsUDPAddr)
	if err != nil {
		return nil, err
	}
	localIP := c.LocalAddr().(*net.UDPAddr).IP
	_ = c.Close()
	return localIP, nil
}
