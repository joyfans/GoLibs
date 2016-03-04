package trace

import (
    "fmt"
)

type Exception struct {
    Message     string
    Traceback   string
    Value       interface{}
}

func (self *Exception) String() string {
    return "(traceback)\n" + self.Traceback + "\n" + self.Message
}

func (self *Exception) Error() string {
    return self.String()
}

/*

    exception types

*/

type BaseException struct {
    Message string
}

type IndexError struct {
    *BaseException
}

type AttributeError struct {
    *BaseException
}

type TimeoutError struct {
    *BaseException
}

type KeyError struct {
    *BaseException
}

type NotImplementedError struct {
    *BaseException
}

type FileGenericError struct {
    *BaseException
}

type FileNotFoundError struct {
    *BaseException
}

type PermissionError struct {
    *BaseException
}

func NewBaseException(format string, args ...interface{}) *BaseException {
    return &BaseException{Message: fmt.Sprintf(format, args...)}
}

func NewIndexError(format string, args ...interface{}) *IndexError {
    return &IndexError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewAttributeError(format string, args ...interface{}) *AttributeError {
    return &AttributeError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewTimeoutError(format string, args ...interface{}) *TimeoutError {
    return &TimeoutError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewKeyError(format string, args ...interface{}) *KeyError {
    return &KeyError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewNotImplementedError(format string, args ...interface{}) *NotImplementedError {
    return &NotImplementedError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewFileGenericError(format string, args ...interface{}) *FileGenericError {
    return &FileGenericError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewFileNotFoundError(format string, args ...interface{}) *FileNotFoundError {
    return &FileNotFoundError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewPermissionError(format string, args ...interface{}) *PermissionError {
    return &PermissionError{
        BaseException: NewBaseException(format, args...),
    }
}


func (self *BaseException) String() string {
    return self.Message
}

func (self *BaseException) Error() string {
    return self.String()
}