package tool

import (
	"bytes"
	"fmt"
	"os"
)

// AppsTool -
type AppsTool struct {

}

// NewAppsTool -
func NewAppsTool() *AppsTool {
    return &AppsTool{}
}

func (app *AppsTool) diffFile(f1, f2 string) error {
    fmt.Println("--------------diffFile begin.")
    file1, _ := os.Open(f1)
    defer file1.Close()
    file2, _ := os.Open(f2)
    defer file2.Close()
    fileinfo, _ := file1.Stat()

    buffer1 := make([]byte, 1024*64)
    buffer2 := make([]byte, 1024*64)
    buffer3 := make([]byte, 1024*64)
    buffer4 := make([]byte, 1024*64)
    index := 0
    for {
        n1, _ := file1.Read(buffer1)
        n2, _ := file2.Read(buffer2)
        if n1 == 0 || n2 == 0 {
            break
        }

        if index / (1024*64) == 3 {
            copy(buffer3, buffer2[0:n2])
        }
        if index / (1024*64) == 4 {
            copy(buffer4, buffer1[0:n1])
        }
        if !bytes.Equal(buffer1[0:n1], buffer2[0:n2]) {
            fmt.Println("--------------index: ", index, index / (1024*64), index % (1024*64), fileinfo.Size(), fileinfo.Size() /(1024*64),  " error")
            count := 10
            for i:=0; i<n1; i++ {
                if buffer1[i] != buffer2[i] {
                    fmt.Printf("%d, ", i)
                    if count <= 0 {
                        break
                    }
                    count--
                }
            }
            fmt.Println("--------------buffer1: ", buffer1[0:300])
            fmt.Println("--------------buffer2: ", buffer2[0:300])
        }
        index += n1
    }
    fmt.Println("--------------buffer3: ", buffer3[0:300])
    fmt.Println("--------------buffer4: ", buffer4[0:300])
    fmt.Println("\n-------------------\n")
    for i:=5*1024; i<10*1024; i++ {
        if buffer3[i] != buffer4[i] {
            fmt.Printf("%d, ", i)
        }
    }
    fmt.Println("\n--------------diffFile end...", fileinfo.Size())
    return nil
}

// Main -
func (app *AppsTool) Main() {
    // app.diffFile("C:\\Users\\Administrator\\Desktop\\Red_Bells.jpg", "C:\\Users\\Administrator\\Desktop\\Red_Bells_bad.jpg")
    app.diffFile("C:\\Users\\Administrator\\Desktop\\2020-11-24143924--3", "C:\\Users\\Administrator\\Desktop\\2020-11-24143924")
}
