package io2

import (
    "ml/io2/filestream"
)

func ReadContent(path string) []byte {
    fs := filestream.Open(path)
    defer fs.Close()

    return fs.ReadAll()
}
