package zlib

import (
    . "ml/trace"
    "io/ioutil"
    "compress/zlib"
    "bytes"
    "ml/compress"
)

func Compress(data []byte, level int) []byte {
    var buffer bytes.Buffer
    w, err := zlib.NewWriterLevel(&buffer, level)
    if err != nil {
        Raise(compress.NewCompressError("%v", err.Error()))
    }

    defer w.Close()
    w.Write(data)

    return buffer.Bytes()
}

func Decompress(compressed []byte) []byte {
    reader, err := zlib.NewReader(bytes.NewReader(compressed))
    if err != nil {
        Raise(compress.NewDecompressError("%v", err.Error()))
    }

    defer reader.Close()
    decompressed, err := ioutil.ReadAll(reader)
    if err != nil {
        Raise(compress.NewDecompressError("%v", err.Error()))
    }

    return decompressed
}
