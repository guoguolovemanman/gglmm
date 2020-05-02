package gglmm

var (
	//StatusInvalid --
	StatusInvalid = ConfigInt8{Value: -128, Name: "无效"}
	// StatusFrozen --
	StatusFrozen = ConfigInt8{Value: -127, Name: "冻结"}
	// StatusValid --
	StatusValid = ConfigInt8{Value: 1, Name: "有效"}
	// Statuses --
	Statuses = []ConfigInt8{StatusValid, StatusFrozen, StatusInvalid}
)
