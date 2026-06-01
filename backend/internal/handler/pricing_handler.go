package handler

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// PricingChannelLister abstracts ListAvailable for testing.
type PricingChannelLister interface {
	ListAvailable(ctx context.Context) ([]service.AvailableChannel, error)
}

// PricingRuntimeProvider abstracts pricing runtime config reading.
type PricingRuntimeProvider interface {
	GetPricingPageRuntime(ctx context.Context) service.PricingPageRuntime
}

// PricingHandler handles the public model pricing page endpoint.
type PricingHandler struct {
	channels PricingChannelLister
	settings PricingRuntimeProvider
}

func NewPricingHandler(channels PricingChannelLister, settings PricingRuntimeProvider) *PricingHandler {
	return &PricingHandler{channels: channels, settings: settings}
}

type pricingModel struct {
	Name           string   `json:"name"`
	InputPrice     *float64 `json:"input_price"`
	OutputPrice    *float64 `json:"output_price"`
	CacheReadPrice *float64 `json:"cache_read_price"`
}

type pricingGroup struct {
	ID               int64          `json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	RateMultiplier   float64        `json:"rate_multiplier"`
	SubscriptionType string         `json:"subscription_type"`
	Models           []pricingModel `json:"models"`
}

type pricingPlatform struct {
	Platform string         `json:"platform"`
	Groups   []pricingGroup `json:"groups"`
}

type pricingResponse struct {
	Enabled            bool              `json:"enabled"`
	CNYRate            float64           `json:"cny_rate"`
	RechargeMultiplier float64           `json:"recharge_multiplier"`
	Platforms          []pricingPlatform `json:"platforms"`
}

// Get handles GET /api/v1/pricing.
func (h *PricingHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	rt := h.settings.GetPricingPageRuntime(ctx)
	if !rt.Enabled {
		response.Success(c, pricingResponse{Enabled: false, CNYRate: rt.CNYRate, RechargeMultiplier: rt.RechargeMultiplier})
		return
	}

	chans, err := h.channels.ListAvailable(ctx)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	allowed := make(map[int64]struct{}, len(rt.GroupIDs))
	for _, id := range rt.GroupIDs {
		allowed[id] = struct{}{}
	}
	groupAllowed := func(id int64) bool {
		if len(allowed) == 0 {
			return true
		}
		_, ok := allowed[id]
		return ok
	}

	type groupAgg struct {
		group *pricingGroup
		seen  map[string]struct{}
	}
	platformOrder := []string{}
	platforms := map[string]map[int64]*groupAgg{}

	for ci := range chans {
		ch := &chans[ci]
		if ch.Status != service.StatusActive {
			continue
		}
		modelsByPlatform := map[string][]service.SupportedModel{}
		for _, m := range ch.SupportedModels {
			modelsByPlatform[m.Platform] = append(modelsByPlatform[m.Platform], m)
		}
		for _, g := range ch.Groups {
			if g.IsExclusive || !groupAllowed(g.ID) || g.Platform == "" {
				continue
			}
			pm, ok := platforms[g.Platform]
			if !ok {
				pm = map[int64]*groupAgg{}
				platforms[g.Platform] = pm
				platformOrder = append(platformOrder, g.Platform)
			}
			agg, ok := pm[g.ID]
			if !ok {
				agg = &groupAgg{
					group: &pricingGroup{
						ID:               g.ID,
						Name:             g.Name,
						Description:      g.Description,
						RateMultiplier:   g.RateMultiplier,
						SubscriptionType: g.SubscriptionType,
						Models:           []pricingModel{},
					},
					seen: map[string]struct{}{},
				}
				pm[g.ID] = agg
			}
			for _, m := range modelsByPlatform[g.Platform] {
				if _, dup := agg.seen[m.Name]; dup {
					continue
				}
				agg.seen[m.Name] = struct{}{}
				pmModel := pricingModel{Name: m.Name}
				if m.Pricing != nil {
					pmModel.InputPrice = m.Pricing.InputPrice
					pmModel.OutputPrice = m.Pricing.OutputPrice
					pmModel.CacheReadPrice = m.Pricing.CacheReadPrice
				}
				agg.group.Models = append(agg.group.Models, pmModel)
			}
		}
	}

	out := pricingResponse{Enabled: true, CNYRate: rt.CNYRate, RechargeMultiplier: rt.RechargeMultiplier, Platforms: []pricingPlatform{}}
	for _, p := range platformOrder {
		pg := pricingPlatform{Platform: p, Groups: []pricingGroup{}}
		for _, agg := range platforms[p] {
			pg.Groups = append(pg.Groups, *agg.group)
		}
		out.Platforms = append(out.Platforms, pg)
	}
	response.Success(c, out)
}
