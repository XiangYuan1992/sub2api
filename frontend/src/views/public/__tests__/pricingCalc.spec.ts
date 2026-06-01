import { describe, it, expect } from 'vitest'
import {
  officialCny,
  officialUsd,
  siteCnyRate,
  groupCny,
  discountLabel,
  discountLabelFromSaving,
  savingPercent,
  modelVersion,
  sortModelsByVersionDesc
} from '../pricingCalc'

describe('pricingCalc', () => {
  it('officialCny converts per-token USD to ¥/1M tokens', () => {
    // 0.0000025 * 1e6 * 7 = 17.5
    expect(officialCny(0.0000025, 7)).toBeCloseTo(17.5, 5)
  })

  it('officialCny returns null for null price', () => {
    expect(officialCny(null, 7)).toBeNull()
  })

  it('officialUsd converts per-token USD to $/1M tokens', () => {
    // 0.0000025 * 1e6 = 2.5
    expect(officialUsd(0.0000025)).toBeCloseTo(2.5, 5)
  })

  it('officialUsd returns null for null price', () => {
    expect(officialUsd(null)).toBeNull()
  })

  it('siteCnyRate = 1 / rechargeMultiplier', () => {
    // multiplier 0.2 → ¥5 per $1
    expect(siteCnyRate(0.2)).toBeCloseTo(5, 5)
    expect(siteCnyRate(1)).toBeCloseTo(1, 5)
  })

  it('siteCnyRate returns 0 for invalid multiplier', () => {
    expect(siteCnyRate(0)).toBe(0)
    expect(siteCnyRate(-1)).toBe(0)
  })

  it('groupCny = official USD × rateMultiplier × siteRate', () => {
    // usd 2.5 * rate 0.5 * siteRate(0.2)=5 = 6.25
    expect(groupCny(0.0000025, 0.2, 0.5)).toBeCloseTo(6.25, 5)
  })

  it('groupCny returns null for null price', () => {
    expect(groupCny(null, 0.2, 0.5)).toBeNull()
  })

  it('discountLabel = price ratio * 10 折', () => {
    expect(discountLabel(0.5)).toBe('5折')
    expect(discountLabel(0.7)).toBe('7折')
    expect(discountLabel(0.35)).toBe('3.5折')
    expect(discountLabel(0.051)).toBe('0.5折')
  })

  it('discountLabelFromSaving derives 折 from saving percent', () => {
    expect(discountLabelFromSaving(94.9)).toBe('0.5折')
    expect(discountLabelFromSaving(64.3)).toBe('3.6折')
    expect(discountLabelFromSaving(0)).toBe('')
  })

  it('savingPercent = (官方¥ − 实际¥) / 官方¥', () => {
    // siteRate(0.2)=5, actual rate = 0.5*5 = 2.5, official = 7 → (1 - 2.5/7)*100 = 64.3
    expect(savingPercent(0.2, 0.5, 7)).toBeCloseTo(64.3, 1)
    // siteRate(1)=1, actual = 0.5, official = 7 → (1 - 0.5/7)*100 = 92.9
    expect(savingPercent(1, 0.5, 7)).toBeCloseTo(92.9, 1)
  })

  it('savingPercent returns 0 when actual >= official', () => {
    // siteRate(1)=1, actual = 7*1 = 7 = official → 0
    expect(savingPercent(1, 7, 7)).toBe(0)
  })

  it('savingPercent returns 0 when cnyRate invalid', () => {
    expect(savingPercent(0.2, 0.5, 0)).toBe(0)
  })

  it('modelVersion extracts first numeric version', () => {
    expect(modelVersion('gpt-5.5')).toBe(5.5)
    expect(modelVersion('gpt-5')).toBe(5)
    expect(modelVersion('gpt-4.1')).toBe(4.1)
    expect(modelVersion('claude-opus')).toBeNull()
  })

  it('sortModelsByVersionDesc sorts descending, no-version last', () => {
    const out = sortModelsByVersionDesc([
      { name: 'gpt-4.1' },
      { name: 'claude-opus' },
      { name: 'gpt-5.5' },
      { name: 'gpt-5' }
    ])
    expect(out.map((m) => m.name)).toEqual(['gpt-5.5', 'gpt-5', 'gpt-4.1', 'claude-opus'])
  })
})
