package main

import (
    "code.google.com/p/go.crypto/ssh"
    // ...
)

// ...

func main() {
    cmd := os.Args[1] // the first argument is a command we’ll execute on all servers 
    hosts := os.Args[2:] // other arguments (starting from the second one) – the list of servers 
    results := make(chan string, 10) // we’ll write results into the buffered channel of strings
    timeout := time.After(5 * time.Second) // in 5 seconds the message will come to timeout channel

    // initialize the structure with the configuration for ssh packat.
    // makeKeyring() function will be written later
    config := &ssh.ClientConfig{
        User: os.Getenv("LOGNAME"),
        Auth: []ssh.ClientAuth{makeKeyring()},
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
