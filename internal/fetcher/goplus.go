package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// GoPlusFetcher checks the address against GoPlus Security APIs.
type GoPlusFetcher struct {
	client *resty.Client
}

// NewGoPlusFetcher configures resty for the free security endpoint.
func NewGoPlusFetcher() *GoPlusFetcher {
	return &GoPlusFetcher{
		client: resty.New().SetTimeout(10 * time.Second).SetRetryCount(2),
	}
}

// GetSecurityRisks queries whether the address is involved in blacklists, phishing, or crypto mixers.
func (f *GoPlusFetcher) GetSecurityRisks(ctx context.Context, address string) ([]string, error) {
	// The GoPlus result for address_security endpoint returns the security struct directly within result.
	type GoPlusResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Result  struct {
			PhishingActivities     string `json:"phishing_activities"`
			BlacklistDoubt         string `json:"blacklist_doubt"`
			Mixer                  string `json:"mixer"`
			HoneypotRelatedAddress string `json:"honeypot_related_address"`
			Cybercrime             string `json:"cybercrime"`
		} `json:"result"`
	}

	var resp GoPlusResponse
	_, err := f.client.R().
		SetContext(ctx).
		SetResult(&resp).
		Get(fmt.Sprintf("https://api.gopluslabs.io/api/v1/address_security/%s?chain_id=1", address))

	if err != nil {
		return nil, err
	}

	if resp.Code != 1 {
		return nil, fmt.Errorf("goplus API error: %s", resp.Message)
	}

	var risks []string
	if resp.Result.PhishingActivities == "1" {
		risks = append(risks, "Phishing")
	}
	if resp.Result.BlacklistDoubt == "1" {
		risks = append(risks, "Blacklist Doubt")
	}
	if resp.Result.Mixer == "1" {
		risks = append(risks, "Mixer")
	}
	if resp.Result.Cybercrime == "1" {
		risks = append(risks, "Cybercrime")
	}
	if resp.Result.HoneypotRelatedAddress == "1" {
		risks = append(risks, "Honeypot")
	}

	return risks, nil
}
