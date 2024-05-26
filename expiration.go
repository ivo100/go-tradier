package tradier

import (
	"github.com/timpalpant/go-tradier/config"
	"time"
)

type ExpirationInfo struct {
	Expirations Expirations `json:"expirations"`
}

type Expirations struct {
	Expiration []ExpirationElement `json:"expiration"`
}

type ExpirationElement struct {
	Date           string  `json:"date"`
	ExpirationType string  `json:"expiration_type"`
	Strikes        Strikes `json:"strikes"`
	//ContractSize   int64   `json:"contract_size"`
}

type Strikes struct {
	Strike []float64 `json:"strike"`
}

const (
	SandboxEndpoint   = "https://sandbox.tradier.com"
	ApiKeyName        = "TRADIER_KEY"
	ApiSandboxKeyName = "TRADIER_KEY_SB"
	AcntName          = "TRADIER_ACNT"
	SandboxAcntName   = "TRADIER_ACNT_SB"
)

// APIEndpoint       = "https://api.tradier.com"
// StreamEndpoint    = "https://stream.tradier.com"

func NewTradierBroker() *Client {
	cfg := config.Instance()
	//util.GrepEnv("TRADIER")
	if cfg.Sandbox {
		token := config.MustGetFromEnv(ApiSandboxKeyName)
		acnt := config.MustGetFromEnv(SandboxAcntName)
		params := DefaultParams(token)
		params.Endpoint = SandboxEndpoint
		params.Account = acnt
		return NewClient(params)
	}
	token := config.MustGetFromEnv(ApiKeyName)
	acnt := config.MustGetFromEnv(AcntName)
	params := DefaultParams(token)
	params.Account = acnt
	return NewClient(params)
}

// GetOptionExpirationDates returns list of an option's expiration dates.
// todo: can be enhanced with strikes and exp.type
func (tc *Client) GetOptionExpirationDates(symbol string) ([]time.Time, error) {
	params := "?symbol=" + symbol
	url := tc.endpoint + "/v1/markets/options/expirations" + params
	var result struct {
		Expirations struct {
			Date []DateTime
		}
	}
	err := tc.getJSON(url, &result)

	times := make([]time.Time, len(result.Expirations.Date))
	for i, dt := range result.Expirations.Date {
		times[i] = dt.Time
	}

	return times, err
}

func (tc *Client) GetOptionExpirations(symbol string) (*ExpirationInfo, error) {
	params := "?symbol=" + symbol + "&strikes=true&expirationType=true"
	url := tc.endpoint + "/v1/markets/options/expirations" + params
	var result ExpirationInfo
	err := tc.getJSON(url, &result)
	return &result, err
}
