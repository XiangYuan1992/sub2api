//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type pricingChannelLister struct{ chans []service.AvailableChannel }

func (l *pricingChannelLister) ListAvailable(_ context.Context) ([]service.AvailableChannel, error) {
	return l.chans, nil
}

type pricingRuntimeProvider struct{ rt service.PricingPageRuntime }

func (p *pricingRuntimeProvider) GetPricingPageRuntime(_ context.Context) service.PricingPageRuntime {
	return p.rt
}

func fptr(v float64) *float64 { return &v }

func TestPricingHandler_DisabledReturnsEnabledFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPricingHandler(
		&pricingChannelLister{},
		&pricingRuntimeProvider{rt: service.PricingPageRuntime{Enabled: false, CNYRate: 7}},
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/pricing", nil)
	h.Get(c)

	if w.Code != http.StatusOK {
		t.Fatalf("code = %d", w.Code)
	}
	var resp struct {
		Data struct {
			Enabled bool `json:"enabled"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Enabled {
		t.Fatal("expected enabled=false")
	}
}

func TestPricingHandler_AggregatesPublicGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	chans := []service.AvailableChannel{{
		Name:   "ch1",
		Status: service.StatusActive,
		Groups: []service.AvailableGroupRef{
			{ID: 1, Name: "Codex Pro", Platform: "openai", RateMultiplier: 0.5, IsExclusive: false},
			{ID: 2, Name: "Secret", Platform: "openai", RateMultiplier: 0.3, IsExclusive: true},
		},
		SupportedModels: []service.SupportedModel{
			{Name: "gpt-5.5", Platform: "openai", Pricing: &service.ChannelModelPricing{
				InputPrice: fptr(0.0000025), OutputPrice: fptr(0.000015), CacheReadPrice: fptr(0.00000025),
			}},
		},
	}}
	h := NewPricingHandler(
		&pricingChannelLister{chans: chans},
		&pricingRuntimeProvider{rt: service.PricingPageRuntime{Enabled: true, CNYRate: 7}},
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/pricing", nil)
	h.Get(c)

	var resp struct {
		Data struct {
			Enabled   bool    `json:"enabled"`
			CNYRate   float64 `json:"cny_rate"`
			Platforms []struct {
				Platform string `json:"platform"`
				Groups   []struct {
					ID             int64   `json:"id"`
					Name           string  `json:"name"`
					RateMultiplier float64 `json:"rate_multiplier"`
					Models         []struct {
						Name       string   `json:"name"`
						InputPrice *float64 `json:"input_price"`
					} `json:"models"`
				} `json:"groups"`
			} `json:"platforms"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, w.Body.String())
	}
	if !resp.Data.Enabled || resp.Data.CNYRate != 7 {
		t.Fatalf("enabled/rate wrong: %+v", resp.Data)
	}
	if len(resp.Data.Platforms) != 1 || resp.Data.Platforms[0].Platform != "openai" {
		t.Fatalf("platforms = %+v", resp.Data.Platforms)
	}
	groups := resp.Data.Platforms[0].Groups
	if len(groups) != 1 || groups[0].ID != 1 {
		t.Fatalf("expected only public group 1, got %+v", groups)
	}
	if len(groups[0].Models) != 1 || groups[0].Models[0].Name != "gpt-5.5" {
		t.Fatalf("models = %+v", groups[0].Models)
	}
}
