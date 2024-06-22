package vocode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type PlanType string

const (
	PlanFree       PlanType = "plan_free"
	PlanDeveloper  PlanType = "plan_developer"
	PlanEnterprise PlanType = "plan_enterprise"
	PlanUnlimited  PlanType = "plan_unlimited"
)

type Usage struct {
	UserID              string   `json:"user_id"`
	PlanType            PlanType `json:"plan_type"`
	MonthlyMinutes      int      `json:"monthly_usage_minutes"`
	MonthlyLimitMinutes int      `json:"monthly_usage_limit_minutes"`
}

func (c *Client) GetUsage(ctx context.Context) (*Usage, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/usage")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	usage := new(Usage)
	if err := json.NewDecoder(resp.Body).Decode(usage); err != nil {
		return nil, err
	}

	return usage, nil
}
