package gglmm

// Status --
type Status struct {
	Value int8   `json:"value"`
	Name  string `json:"name"`
}

var (
	//StatusInvalid --
	StatusInvalid = Status{-128, "无效"}
	// StatusFrozen --
	StatusFrozen = Status{-127, "冻结"}
	// StatusValid --
	StatusValid = Status{1, "有效"}
	// Statuses --
	Statuses = []Status{StatusValid, StatusFrozen, StatusInvalid}
)
