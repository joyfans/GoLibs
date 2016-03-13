package json

import (
    . "ml/trace"
    . "ml/dict"

    "os"
    "io/ioutil"
)

func MustMarshal(v interface{}) []byte {
    data, err := Marshal(v)
    RaiseIf(err)
    return data
}

func MustMarshalIndent(v interface{}, prefix, indent string) []byte {
    data, err := MarshalIndent(v, prefix, indent)
    RaiseIf(err)
    return data
}

func loadFile(file string, v interface{}) {
    f, err := os.Open(file)
    if err != nil {
        Raise(NewFileNotFoundError(err.Error()))
    }

    defer f.Close()

    data, err := ioutil.ReadAll(f)
    if err != nil {
        Raise(NewBaseException(err.Error()))
    }

    loadData(data, v)
}

func loadData(data []byte, v interface{}) {
    err := Unmarshal(data, v)
    if err != nil {
        Raise(NewJSONDecodeError(err.Error()))
    }
}

func loadString(text string, v interface{}) {
    loadData([]byte(text), v)
}


func LoadFileDict(file string) (v JsonDict) {
    loadFile(file, &v)
    return
}

func LoadDataDict(data []byte) (v JsonDict) {
    loadData(data, &v)
    return
}

func LoadStringDict(text string) (v JsonDict) {
    loadString(text, &v)
    return
}

func LoadFileArray(file string) (v JsonArray) {
    loadFile(file, &v)
    return
}

func LoadDataArray(data []byte) (v JsonArray) {
    loadData(data, &v)
    return
}

func LoadStringArray(text string) (v JsonArray) {
    loadString(text, &v)
    return
}
