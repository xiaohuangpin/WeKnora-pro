import { get } from '@/utils/request'

export interface SystemInfo {
  version: string
  commit_id?: string
  build_time?: string
  go_version?: string
}

export function getSystemInfo(): Promise<{ data: SystemInfo }> {
  return get('/api/v1/system/info')
}
