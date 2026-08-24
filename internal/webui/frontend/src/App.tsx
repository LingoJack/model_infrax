import { useEffect, useMemo, useState } from 'react'
import { fetchConfig, saveConfig, generate } from './api'
import type { Snapshot } from './types'
import { FieldGroup } from './components/FieldGroup'
import { MethodMatrix } from './components/MethodMatrix'
import { OutputView } from './components/OutputView'
import { Toolbar } from './components/Toolbar'

export default function App() {
  const [snap, setSnap] = useState<Snapshot | null>(null)
  const [values, setValues] = useState<Record<string, unknown>>({})
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [busy, setBusy] = useState(false)
  const [msg, setMsg] = useState<{ text: string; ok: boolean } | null>(null)
  const [tab, setTab] = useState<'config' | 'output'>('config')
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    fetchConfig()
      .then((s) => {
        setSnap(s)
        setValues(s.values)
        const methods: string[] =
          Array.isArray(s.values['generate_option.dao_methods']) && s.values['generate_option.dao_methods'].length
            ? (s.values['generate_option.dao_methods'] as string[])
            : s.dao_methods.map((m) => m.id)
        setSelected(new Set(methods))
      })
      .catch((e) => setMsg({ text: String(e.message ?? e), ok: false }))
  }, [])

  // 按分组组织字段；数据库连接组仅在 database 模式显示
  const groups = useMemo(() => {
    if (!snap) return []
    const mode = values['generate_config.generate_mode']
    const m = new Map<string, typeof snap.fields>()
    for (const f of snap.fields) {
      if (f.group === '数据库连接' && mode !== 'database') continue
      const list = m.get(f.group) ?? []
      list.push(f)
      m.set(f.group, list)
    }
    return [...m.entries()]
  }, [snap, values])

  const setValue = (key: string, value: unknown) => {
    setValues((prev) => ({ ...prev, [key]: value }))
  }

  const toggleMethod = (id: string, checked: boolean) => {
    setSelected((prev) => {
      const next = new Set(prev)
      if (checked) next.add(id)
      else next.delete(id)
      return next
    })
  }

  const toggleGroup = (group: string) => {
    setSelected((prev) => {
      const ids = (snap?.dao_methods ?? []).filter((m) => m.group === group).map((m) => m.id)
      const allOn = ids.every((id) => prev.has(id))
      const next = new Set(prev)
      ids.forEach((id) => (allOn ? next.delete(id) : next.add(id)))
      return next
    })
  }

  const selectAll = () => setSelected(new Set((snap?.dao_methods ?? []).map((m) => m.id)))

  const collect = (): Record<string, unknown> => {
    const all = snap?.dao_methods.map((m) => m.id) ?? []
    const ms = all.filter((id) => selected.has(id))
    // 全选 = 提交空列表（后端语义：全部生成）
    return { ...values, 'generate_option.dao_methods': ms.length === all.length ? [] : ms }
  }

  const doSave = async (withGenerate: boolean) => {
    setBusy(true)
    setMsg(null)
    try {
      await saveConfig(collect())
      if (!withGenerate) {
        setMsg({ text: '已保存', ok: true })
        return
      }
      setMsg({ text: '已保存，正在生成…', ok: true })
      await generate()
      setReloadKey((k) => k + 1)
      setTab('output')
      setMsg({ text: '生成完成', ok: true })
    } catch (e) {
      setMsg({ text: '失败: ' + (e instanceof Error ? e.message : String(e)), ok: false })
    } finally {
      setBusy(false)
    }
  }

  if (!snap) {
    return (
      <div className="wrap">
        <p>{msg && !msg.ok ? msg.text : '加载中…'}</p>
      </div>
    )
  }

  return (
    <>
      {msg && <div id="msg" className={msg.ok ? 'ok' : 'err'}>{msg.text}</div>}
      <div className="wrap">
        <header>
          <h1>jen 生成配置</h1>
          <span className="path">{snap.config_path}</span>
          <nav className="tabs">
            <button
              type="button"
              className={tab === 'config' ? 'active' : ''}
              onClick={() => setTab('config')}
            >
              配置
            </button>
            <button
              type="button"
              className={tab === 'output' ? 'active' : ''}
              onClick={() => {
                setTab('output')
                setReloadKey((k) => k + 1)
              }}
            >
              生成产物
            </button>
          </nav>
        </header>
        {tab === 'config' ? (
          <>
            {groups.map(([title, fields]) => (
              <FieldGroup key={title} title={title} fields={fields} values={values} onChange={setValue} />
            ))}
            <MethodMatrix
              methods={snap.dao_methods}
              selected={selected}
              onToggle={toggleMethod}
              onToggleGroup={toggleGroup}
              onSelectAll={selectAll}
            />
          </>
        ) : (
          <OutputView reloadKey={reloadKey} onError={(t) => setMsg({ text: t, ok: false })} />
        )}
      </div>
      {tab === 'config' && (
        <Toolbar busy={busy} onSave={() => doSave(false)} onSaveAndGenerate={() => doSave(true)} />
      )}
    </>
  )
}
