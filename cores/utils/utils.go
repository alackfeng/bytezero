package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// InetAToI - ip to int64.
func InetAToI(ip string) int64 {
	ret := big.NewInt(0)
	nip := net.ParseIP(ip)
	if nip.To4() != nil {
		ret.SetBytes(nip.To4())
	} else {
		ret.SetBytes(nip.To16())
	}
	return ret.Int64()
}

// StringToInt64 -
func StringToInt64(s string) int64 {
	si, _ := strconv.ParseInt(s, 10, 64)
	return si
}

// StringToInt -
func StringToInt(s string) int {
	si, _ := strconv.Atoi(s)
	return si
}

// IntToString -
func IntToString(n int) string {
	return strconv.Itoa(n)
}

// Int64ToString -
func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

// JSONToBytes -
func JSONToBytes(jn interface{}) []byte {
	byt, err := json.Marshal(jn)
	if err != nil {
		return nil
	}
	return byt
}

// BytesToJSON -
func BytesToJSON(byt []byte, jn interface{}) error {
	err := json.Unmarshal(byt, jn)
	if err != nil {
		return err
	}
	return nil
}

// StringToSlice -
func StringToSlice(d string, s string) []string {
	return strings.Split(d, s)
}

// NowMs - 当前时间ms.
func NowMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// TimestampFormat - ms format.
func TimestampFormat(ts int64) string {
	return time.Unix(0, ts*1000000).Format(time.RFC3339Nano)
}

// MsFormat - ms format.
func MsFormat(ts int64) string {
	return time.Unix(0, ts*1000000).Format(time.RFC3339Nano)
}

// MsFormatMonth - ms format.
func MsFormatMonth(ts int64) string {
	currMs := time.Unix(0, ts*1000000)
	return fmt.Sprintf("%04d-%02d-%02d", currMs.Year(), currMs.Month(), currMs.Day())
}

// MsFormatFN - ms format filename.
func MsFormatFN(ts int64) string {
	currMs := time.Unix(0, ts*1000000)
	return fmt.Sprintf("%02d-%02d-%02d", currMs.Year(), currMs.Month(), currMs.Day())
	// return fmt.Sprintf("%02d-%02d-%02d_%02d-%02d-%02d", currMs.Year(), currMs.Month(), currMs.Day(), currMs.Hour(), currMs.Minute(), currMs.Second())
}

// MsFormatDay - ms format day peroid.
func MsFormatDay(ts int64) string {
	currMs := time.Unix(0, ts*1000000)
	var peroid int = 3 // 3 day.
	fmt.Println("day ", currMs.Day())
	day := currMs.Day() / peroid
	return fmt.Sprintf("%02d%02d%02d", currMs.Year(), currMs.Month(), day)
}

// IsExistFile - 文件是否存在
func IsExistFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if ok := os.IsNotExist(err); ok {
		return false, nil
	}
	return false, err
}

// DirIsExistThenMkdir -
func DirIsExistThenMkdir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if ok := os.IsNotExist(err); !ok {
		return err
	}
	// 不存在创建一个多级目录.
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func MapStringToSlice(om map[string]interface{}) (ret []string) {
	for key := range om {
		ret = append(ret, key)
	}
	return ret
}

// GenMD5 - 计算MD5.
func GenMD5(s string) string {
	w := md5.New()
	io.WriteString(w, s)
	return fmt.Sprintf("%X", w.Sum(nil))
}

func Base64Encode(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}

// Base64Decode -
func Base64Decode(s string) []byte {
	r, _ := base64.StdEncoding.DecodeString(s)
	return r
}

// GenNonce -  随机nonce数字.
func GenNonce() int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(1000000000))
	return result.Int64()
}

var numeric = [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

// GenVerifyCode - 随机验证码.
func GenVerifyCode() string {
	// for i := 0; i < len; i++ {
	// 	result, _ := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
	// 	fmt.Fprintf(&sb, "%d", numeric[result.Int64()%10])
	// }
	// return sb.String()

	var sb strings.Builder
	result, _ := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
	fmt.Fprintf(&sb, "%06d", result.Int64())
	return sb.String()[0:6]
}
