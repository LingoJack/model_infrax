import type { Snapshot } from './types'

export async function fetchConfig(): Promise<Snapshot> {
  const r = await fetch('/api/config')
  if (!r.ok) throw new Error(`加载配置失败: HTTP ${r.status}`)
  return r.json()
}

export async function saveConfig(values: Record<string, unknown>): Promise<void> {
  const r = await fetch('/api/config', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ values }),
  })
  const body = await r.json()
  if (body.error) throw new Error(body.error)
}

export async function generate(): Promise<void> {
  const r = await fetch('/api/generate', { method: 'POST' })
  const body = await r.json()
  if (body.error) throw new Error(body.error)
}
