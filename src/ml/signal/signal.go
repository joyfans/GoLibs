package signal

import (
    "os"
    "os/signal"
)

type Handler struct {
    ch chan os.Signal
}

func Notify(handler func(), sig ...os.Signal) *Handler {
    h := &Handler{
        ch: make(chan os.Signal, 10),
    }

    signal.Notify(h.ch, sig...)

    go func() {
        for {
            _, ok := <-h.ch
            if ok == false {
                break
            }

            handler()
        }
    }()

    return h
}

func (self *Handler) Close() {
    close(self.ch)
}
