package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/go-uuid"
)

// Credential -
type Credential struct {
    Name string
    ExpireS int64
}

// NewCredential -
func NewCredential(expiredMs int64) *Credential {
    name, _ := uuid.GenerateRandomBytes(12)
    return &Credential{
        Name: fmt.Sprintf("%x", name),
        ExpireS: expiredMs / 1000,
    }
}

// Username -
func (c *Credential) Username() string {
    return fmt.Sprintf("%d:%s", c.ExpireS, c.Name)
}

// Sign - sha1-hmac.
func (c *Credential) Sign(pass string) (s string) {
    mac := hmac.New(sha1.New, []byte(pass))
    mac.Write([]byte(c.Username()))
    return base64.StdEncoding.EncodeToString(mac.Sum([]byte("")))
}


var ErrCredentialParamNotEngouth = errors.New("Credential Sign Param Not Engouth")
var ErrCredentialUsernameNULL = errors.New("Credential Username is nil")
var ErrCredentialSignNotMatch = errors.New("Credential Sign Not Match")
var ErrCredentialExpired = errors.New("User %s Sign<%s> expired now.")

// CredentialVerify -
func CredentialVerify(user string, sign string, pass string, get_fileds_sign func(string)[]byte) error {
    ss := StringToSlice(user, ":")
    if len(ss) != 3 {
        return ErrCredentialParamNotEngouth
    }
    cred := &Credential{
        Name: ss[1],
        ExpireS: StringToInt64(ss[0]),
    }
    if cred.Name == "" {
        return ErrCredentialUsernameNULL
    }
    if cred.ExpireS < NowMs() / 1000 {
        return fmt.Errorf(ErrCredentialExpired.Error(), cred.Name, user)
    }
    secretType := ss[2]
    secretKey := cred.Sign(pass)

    m5 := Md5Sign([]byte(user), get_fileds_sign(secretType), []byte(secretKey))
    if m5 != sign {
        return ErrCredentialSignNotMatch
    }
    return nil
}

// Md5Sign -
func Md5Sign(user []byte, fields []byte, pass []byte) string {
    hs := md5.New()
    hs.Write(user)
    hs.Write(fields)
    hs.Write(pass)
    m5 := fmt.Sprintf("%X", hs.Sum(nil)[4:12])
    return m5
}

