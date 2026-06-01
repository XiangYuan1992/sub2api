import { describe, it, expect } from 'vitest'
import { officialCny, groupCny, discountLabel, savingPercent } from '../pricingCalc'

describe('pricingCalc', () => {
  it('officialCny converts per-token USD to ¥/1M tokens', () => {
    // 0.0000025 * 1e6 * 7 = 17.5
    expect(officialCny(0.0000025, 7)).toBeCloseTo(17.5, 5)
  })

  it('officialCny returns null for null price', () => {
    expect(officialCny(null, 7)).toBeNull()
  })

  it('groupCny applies rate multiplier', () => {
    // official 17.5 * 0.5 = 8.75
    expect(groupCny(0.0000025, 7, 0.5)).toBeCloseTo(8.75, 5)
  })

  it('groupCny returns null for null price', () => {
    expect(groupCny(null, 7, 0.5)).toBeNull()
  })

  it('discountLabel = multiplier * 10 折', () => {
    expect(discountLabel(0.5)).toBe('5折')
    expect(discountLabel(0.7)).toBe('7折')
    expect(discountLabel(0.35)).toBe('3.5折')
  })

  it('savingPercent rounds (1-mult)*100', () => {
    expect(savingPercent(0.5)).toBe(50)
    expect(savingPercent(0.07)).toBe(93)
  })

  it('savingPercent returns 0 when multiplier >= 1', () => {
    expect(savingPercent(1)).toBe(0)
    expect(savingPercent(1.2)).toBe(0)
  })
})
