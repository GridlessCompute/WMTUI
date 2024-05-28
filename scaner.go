package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type IPRange struct {
	StartIP string
	EndIP   string
	StartO3 int
	Length  int
}

func Ipgen() []int {
	var ips []int
	for i := 1; i <= 255; i++ {
		ips = append(ips, i)
	}

	return ips
}

func GenerateRangeIPs(r IPRange) []string {
	// Something about the O3 and the ipgen
	//
	var ips []string
	split := strings.Split(r.StartIP, ".")

	base := fmt.Sprint(split[0] + "." + split[1] + ".")

	jumps := (r.Length - 255) / 255

	subs := Ipgen()

	for i := r.StartO3; i < (r.StartO3 + jumps + 1); i++ {
		O3 := strconv.Itoa(i)
		for _, ip := range subs {
			O4 := strconv.Itoa(ip)

			ips = append(ips, fmt.Sprint(base+O3+"."+O4))
		}
	}
	return ips
}

func GenIpRange(r FarmStruct) IPRange {
	var iprange IPRange

	start := strings.Split(r.Start, ".")
	end := strings.Split(r.End, ".")

	iprange.StartIP = r.Start
	iprange.EndIP = r.End
	iprange.StartO3, _ = strconv.Atoi(start[2])

	fmt.Println(iprange.StartO3)
	a, _ := strconv.Atoi(end[2])

	length := a - iprange.StartO3

	iprange.Length = (length * 255) + 255

	return iprange
}

func ScanRange(ips []string, wg *sync.WaitGroup, hashChannel chan MinerObj) {
	for _, ip := range ips {
		wg.Add(1)
		go InitScanOne(ip, hashChannel, wg)
	}
}

func PopulateRange(minerList chan MinerObj, wg *sync.WaitGroup, hashChannel chan MinerObj) {
	for m := range minerList {
		wg.Add(1)
		go GetMinerData(wg, m, hashChannel)
	}
}

func PopulateRangeSlice(minerList []MinerObj, wg *sync.WaitGroup, hashChannel chan MinerObj) {
	for _, m := range minerList {
		wg.Add(1)
		go GetMinerData(wg, m, hashChannel)
	}
}

func ScanMaster(scanwg *sync.WaitGroup, popwg *sync.WaitGroup, mainChannel chan MinerObj, swg *sync.WaitGroup, f FarmStruct) {
	mnrOChannel := make(chan MinerObj, 150)
	popChannel := make(chan MinerObj, 150)

	defer swg.Done()

	ips := GenerateRangeIPs(GenIpRange(f))

	ScanRange(ips, scanwg, mnrOChannel)
	scanwg.Wait()
	close(mnrOChannel)

	PopulateRange(mnrOChannel, popwg, popChannel)
	popwg.Wait()
	close(popChannel)

	for mnr := range popChannel {
		//
		//fmt.Println(mnr)
		mainChannel <- mnr
	}

}
