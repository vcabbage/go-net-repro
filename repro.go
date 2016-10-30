package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("IN MAIN")
	serverConn, serverAddr, clientConn, clientAddr := testConns()

	errChan := make(chan error)
	go func() {
		_, err := clientConn.WriteTo([]byte{}, serverAddr)
		clientConn.Close()
		errChan <- err
	}()

	_, serverErr := serverConn.WriteTo([]byte{}, clientAddr)
	serverConn.Close()

	clientErr := <-errChan
	if serverErr != nil || clientErr != nil {
		fmt.Println("serverConn:", serverErr)
		fmt.Println("clientConn:", clientErr)
		os.Exit(1)
	}
}

func testConns() (*net.UDPConn, *net.UDPAddr, *net.UDPConn, *net.UDPAddr) {
	clientConn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}
	clientAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: clientConn.LocalAddr().(*net.UDPAddr).Port,
	}

	serverConn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: serverConn.LocalAddr().(*net.UDPAddr).Port,
	}

	return serverConn, serverAddr, clientConn, clientAddr
}
