package stcp

import (
	"fmt"
	"io"
	"time"
)

func heartbeats(typ uint8, resetC chan time.Duration, w io.Writer) {
	defer func() {
		fmt.Println("exit heartbeat")
	}()

	var interval time.Duration = HeartbeatTimeout

	hb := []byte{0, 1, typ} // binary data
	timer := time.NewTimer(interval)
	for {

		select {
		case <-timer.C:
			_, err := w.Write(hb) //c.WriteBytes(&hb)
			if err != nil {
				return
			}
		case newInterval, ok := <-resetC:
			if !timer.Stop() {
				<-timer.C
			}
			if !ok {
				return
			}
			interval = newInterval
		}
		_ = timer.Reset(interval)
	}
}
