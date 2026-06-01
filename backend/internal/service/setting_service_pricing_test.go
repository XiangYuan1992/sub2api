package service

import (
	"context"
	"testing"
)

func TestGetPricingPageRuntime_ParsesValues(t *testing.T) {
	repo := &pricingRuntimeRepoStub{vals: map[string]string{
		SettingKeyPricingPageEnabled: "true",
		SettingKeyPricingCNYRate:     "7.3",
		SettingKeyPricingGroupIDs:    "[1,2,3]",
	}}
	svc := &SettingService{settingRepo: repo}
	rt := svc.GetPricingPageRuntime(context.Background())
	if !rt.Enabled {
		t.Fatal("expected enabled")
	}
	if rt.CNYRate != 7.3 {
		t.Fatalf("rate = %v, want 7.3", rt.CNYRate)
	}
	if len(rt.GroupIDs) != 3 || rt.GroupIDs[0] != 1 {
		t.Fatalf("group ids = %v", rt.GroupIDs)
	}
}

func TestGetPricingPageRuntime_DefaultsRateWhenInvalid(t *testing.T) {
	repo := &pricingRuntimeRepoStub{vals: map[string]string{
		SettingKeyPricingPageEnabled: "false",
		SettingKeyPricingCNYRate:     "",
	}}
	svc := &SettingService{settingRepo: repo}
	rt := svc.GetPricingPageRuntime(context.Background())
	if rt.Enabled {
		t.Fatal("expected disabled")
	}
	if rt.CNYRate != 7.0 {
		t.Fatalf("rate = %v, want 7.0 default", rt.CNYRate)
	}
	if len(rt.GroupIDs) != 0 {
		t.Fatalf("expected empty group ids, got %v", rt.GroupIDs)
	}
}

// pricingRuntimeRepoStub implements SettingRepository for pricing runtime tests.
type pricingRuntimeRepoStub struct {
	vals map[string]string
}

func (s *pricingRuntimeRepoStub) Get(_ context.Context, _ string) (*Setting, error) {
	return nil, nil
}

func (s *pricingRuntimeRepoStub) GetValue(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (s *pricingRuntimeRepoStub) Set(_ context.Context, _, _ string) error {
	return nil
}

func (s *pricingRuntimeRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := s.vals[k]; ok {
			out[k] = v
		}
	}
	return out, nil
}

func (s *pricingRuntimeRepoStub) SetMultiple(_ context.Context, _ map[string]string) error {
	return nil
}

func (s *pricingRuntimeRepoStub) GetAll(_ context.Context) (map[string]string, error) {
	return nil, nil
}

func (s *pricingRuntimeRepoStub) Delete(_ context.Context, _ string) error {
	return nil
}
