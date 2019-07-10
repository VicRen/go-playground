package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	intervalTimeout = 1
)

var host = flag.String("host", "0.0.0.0", "host")
var port = flag.String("port", "8221", "port")
var num = flag.Int("n", 100, "numbers of sending")

func main() {
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Server listening on: ", addr)
	file, err := os.OpenFile("recored.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var i, sum, n, sum2 int
	var min = int64(99999999)
	var max, sumTime int64
	var avg float64
	var remoteAddr *net.UDPAddr
	finishCh := make(chan struct{})
	data := make([]byte, 32*1024)
	go func() {
		for {
			if i >= *num {
				finishCh <- struct{}{}
				return
			}

			n, remoteAddr, err = conn.ReadFromUDP(data)
			if err != nil {
				finishCh <- struct{}{}
				return
			}
			err = conn.SetDeadline(time.Now().Add(intervalTimeout * time.Second))
			if err != nil {
				finishCh <- struct{}{}
				return
			}

			tRecv := time.Now().UnixNano()
			tSend := int64(binary.BigEndian.Uint64(data[:8]))
			tt := (tRecv - tSend) / 1e6
			if min > tt {
				min = tt
			}
			if max < tt {
				max = tt
			}
			fmt.Println("tRecv", tRecv, "tSend", tSend, "diff", tt)
			sumTime += tt
			i++
			sum += n
			avg = float64(sumTime) / float64(i)
			s := fmt.Sprintf("%v %v %v Mbit/s %vms %v ms\n", i, sum, float64(sum)*8/30/(1024*1024), tt, avg)
			_, err = file.WriteString(s)
			if err != nil {
				fmt.Println(err)
				finishCh <- struct{}{}
				return
			}
			raw := data[:n]
			d1 := make([]byte, 1024)
			var d1Len int
			for {
				if len(raw) > 1024 {
					d1Len = 1024
					copy(d1, raw[:1024])
					raw = raw[1024:]
				} else if len(raw) > 1 {
					d1Len = len(raw)
					copy(d1, raw)
					raw = raw[d1Len:]
				} else {
					break
				}
				nw, err := conn.WriteToUDP(d1[:d1Len], remoteAddr)
				if err != nil {
					fmt.Println(err)
					finishCh <- struct{}{}
					return
				}
				sum2 += nw
			}
		}
	}()

	<-finishCh
	fmt.Println("已收到设定的包数！")

	shouldSize := strconv.Itoa(*num * 1024)
	fmt.Println("双向发送测试结果server:")
	fmt.Printf(
		"\t设定收发包数: %v\n"+
			"\t接收次数: %v\n"+
			"\t接收总字节数: %v/"+shouldSize+"\n"+
			"\t丢包率: %v %%\n"+
			"\t码率: %v Mbit/s\n"+
			"\t平均时延: %v ms\n"+
			"\t最小时延: %v\n"+
			"\t最大时延: %v\n",
		*num,
		i,
		sum,
		float64((*num*1024-sum)/1024)*100/float64(*num),
		float64(sum)*8/30/(1024*1024), avg, min, max)
	fmt.Println("发送总字节数：", sum2, "包数", sum2/1024)
	fmt.Println("end")
}

func handleClient(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("Server failed to read UDP msg: ", err.Error())
		return
	}

	fmt.Printf("Server Got \"%s\" from %s\n", string(data[:n]), remoteAddr)
	_, err = conn.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		fmt.Println("Server failed to write UDP msg:", err.Error())
	}
}
