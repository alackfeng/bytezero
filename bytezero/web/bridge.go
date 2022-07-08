package web

import (
	"fmt"
	"net/url"
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
    URL         string    `form:"URL" json:"URL" xml:"URL" bson:"URL" binding:"required"`
    User        string    `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
    Pass        string    `form:"Pass" json:"Pass" xml:"Pass" bson:"Pass" binding:"required"`
    Expired     int64    `form:"Expired" json:"Expired" xml:"Expired" bson:"Expired" binding:"required"`
}

// Scheme -
func (c CredentialURL) Scheme() string {
    u, err := url.Parse(c.URL)
    if err != nil {
        return ""
    }
    return u.Scheme
}

// Host -
func (c CredentialURL) Host() string {
    u, err := url.Parse(c.URL)
    if err != nil {
        return ""
    }
    return u.Hostname()
}

// Port -
func (c CredentialURL) Port() string {
    u, err := url.Parse(c.URL)
    if err != nil {
        return "0"
    }
    return u.Port()
}

func (c CredentialURL) Secret() bool {
    scheme := c.Scheme()
    return scheme == "tls" || scheme == "wss" || scheme == "https"
}

// String -
func (c CredentialURL) String() string {
    return fmt.Sprintf("URL[%s://%s:%s (%v)] User[%s] Pass[%s] Expired[%d]", c.Scheme(), c.Host(), c.Port(), c.Secret(), c.User, c.Pass, c.Expired)
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
