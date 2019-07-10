package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	intervalTimeout = 1
)

var host = flag.String("host", "0.0.0.0", "host")
var port = flag.String("port", "8221", "port")
var mode = flag.String("mode", "1024", "mode")
var rate = flag.Int("rate", 1000, "rate")
var num = flag.Int("n", 100, "numbers of sending")

func main() {
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var slice512 = []byte("901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012")

	var slice1024 = []byte("90123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234")
	var data []byte
	if *mode == "1024" {
		data = slice1024
	} else {
		data = slice512
	}
	file, err := os.OpenFile("recored.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var i, sum, n int
	var min = int64(99999999)
	var max, sumTime int64
	var avg float64
	buf := make([]byte, 32*1024)
	ch := make(chan struct{})
	go func() {
		for {
			if i >= *num {
				ch <- struct{}{}
				return
			}

			n, err = conn.Read(buf)
			if err != nil {
				//fmt.Println("Client failed to read UDP msg: ", err.Error())
				ch <- struct{}{}
				return
			}
			err = conn.SetDeadline(time.Now().Add(intervalTimeout * time.Second))
			if err != nil {
				ch <- struct{}{}
				return
			}

			tRecv := time.Now().UnixNano()
			tSend := int64(binary.BigEndian.Uint64(buf[:8]))
			tt := (tRecv - tSend) / 1e6
			if min > tt {
				min = tt
			}
			if max < tt {
				max = tt
			}
			sumTime += tt
			i++
			sum += n
			avg = float64(sumTime) / float64(i)
			s := fmt.Sprintf("%v %v %v Mbit/s %v ms %v ms\n", i, sum, float64(sum)*8/30/(1024*1024), tt, avg)
			_, err = file.WriteString(s)
			if err != nil {
				fmt.Println(err)
				ch <- struct{}{}
				return
			}
		}
	}()

	var i2, sum2, n2 int
	t1 := time.NewTimer(0)
	defer t1.Stop()
	//data2 := make([]byte, 1024)
L:
	for {
		select {
		case <-t1.C:
			tS := time.Now().UnixNano()
			//fmt.Println(t1)
			s := make([]byte, 8)
			binary.BigEndian.PutUint64(s, uint64(tS))
			data := append(s, data...)
			n2, err = conn.Write(data)
			//checkError(err)
			if err != nil {
				fmt.Println(err)
				break L
			}
			i2++
			sum2 += n2
			fmt.Println(i2)
			if i2 >= *num {
				fmt.Println("已发送指定包数！")
				break L
			}
			//tE := time.Now().UnixNano()
			//left := 30 - (tE-tS)/1e6
			//fmt.Println("left", left)
			t1.Reset(time.Duration(*rate) * time.Microsecond)
		}
	}

	<-ch
	fmt.Println("退出！")

	fmt.Println("双向发送测试结果client:")
	fmt.Printf(
		"\t设定收发包数: %v\n"+
			"\t发送包个数: %v\n"+
			"\t发送总字节数: %v\n"+
			"\t接收次数: %v\n"+
			"\t接收字节数: %v\n"+
			"\t丢包率: %v %%\n"+
			"\t码率: %v Mbit/s\n"+
			"\t平均时延: %v ms\n\t最小时延: %v\n\t最大时延: %v\n",
		*num,
		i2,
		sum2,
		i,
		sum,
		float64((sum2-sum)/1024)*100/float64(i2),
		float64(sum)*8/30/(1024*1024), avg, min, max)
	fmt.Println("end")
}

/*
func sendData(addr string) {

	fmt.Printf("Client send msg to %s from %s\n", conn.RemoteAddr(), conn.LocalAddr())
	if err != nil {
		fmt.Println("Client failed to write UDP msg: ", err.Error())
		return
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Client failed to read UDP msg: ", err.Error())
		return
	}

	fmt.Printf("Client Get \"%s\" %v\n", string(data[:n]), data[:n])
}
*/
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
