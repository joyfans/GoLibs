package dict

type JsonDict map[string]interface{};
type JsonArray []interface{};

func toJsonDict(v interface{}) JsonDict {
    switch t := v.(type) {
        case nil:
            return nil

        case JsonDict:
            return t

        case *JsonDict:
            return *t
    }

    return v.(map[string]interface{})
}

func toJsonArray(v interface{}) JsonArray {
    switch t := v.(type) {
        case nil:
            return nil

        case JsonArray:
            return t

        case *JsonArray:
            return *t
    }

    return v.([]interface{})
}

func (self JsonDict) Map(key string) JsonDict {
    return toJsonDict(self[key])
}

func (self JsonDict) Array(key string) JsonArray {
    return toJsonArray(self[key])
}

func (self JsonArray) Map(index int) JsonDict {
    return toJsonDict(self[index])
}

func (self JsonArray) Array(index int) JsonArray {
    return toJsonArray(self[index])
}

func (self *JsonArray) Append(values ...interface{}) *JsonArray {
    *self = append(*self, values...)
    return self
}
