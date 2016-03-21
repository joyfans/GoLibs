package compress

import (
    . "ml/trace"
)

type CompressError struct {
    *BaseException
}

type DecompressError struct {
    *BaseException
}

func NewCompressError(format string, args ...interface{}) *CompressError {
    return &CompressError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewDecompressError(format string, args ...interface{}) *DecompressError {
    return &DecompressError{
        BaseException: NewBaseException(format, args...),
    }
}
