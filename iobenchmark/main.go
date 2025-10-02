package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func benchSocket(iter int) time.Duration {
	socketPath := "/tmp/demo.sock"
	os.Remove(socketPath)
	ln, _ := net.Listen("unix", socketPath)
	defer ln.Close()
	defer os.Remove(socketPath)

	done := make(chan struct{})
	go func() {
		conn, _ := ln.Accept()
		defer conn.Close()
		buf := make([]byte, 16)
		for i := 0; i < iter; i++ {
			conn.Read(buf)
			conn.Write([]byte("ok"))
		}
		close(done)
	}()

	conn, _ := net.Dial("unix", socketPath)
	defer conn.Close()

	start := time.Now()
	for i := 0; i < iter; i++ {
		conn.Write([]byte("x"))
		buf := make([]byte, 16)
		conn.Read(buf)
	}
	<-done
	return time.Since(start)
}

func benchFile(iter int) time.Duration {
	fpath := "/tmp/demo.txt"
	start := time.Now()
	for i := 0; i < iter; i++ {
		os.WriteFile(fpath, []byte("x"), 0o644)
		_, _ = os.ReadFile(fpath)
	}
	return time.Since(start)
}

func main() {
	iters := 10000
	d1 := benchSocket(iters)
	d2 := benchFile(iters)

	fmt.Printf("Socket round-trips: %v for %d ops (avg %v)\n",
		d1, iters, d1/time.Duration(iters))
	fmt.Printf("File writes+reads: %v for %d ops (avg %v)\n",
		d2, iters, d2/time.Duration(iters))
}
