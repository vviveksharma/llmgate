package usage

import "time"

type Record struct {
	TenantID  string    `json:"tenant_id"`
	Model     string    `json:"model"`
	TokensIn  int64     `json:"tokens_in"`
	TokensOut int64     `json:"tokens_out"`
	At        time.Time `json:"at"`
}
