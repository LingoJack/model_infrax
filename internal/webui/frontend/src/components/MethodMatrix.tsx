import { useMemo } from 'react'
import type { MethodMeta } from '../types'

interface Props {
  methods: MethodMeta[]
  selected: Set<string>
  onToggle: (id: string, checked: boolean) => void
  onToggleGroup: (group: string) => void
  onSelectAll: () => void
  onClearAll: () => void
}

// DAO 方法生成矩阵：按分组勾选；全部勾选 = 全部生成（提交时删除配置键）
export function MethodMatrix({ methods, selected, onToggle, onToggleGroup, onSelectAll, onClearAll }: Props) {
  const groups = useMemo(() => {
    const m = new Map<string, MethodMeta[]>()
    for (const item of methods) {
      const list = m.get(item.group) ?? []
      list.push(item)
      m.set(item.group, list)
    }
    return [...m.entries()]
  }, [methods])

  return (
    <div className="card">
      <h2>
        DAO 方法生成矩阵
        <span className="sub">
          （全部勾选 = 全部生成）
          <button type="button" className="link" onClick={onSelectAll}>
            全选
          </button>
          <button type="button" className="link" onClick={onClearAll}>
            全不选
          </button>
        </span>
      </h2>
      <div className="methods">
        {groups.map(([group, items]) => (
          <div className="mgroup" key={group}>
            <b>
              {group}
              <span className="ops" onClick={() => onToggleGroup(group)}>
                组内切换
              </span>
            </b>
            <div className="mlist">
              {items.map((m) => (
                <label className="m" key={m.id}>
                  <input
                    type="checkbox"
                    checked={selected.has(m.id)}
                    onChange={(e) => onToggle(m.id, e.target.checked)}
                  />
                  <span className="info">
                    <span className="name">{m.id}</span>
                    <span className="d">{m.desc}</span>
                  </span>
                </label>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
