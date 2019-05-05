package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sort"
	"strconv"

	"github.com/zhu/qvm/server/errors"
)

// left close & right open range section
type section struct {
	from uint32
	to   uint32
}

// find an avaiable subnet in the given parent net which is not conflict with existing subnets and mask is parentMask + maskInc
func SearchAvailableSubnet(netStr string, subnetStrs []string, maskInc int) (ipNet *net.IPNet, err error) {
	if maskInc <= 0 || maskInc >= 24 {
		return nil, fmt.Errorf("SearchAvailableSubnet(%v, %v, %d), invalid maskInc value, it should between 0 ~ 24!", netStr, subnetStrs, maskInc)
	}

	// Step 1: transform params
	parentNet, existSubnets, err := transformParams(netStr, subnetStrs)
	if err != nil {
		return nil, fmt.Errorf("transformParams(%v, %v): %v\n", netStr, subnetStrs, err)
	}

	parentMask, _ := parentNet.Mask.Size()
	targetMask := parentMask + maskInc

	if targetMask >= 32 {
		return nil, fmt.Errorf("SearchAvailableSubnet(%v, %v, %d), invalid maskInc value, it should less than 32-parentMask", netStr, subnetStrs, maskInc)
	}

	// step 2: compute used sections, and sort sections by asc and merge to smallest
	sections, err := computeUsedSections(parentMask, targetMask, existSubnets)
	if err != nil {
		return nil, fmt.Errorf("computeUsedSections(%d, %d, %v): %v\n", parentMask, targetMask, existSubnets, err)
	}

	// step 3: find a gap in sections
	max := uint32(1 << uint(maskInc))

	value := uint32(0)
	for i := 0; value < max && i < len(sections); {
		sec := sections[i]

		if value < sec.from {
			break
		}

		i++
		value = sec.to + 1
	}

	// return nil if there are no avaiable gap
	if value >= max {
		return nil, errors.NoGlobalInconlictSubnet
	}

	// step 4: give result
	value = value << uint(32-targetMask)

	parentIpValue := IPv4ToUint32(&parentNet.IP)
	resultIp := Uint32ToIPv4(parentIpValue | value)

	ipNet = &net.IPNet{
		IP:   resultIp,
		Mask: net.CIDRMask(targetMask, 32),
	}

	return
}

// check intersect of two given nets
func intersect(n1, n2 *net.IPNet) bool {
	return n2.Contains(n1.IP) || n1.Contains(n2.IP)
}

// transform ipv4 address to binary string
func toBineryString(ip *net.IP) string {
	var buf bytes.Buffer

	for _, b := range *ip {
		fmt.Fprintf(&buf, "%08b", b)
	}

	return buf.String()
}

// Uint32 to IPv4
func Uint32ToIPv4(i uint32) net.IP {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)

	return net.IPv4(buf[0], buf[1], buf[2], buf[3])
}

// IPv4 to Uint32
func IPv4ToUint32(ip *net.IP) uint32 {
	return binary.BigEndian.Uint32(*ip)
}

func transformParams(netStr string, subnetStrs []string,
) (parentNet *net.IPNet, existSubnets []*net.IPNet, err error) {
	_, parentNet, err = net.ParseCIDR(netStr)
	if err != nil {
		return nil, nil, fmt.Errorf("net.ParseCIDR(%s): %v\n", netStr, err)
	}

	existSubnets = []*net.IPNet{}
	if subnetStrs != nil && len(subnetStrs) > 0 {
		for _, str := range subnetStrs {
			_, ipNet, err := net.ParseCIDR(str)
			if err != nil {
				return nil, nil, fmt.Errorf("net.ParseCIDR(%s): %v\n", str, err)
			}

			if intersect(parentNet, ipNet) {
				existSubnets = append(existSubnets, ipNet)
			}
		}
	}

	return
}

// compute used sections, and sort sections by asc and merge to smallest
func computeUsedSections(parentMask, targetMask int, existSubnets []*net.IPNet) (mergedSecs []section, err error) {

	sections := []section{}

	if existSubnets == nil || len(existSubnets) == 0 {
		return sections, nil
	}

	for _, subnet := range existSubnets {
		ipBineryStr := toBineryString(&subnet.IP)
		mask, _ := subnet.Mask.Size()

		from, err := strconv.ParseUint(ipBineryStr[parentMask:targetMask], 2, 64)
		if err != nil {
			return nil, fmt.Errorf("strconv.ParseUint(%s, %d, %d): %v\n", ipBineryStr[0:targetMask-parentMask], 2, 64, err)
		}

		to := from

		if mask < targetMask {
			diff := targetMask - mask
			to = 1<<uint(diff) + from - 1
		}

		sections = append(sections, section{
			from: uint32(from),
			to:   uint32(to),
		})
	}

	// sort by asc
	sort.SliceStable(sections, func(i, j int) bool {
		return sections[i].from < sections[j].from
	})

	// merge
	mergedSecs = []section{sections[0]}
	for _, sec := range sections {
		pre := mergedSecs[len(mergedSecs)-1]

		if sec.from > pre.to {
			mergedSecs = append(mergedSecs, sec)
			continue
		}

		if sec.to > pre.to {
			pre.to = sec.to
		}
	}

	return
}
