package sync2

type Event struct {
    ch chan int
}

func NewEvent() *Event {
    return &Event{ch: make(chan int, 1)}
}

func (self *Event) Wait() {
    <- self.ch
}

func (self *Event) TryWait() bool {
    select {
        case <- self.ch:
            return true

        default:
            return false
    }
}

func (self *Event) Signal() {
    self.trySignal()
}

func (self *Event) Broadcast() {
    for self.trySignal() {}
}

func (self *Event) trySignal() bool {
    select {
        case self.ch <- 0:
            return true

        default:
            return false
    }
}
