// 与后端 internal/webui/server.go 的 configField/MethodMeta JSON 对齐

export type FieldType = 'text' | 'password' | 'number' | 'bool' | 'select' | 'list'

export interface ConfigField {
  key: string
  label: string
  type: FieldType
  opts?: string[]
  group: string
  hint?: string
}

export interface MethodMeta {
  id: string
  desc: string
  group: string
}

export interface Snapshot {
  config_path: string
  values: Record<string, unknown>
  fields: ConfigField[]
  dao_methods: MethodMeta[]
}
