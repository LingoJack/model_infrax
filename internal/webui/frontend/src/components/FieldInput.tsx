import type { ConfigField } from '../types'
import { Select } from './Select'

interface Props {
  field: ConfigField
  value: unknown
  onChange: (key: string, value: unknown) => void
}

// 按字段类型渲染表单控件
export function FieldInput({ field, value, onChange }: Props) {
  const { key, label, type, opts, hint } = field

  let ctl: JSX.Element
  switch (type) {
    case 'select':
      ctl = (
        <Select
          value={String(value || (opts ?? [])[0] || '')}
          options={opts ?? []}
          onChange={(v) => onChange(key, v)}
        />
      )
      break
    case 'bool':
      ctl = (
        <input
          type="checkbox"
          checked={Boolean(value)}
          onChange={(e) => onChange(key, e.target.checked)}
        />
      )
      break
    case 'list':
      ctl = (
        <input
          type="text"
          value={Array.isArray(value) ? value.join(', ') : String(value ?? '')}
          placeholder="逗号分隔"
          onChange={(e) =>
            onChange(
              key,
              e.target.value.split(/[,，]/).map((s) => s.trim()).filter(Boolean),
            )
          }
        />
      )
      break
    case 'number':
      ctl = (
        <input
          type="number"
          value={Number(value ?? 0)}
          onChange={(e) => onChange(key, Number(e.target.value || 0))}
        />
      )
      break
    default:
      ctl = (
        <input
          type={type}
          value={String(value ?? '')}
          onChange={(e) => onChange(key, e.target.value)}
        />
      )
  }

  return (
    <div className="row">
      <label>{label}</label>
      {ctl}
      {hint && <span className="hint">{hint}</span>}
    </div>
  )
}
