package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
)

// Credential -
type Credential struct {
    Name string
    ExpireS int64
}

// NewCredential -
func NewCredential(name string, expiredMs int64) *Credential {
    return &Credential{
        Name: name,
        ExpireS: expiredMs / 1000,
    }
}

// Username -
func (c *Credential) Username() string {
    return fmt.Sprintf("%d:%s", c.ExpireS, c.Name)
}

// Sign -
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
func CredentialVerify(s string, pass string) error {
    ss := StringToSlice(s, ":")
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
        return fmt.Errorf(ErrCredentialExpired.Error(), cred.Name, s)
    }
    if cred.Sign(pass) != ss[2] {
        return ErrCredentialSignNotMatch
    }
    return nil
}

