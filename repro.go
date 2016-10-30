package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	serverConn, serverAddr, clientConn, clientAddr := testConns()
	defer clientConn.Close()
	defer serverConn.Close()

	errChan := make(chan error)
	go func() {
		_, err := clientConn.WriteTo([]byte{}, serverAddr)
		errChan <- err
	}()

	_, serverErr := serverConn.WriteTo([]byte{}, clientAddr)

	clientErr := <-errChan
	if serverErr != nil || clientErr != nil {
		fmt.Println("serverConn:", serverErr)
		spew.Dump(serverConn)
		fmt.Println("clientConn:", clientErr)
		spew.Dump(clientConn)
		os.Exit(1)
	}
}

func testConns() (*net.UDPConn, *net.UDPAddr, *net.UDPConn, *net.UDPAddr) {
	clientConn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}

	clientPort := clientConn.LocalAddr().(*net.UDPAddr).Port
	clientAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:"+strconv.Itoa(clientPort))
	if err != nil {
		panic(err)
	}

	serverConn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}

	serverPort := serverConn.LocalAddr().(*net.UDPAddr).Port
	serverAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:"+strconv.Itoa(serverPort))
	if err != nil {
		panic(err)
	}

	return serverConn, serverAddr, clientConn, clientAddr
}
