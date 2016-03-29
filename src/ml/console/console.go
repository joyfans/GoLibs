package console

import (
    "fmt"
    "bufio"
    "os"
)

func ReadLine() string {
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    return text
}

func SetTitle(text interface{}) {
    setTitle(fmt.Sprintf("%v", text))
}

func Pause(text ...interface{}) {
    if len(text) != 0 {
        fmt.Printf("%v\n", text[0])
    }

    pause()
}
