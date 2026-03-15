package dto

type CodeRequest struct {
	Code          string `json:"code"`
	Percent       int    `json:"percent"`
	AmountOfUsage int    `json:"amount_of_usage"`
	ExpiredAt     string `json:"expired_at"`
}

type CodeResponse struct {
	Code          string `json:"code"`
	Percent       int    `json:"percent"`
	AmountOfUsage int    `json:"amount_of_usage,omitempty"`
	ExpiredAt     string `json:"expired_at,omitempty"`
}

type CodeListResponse struct {
	Codes []CodeResponse `json:"codes"`
}
