package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	help    = flag.String("h", "", "Displays usage of scanner")
	host    = flag.String("t", "", "Host or IP address to scan")
	minPort = flag.Int("min", 1, "Port from which scanning begins")
	maxPort = flag.Int("max", 1024, "Port in which scanning finished")
	verbose = flag.Bool("V", false, "Verbose: immediately let the user know if the current port is open")
)

func printUsage() {
	log.Println("Usage:")
	log.Println()
	log.Println("	go run scanner.go <host> [OPTIONS]")
}

func testTCPConnection(host string, port int, doneChannel chan bool) {
	timeoutLength := 5 * time.Second

	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), timeoutLength)

	if err != nil {
		doneChannel <- false
		return
	}
	conn.Close()
	log.Printf("[+] %d open", port)
	doneChannel <- true
}

func main() {
	flag.Parse()
	if *help != "" {
		printUsage()
		os.Exit(1)
	}
	if *host == "" {
		log.Println("No target host provided")
		printUsage()
		os.Exit(1)
	}
	doneChannel := make(chan bool)
	activeThreadCount := 0

	fmt.Println("Scanning host: " + *host)

	for portNumber := *minPort; portNumber <= *maxPort; portNumber++ {
		activeThreadCount++
		if *verbose {
			log.Println("Listening on port: " + strconv.Itoa(portNumber))
		}
		go testTCPConnection(*host, portNumber, doneChannel)
	}

	for {
		<-doneChannel
		activeThreadCount--
		if activeThreadCount == 0 {
			break
		}
	}
	log.Println("Scan completed")
}
