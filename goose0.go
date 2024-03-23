package main

import (
"strconv"
"os"
"log"
"fmt"
"net"
"time"
)

func Portscan(ip string) (ports [1024]int) {
	fmt.Println("\n")
	startPort := 1
	endPort := 1024
	for loopPort := startPort; loopPort <= endPort; loopPort++ {
		address := fmt.Sprintf("%s:%d", ip, loopPort)
		conn, err := net.DialTimeout("tcp", address, time.Millisecond*100)
		if err == nil {
			conn.Close()
			ports[loopPort] = loopPort
		}
	}
	return ports
}

func Ipscan() (iplist [100]net.IP) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for counter, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil || ipnet.IP.To16 != nil {
				iplist[counter] = ipnet.IP.To16()
			}
		}
	}
	return iplist
}

func ScanNetwork() {
	//lets test out this bitch 🔪
	IPlist := Ipscan()
	type IPWithPorts struct {
		IP string
		openPorts [1024]int
	}
	var allIPsWithPorts = [len(IPlist)]IPWithPorts{}
	for i := range IPlist {
		ports := Portscan(IPlist[i].String())
		for p := range ports {
			if ports[p] != 3000 {
				allIPsWithPorts[i].openPorts = ports
			}
		}
		allIPsWithPorts[i].IP = IPlist[i].String()
		fmt.Println("We're on ", i, "of ", len(IPlist))
	}
	
	var lines = []string{
		"IPs with their open ports:",
	}
	for a := range allIPsWithPorts {
		lines = append(lines, allIPsWithPorts[a].IP)
		for range allIPsWithPorts[a].openPorts {
			if allIPsWithPorts[a].openPorts[a] != 0 {
				lines = append(lines, strconv.Itoa(allIPsWithPorts[a].openPorts[a]))
			}
		}
		//lines = append(lines, "\n")
	}

	//lets write them out to a file now
	f, err := os.Create("iplist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	for _, line := range lines {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Open ports logged to file. Go get em' killer.")
}

func main() {
	//adds some FLAVOR to the startup
	ascii := "✩░▒▓▆▅▃▂▁𝐆𝐨𝐨𝐬𝐞 𝐯𝟏▁▂▃▅▆▓▒░✩"
	fmt.Println(ascii)

	ScanNetwork()
}
