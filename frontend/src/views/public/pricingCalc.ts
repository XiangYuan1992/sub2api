const PER_M = 1_000_000

/** per-token USD → ¥/1M tokens；price 为 null 返回 null。 */
export function officialCny(perTokenUsd: number | null, cnyRate: number): number | null {
  if (perTokenUsd == null) return null
  return perTokenUsd * PER_M * cnyRate
}

/** 分组价 = 官方价 × 倍率。 */
export function groupCny(
  perTokenUsd: number | null,
  cnyRate: number,
  multiplier: number
): number | null {
  const official = officialCny(perTokenUsd, cnyRate)
  return official == null ? null : official * multiplier
}

/** 折扣标签：倍率 × 10 折（去掉末尾 .0）。 */
export function discountLabel(multiplier: number): string {
  const v = Math.round(multiplier * 10 * 10) / 10
  return `${v}折`
}

/** 节省百分比：round((1-倍率)*100)；倍率≥1 返回 0。 */
export function savingPercent(multiplier: number): number {
  if (multiplier >= 1) return 0
  return Math.round((1 - multiplier) * 100)
}
