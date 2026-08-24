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

export interface FileList {
  output_path: string
  files: string[]
}

export async function fetchFiles(): Promise<FileList> {
  const r = await fetch('/api/files')
  if (!r.ok) throw new Error(`加载文件列表失败: HTTP ${r.status}`)
  return r.json()
}

export async function fetchFileContent(path: string): Promise<string> {
  const r = await fetch(`/api/file?path=${encodeURIComponent(path)}`)
  const body = await r.json()
  if (body.error) throw new Error(body.error)
  return body.content
}
