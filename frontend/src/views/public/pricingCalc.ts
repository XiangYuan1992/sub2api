const PER_M = 1_000_000

/** per-token USD → ¥/1M tokens（官方价，按官方汇率）；price 为 null 返回 null。 */
export function officialCny(perTokenUsd: number | null, cnyRate: number): number | null {
  if (perTokenUsd == null) return null
  return perTokenUsd * PER_M * cnyRate
}

/** per-token USD → $/1M tokens（官方美元价）；price 为 null 返回 null。 */
export function officialUsd(perTokenUsd: number | null): number | null {
  if (perTokenUsd == null) return null
  return perTokenUsd * PER_M
}

/** 本站汇率（¥ / $1）= 1 / 充值倍率（每 1 CNY 充多少 USD）；非法时返回 0。 */
export function siteCnyRate(rechargeMultiplier: number): number {
  return rechargeMultiplier > 0 ? 1 / rechargeMultiplier : 0
}

/** 分组实际价(¥) = 官方美元价 × 分组倍率 × 本站汇率。 */
export function groupCny(
  perTokenUsd: number | null,
  rechargeMultiplier: number,
  rateMultiplier: number
): number | null {
  const usd = officialUsd(perTokenUsd)
  return usd == null ? null : usd * rateMultiplier * siteCnyRate(rechargeMultiplier)
}

/** 折扣标签：相对官方价的实际比例 × 10 折（去掉末尾 .0）。 */
export function discountLabel(priceRatio: number): string {
  const v = Math.round(priceRatio * 10 * 10) / 10
  return `${v}折`
}

/** 按节省幅度反推折扣标签：实际付费占官方价比例 × 10 折。 */
export function discountLabelFromSaving(savingPct: number): string {
  if (savingPct <= 0) return ''
  return discountLabel(1 - savingPct / 100)
}

/** 节省幅度 % = (官方¥ − 实际¥) / 官方¥ × 100，保留 1 位小数；≤0 返回 0。 */
export function savingPercent(
  rechargeMultiplier: number,
  rateMultiplier: number,
  cnyRate: number
): number {
  if (cnyRate <= 0) return 0
  const actual = rateMultiplier * siteCnyRate(rechargeMultiplier)
  const pct = (1 - actual / cnyRate) * 100
  return pct > 0 ? Math.round(pct * 10) / 10 : 0
}

/** 从模型名提取版本号（首个 \d+(.\d+)? ）；无则返回 null。 */
export function modelVersion(name: string): number | null {
  const m = name.match(/\d+(?:\.\d+)?/)
  return m ? parseFloat(m[0]) : null
}

/** 按版本号降序排序模型；无版本号的排最后，保持原相对顺序。 */
export function sortModelsByVersionDesc<T extends { name: string }>(models: T[]): T[] {
  return models
    .map((m, i) => ({ m, i, v: modelVersion(m.name) }))
    .sort((a, b) => {
      if (a.v == null && b.v == null) return a.i - b.i
      if (a.v == null) return 1
      if (b.v == null) return -1
      return b.v - a.v || a.i - b.i
    })
    .map((x) => x.m)
}
