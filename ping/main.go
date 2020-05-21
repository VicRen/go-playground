package main

import (
	"fmt"
	"log"
	"time"

	ping "github.com/caucy/batch_ping"
)

func main() {
	var ipSlice []string
	// ip list should not more than 65535

	ipSlice = append(ipSlice, "2400:da00:2::29") //support ipv6
	ipSlice = append(ipSlice, "baidu.com")
	for i := 0; i < 256; i++ {
		ipSlice = append(ipSlice, fmt.Sprintf("47.102.45.%d", i))
	}

	bp, err := ping.NewBatchPinger(ipSlice, false) // true will need to be root

	if err != nil {
		log.Fatalf("new batch ping err %v", err)
	}
	bp.SetDebug(false) // debug == true will fmt debug log

	bp.SetSource("") // if hava multi source ip, can use one isp
	bp.SetCount(10)

	bp.OnFinish = func(stMap map[string]*ping.Statistics) {
		count := 0
		for ip, st := range stMap {
			count++
			log.Printf("\n--- %s ping statistics ---\n", st.Addr)
			log.Printf("[%d]: ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n", count, ip,
				st.PacketsSent, st.PacketsRecv, st.PacketLoss)
			log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				st.MinRtt, st.AvgRtt, st.MaxRtt, st.StdDevRtt)
			log.Printf("rtts is %v \n", st.Rtts)
		}

	}

	timeStart := time.Now()
	err = bp.Run()
	if err != nil {
		log.Printf("run err %v \n", err)
	}
	log.Printf("cost: %s", time.Since(timeStart))
	bp.OnFinish(bp.Statistics())
}
