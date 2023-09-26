package web

// action定义.
const (
	ActionResultTypeJSON = "JSON"
	ActionResultTypeHTML = "HTML"
)

// ResultMap -
type ResultMap map[string]interface{}

// ActionResult -
type ActionResult struct {
	Code    ErrCode     `form:"code" json:"code" xml:"code" bson:"code" binding:"required"`
	Message string      `form:"info" json:"info" xml:"info" bson:"info" binding:"required"`
	Data    interface{} `form:"data" json:"data" xml:"data" bson:"data" binding:"required"`
}

// ResultNONE -
type ResultNONE struct{}

// NewActionResult -
func NewActionResult() *ActionResult {
	return &ActionResult{
		Code:    ErrCode_success,
		Message: "ok",
		Data:    ResultNONE{},
	}
}

// Set -
func (ar *ActionResult) Set(code ErrCode, message string) *ActionResult {
	ar.Code = code
	ar.Message = message
	return ar
}

// SetData -
func (ar *ActionResult) SetData(data interface{}) *ActionResult {
	ar.Data = data
	return ar
}
