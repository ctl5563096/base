package str

import (
    "bytes"
    "math/rand"
    "time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// 高性能生成随机字符串
func RandStringBytesMaskImprSrcUnsafe(size int) string {
    source := rand.NewSource(time.Now().UnixNano()) // 产生随机种子
    var s bytes.Buffer
    for i := 0; i < size; i++ {
        s.WriteByte(letterBytes[source.Int63()%int64(len(letterBytes))])
    }
    return s.String()
}

