package web

// CredentialGetReq -
type CredentialGetReq struct {
}

// CredentialResult -
type CredentialResult struct {
    User        string    `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
    Pass        string    `form:"Pass" json:"Pass" xml:"Pass" bson:"Pass" binding:"required"`
    Expired     int64    `form:"Expired" json:"Expired" xml:"Expired" bson:"Expired" binding:"required"`
}
