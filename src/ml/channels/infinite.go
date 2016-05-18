package channels

import (
    "fmt"

    "time"

    "ml/array"
    "ml/sync2"
)

type InfiniteChannel struct {
    input   chan interface{}
    output  chan interface{}
    buffer *array.Array
    event  *sync2.Event
}

func NewInfiniteChannel() *InfiniteChannel {
    ch := &InfiniteChannel{
                input   : make(chan interface{}),
                output  : make(chan interface{}),
                buffer  : array.NewArray(),
                event   : sync2.NewEvent(),
        }

    go ch.infiniteBuffer()
    return ch
}

func (self *InfiniteChannel) In() chan <- interface{} {
    return self.input
}

func (self *InfiniteChannel) Out() <-chan interface{} {
    return self.output
}

func (self *InfiniteChannel) Length() int {
    return self.buffer.Length()
}

func (self *InfiniteChannel) Close() {
    fmt.Println("close input")
    close(self.input)
    fmt.Println("wait loop")
    self.event.Wait()
}

func (self *InfiniteChannel) shutdown() {
    fmt.Println("shutdown")
FLUSH:
    for _, v := range *self.buffer {
        select {
            case self.output <- v:

            default:
                break FLUSH
        }
    }

    fmt.Println("close output")
    close(self.output)
}

func (self *InfiniteChannel) infiniteBuffer() {

INFINITE_LOOP:
    for {
        time.Sleep(time.Millisecond)

        switch self.buffer.Length() {
            case 0:
                select {
                    case elem, open := <-self.input:
                        if open == false {
                            break INFINITE_LOOP
                        }
                        self.buffer.Append(elem)

                    default:
                }

            default:
                select {
                    case elem, open := <-self.input:
                        if open == false {
                            break INFINITE_LOOP
                        }
                        self.buffer.Append(elem)

                    case self.output <- self.buffer.Peek(0):
                        self.buffer.Pop(0)
                }
        }
    }

    self.shutdown()
    fmt.Println("Broadcast")
    self.event.Broadcast()
}
