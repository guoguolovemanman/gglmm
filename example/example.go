package example

import "github.com/weihongguo/gglmm"

// Example --
type Example struct {
	gglmm.Model
	IntValue    int     `json:"intValue"`
	FloatValue  float64 `json:"floatValue"`
	StringValue string  `json:"stringValue"`
}
