package web

import (
	"fmt"
)

// CredentialGetReq -
type CredentialGetReq struct {
}

// // CredentialURL -
// type CredentialResult struct {
//     User        string    `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
//     Pass        string    `form:"Pass" json:"Pass" xml:"Pass" bson:"Pass" binding:"required"`
//     Expired     int64    `form:"Expired" json:"Expired" xml:"Expired" bson:"Expired" binding:"required"`
// }

// CredentialURL -
type CredentialURL struct {
    Scheme      string    `form:"Scheme" json:"Scheme" xml:"Scheme" bson:"Scheme" binding:"required"`
    IP          string    `form:"IP" json:"IP" xml:"IP" bson:"IP" binding:"required"`
    Port        string    `form:"Port" json:"Port" xml:"Port" bson:"Port" binding:"required"`
    User        string    `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
    Pass        string    `form:"Pass" json:"Pass" xml:"Pass" bson:"Pass" binding:"required"`
    Expired     int64     `form:"ExpiredMs" json:"ExpiredMs" xml:"ExpiredMs" bson:"ExpiredMs" binding:"required"`
}

func (c CredentialURL) Secret() bool {
    return c.Scheme == "tls" || c.Scheme == "wss" || c.Scheme == "https"
}

// String -
func (c CredentialURL) String() string {
    return fmt.Sprintf("URL[%s://%s:%s (%v)] User[%s] Pass[%s] Expired[%d]", c.Scheme, c.IP, c.Port, c.Secret(), c.User, c.Pass, c.Expired)
}


// CredentialUrlResult -
type CredentialUrlResult []CredentialURL

// Get -
func (c CredentialUrlResult) Get(i int) CredentialURL {
    if i >= len(c) {
        return CredentialURL{}
    }
    return c[i]
}
