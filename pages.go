package vocode

type Pages struct {
	Page        int  `json:"page"`
	Size        int  `json:"size"`
	HasMore     bool `json:"has_more"`
	Total       int  `json:"total"`
	IsEstimated bool `json:"total_is_estimated"`
}
