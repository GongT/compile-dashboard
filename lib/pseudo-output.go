package lib

import (
	"time"
	"encoding/hex"
	"crypto/rand"
)

type PseudoOutput struct {
	Tunnel chan string
	ticker *time.Ticker
}

func NewPseudoOutput() *PseudoOutput {
	ret := new(PseudoOutput)
	ret.Tunnel = make(chan string)
	return ret
}

func (po *PseudoOutput) Stop() {
	if po.ticker == nil {
		return
	}
	po.ticker.Stop()
}
func (po *PseudoOutput) Start() {
	if po.ticker != nil {
		return
	}
	po.ticker = time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-po.ticker.C:
				by := make([]byte, 32)
				rand.Read(by)
				po.Tunnel <- ">" + hex.EncodeToString(by) + "[[" + string(0x1B) + "]]\n"
			case <-po.Tunnel:
				po.ticker.Stop()
				return
			}
		}
	}()
}
