package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/gelian/gldm/src/utils/crc"
)

type MacData struct {
	Mac  string
	Crc8 byte
	Rssi byte
	Next *MacData
}

type MacStat struct {
	TotalCount int
	DiffCount  int
	Macs       *MacData
}

func addMac(ms *MacStat, mac *MacData) {

	ms.TotalCount++
	mac.Crc8 = crc.Crc8([]byte(mac.Mac))
	if ms.Macs == nil {
		ms.Macs = mac
		return
	}

	next := ms.Macs
	for {

		if next.Crc8 == mac.Crc8 {
			if next.Mac == mac.Mac {
				break
			}
		}

		if nil == next.Next {
			next.Next = mac
			ms.DiffCount++
			break
		}

		next = next.Next
	}
}

func statLine(ms *MacStat, l string) {
	regex := `"[0-9a-fA-F]{12}"|"[0-9a-fA-F:-]{17}"`

	reg := regexp.MustCompile(regex)
	matchs := reg.FindAllString(l, -1)

	for _, m := range matchs {
		mac := &MacData{}
		mac.Mac = m
		//fmt.Println("mac:", m)
		addMac(ms, mac)
	}

}

func (ms *MacStat) StatFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
			} else {
				fmt.Println("Read file error!", err)
			}
			break
		}
		statLine(ms, line)
	}
}

func (ms *MacStat) Show() {
	fmt.Println("TotalCount:", ms.TotalCount)
	fmt.Println("DiffCount:", ms.DiffCount)
}
