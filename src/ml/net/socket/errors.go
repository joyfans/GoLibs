package socket

import (
    . "ml/trace"
    . "fmt"
    "net"
)

type SocketError struct {
    *BaseException
}

type SocketTimeoutError struct {
    *SocketError
}

func RaiseSocketError(err error) {
    if err == nil {
        return
    }

    if e, ok := err.(*net.OpError); ok {
        if e.Timeout() {
            Raise(NewSocketTimeoutError(e.Error()))
        }
    }

    Raise(NewSocketError(err.Error()))
}

func NewSocketError(msg string) *SocketError {
    return &SocketError{
        BaseException: NewBaseException(msg),
    }
}

func NewSocketTimeoutError(msg string) *SocketTimeoutError {
    return &SocketTimeoutError{
        SocketError: NewSocketError(msg),
    }
}
