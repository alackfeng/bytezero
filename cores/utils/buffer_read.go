package utils

import "fmt"

// BufferRead -
type BufferRead struct {
    buffer []byte
    currOffset int
    readOffset int
    remainLen int
    needRead bool
    maxBufferLen int
}

// NewBufferRead -
func NewBufferRead(l int) *BufferRead {
    return &BufferRead{
        maxBufferLen: l,
        buffer: make([]byte, l),
        needRead: true,
    }
}

// Get -
func (b *BufferRead) Get() []byte {
    if b.currOffset >= b.readOffset {
        fmt.Println("BufferRead::Get - nil.")
        return []byte{}
    }
    fmt.Println("BufferRead::Get - ", b.currOffset, b.readOffset, len(b.buffer[b.currOffset:b.readOffset]))
    return b.buffer[b.currOffset:b.readOffset]
}

// Empty -
func (b *BufferRead) Empty() bool {
    return b.remainLen == 0
}

// Length -
func (b *BufferRead) Length() int {
    return b.remainLen
}

// Read -
func (b *BufferRead) Read(readFunc func([]byte) (int, error)) (int, error) {
    if !b.needRead {
       fmt.Printf("BufferRead::read - currOffset %d, readOffset %d, remainLen %d, cap %d, engouht----------------.\n", b.currOffset, b.readOffset, b.remainLen, len(b.buffer))
       return b.remainLen, nil
    }
    if readFunc != nil {
        l, err := readFunc(b.buffer[b.readOffset:])
        if err != nil {
            fmt.Printf("BufferRead::read - currOffset %d, readOffset %d, remainLen %d, curr read len %d, error----------------.\n", b.currOffset, b.readOffset, b.remainLen, l)
            return b.remainLen, err
        }
        b.readOffset += l
        b.remainLen = b.readOffset - b.currOffset
        if b.readOffset > b.maxBufferLen - b.maxBufferLen / 10 {
            copy(b.buffer, b.buffer[b.currOffset:b.readOffset])
            b.currOffset = 0; b.readOffset = b.remainLen
        }
        b.needRead = true
        fmt.Printf("BufferRead::read - currOffset %d, readOffset %d, remainLen %d, curr read len %d, cap %d.\n", b.currOffset, b.readOffset, b.remainLen, l, len(b.buffer))
        return b.remainLen, nil
    }
    return 0, fmt.Errorf("BufferRead need read func")
}

// Next -
func (b *BufferRead) Next(l int) {
    b.currOffset += l
    b.remainLen = b.readOffset - b.currOffset
    b.needRead = b.remainLen == 0
}

// Step -
func (b *BufferRead) Step() {
    b.needRead = true
}
