package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"time"
)

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	conn, _ := ssh.Dial("tcp", hostname+":22", config)
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	return hostname + ": " + stdoutBuf.String()
}

func GetPublicKey(privateKeyPath string) (pub ssh.Signer, err error) {
	priv, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	pub, err = ssh.ParsePrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func main() {
	user := os.Args[1]
	privateKeyPath := os.Args[2]
	pub, err := GetPublicKey(privateKeyPath)
	if err != nil {
		panic(err)
	}
	cmd := os.Args[3]                      // the first argument is a command we will execute on all servers
	hosts := os.Args[4:]                   // other arguments (starting from the second one) ? the list of servers
	results := make(chan string, 10)       // we’ll write results into the buffered channel of strings
	timeout := time.After(5 * time.Second) // in 5 seconds the message will come to timeout channel

	// initialize the structure with the configuration for ssh packat.
	// makeKeyring() function will be written later
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pub),
		},
	}

	// running one goroutine (light-weight alternative of OS thread) per server,
	// executeCmd() function will be written later
	for _, hostname := range hosts {
		go func(hostname string) {
			results <- executeCmd(cmd, hostname, config)
		}(hostname)
	}

	// collect results from all the servers or print "Timed out",
	// if the total execution time has expired
	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}
