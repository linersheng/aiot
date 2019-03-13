package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	SERVER_URL          = "www.youserver.com"
	SERVER_IP           = 8080
	SEND_COUNT          = 10000
	READ_TIMEOUT_SECOND = 3
)

//869060030115533
var dev_id uint64 = 869060330000000
var fail_count uint32 = 0
var logger *log.Logger
var finish_count uint32 = 0

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func sendMessage(udpConn *net.UDPConn, id uint64, c chan int) {
	data := fmt.Sprintf("%d,%d", dev_id+id, 32)
	udpConn.Write([]byte(data))
	//fmt.Println("send:", string(data))
	c <- 1
}

func readMessage(udpConn *net.UDPConn, id uint64, c chan int) {
	buff := make([]byte, 64)
	n, _, err := udpConn.ReadFromUDP(buff)
	if nil != err || n <= 0 {
		logger.Println("recv_err:", id, err)
		atomic.AddUint32(&fail_count, 1)
		c <- 2
		return
	}
	c <- 2
	logger.Printf("recv_ok:%d, %d, %s\n", id, n, string(buff[:n]))

}

func startUdpClient(id uint64) {
	udpAddr, err := net.ResolveUDPAddr("udp4", SERVER_URL+":"+strconv.Itoa(SERVER_IP))

	if err != nil {
		fmt.Println(err)
		//atomic.AddUint32(&finish_count, 1)
		return
	}

	//udp连接
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		//atomic.AddUint32(&finish_count, 1)
		return
	}
	defer udpConn.Close()

	wrChan := make(chan int)
	go sendMessage(udpConn, id, wrChan)
	<-wrChan

	t := time.Now()
	udpConn.SetReadDeadline(t.Add(time.Duration(READ_TIMEOUT_SECOND * time.Second)))
	go readMessage(udpConn, id, wrChan)
	<-wrChan

	atomic.AddUint32(&finish_count, 1)
	close(wrChan)
	//atomic.AddUint32(&finish_count, 1)

}

func main() {

	file, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE, 666)
	if nil != err {
		fmt.Println("Open file err.")
		os.Exit(2)
	}
	defer file.Close()
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	var i uint64 = 0
	for i = 0; i < SEND_COUNT; i++ {
		go startUdpClient(i)
		//time.Sleep(time.Millisecond * 10)
	}

	i = 0
	for {
		if finish_count >= SEND_COUNT {
			return
		}
		time.Sleep(time.Second * 1)
		i++
		fmt.Printf("wait second:%d, finish_count:%d, fail_count:%d\n", i, finish_count, fail_count)

	}

}
