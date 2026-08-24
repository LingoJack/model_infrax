import type { ConfigField } from '../types'
import { FieldInput } from './FieldInput'

interface Props {
  title: string
  fields: ConfigField[]
  values: Record<string, unknown>
  onChange: (key: string, value: unknown) => void
}

// 分组卡片
export function FieldGroup({ title, fields, values, onChange }: Props) {
  return (
    <div className="card">
      <h2>{title}</h2>
      {fields.map((f) => (
        <FieldInput key={f.key} field={f} value={values[f.key]} onChange={onChange} />
      ))}
    </div>
  )
}
