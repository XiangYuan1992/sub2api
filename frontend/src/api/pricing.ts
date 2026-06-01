import { apiClient } from './client'

export interface PricingModel {
  name: string
  input_price: number | null
  output_price: number | null
  cache_read_price: number | null
}

export interface PricingGroup {
  id: number
  name: string
  description: string
  rate_multiplier: number
  subscription_type: string
  models: PricingModel[]
}

export interface PricingPlatform {
  platform: string
  groups: PricingGroup[]
}

export interface PricingResponse {
  enabled: boolean
  cny_rate: number
  platforms: PricingPlatform[]
}

export async function getPricing(): Promise<PricingResponse> {
  const { data } = await apiClient.get<PricingResponse>('/pricing')
  return data
}
