import { useEffect, useRef, useState } from 'react'

interface Props {
  value: string
  options: string[]
  onChange: (v: string) => void
}

// 自定义下拉选择（替代原生 select 的系统下拉列表，保持整体风格一致）
export function Select({ value, options, onChange }: Props) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  // 点击组件外部时收起
  useEffect(() => {
    if (!open) return
    const close = (e: MouseEvent) => {
      if (!ref.current?.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', close)
    return () => document.removeEventListener('mousedown', close)
  }, [open])

  return (
    <div className="sel" ref={ref}>
      <button
        type="button"
        className={`sel-btn${open ? ' open' : ''}`}
        onClick={() => setOpen((o) => !o)}
      >
        <span className="sel-value">{value}</span>
        <svg
          className="sel-arrow"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="m6 9 6 6 6-6" />
        </svg>
      </button>
      {open && (
        <ul className="sel-list">
          {options.map((o) => (
            <li
              key={o}
              className={o === value ? 'active' : ''}
              onMouseDown={(e) => {
                e.preventDefault()
                onChange(o)
                setOpen(false)
              }}
            >
              <span>{o}</span>
              {o === value && (
                <svg
                  width="13"
                  height="13"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M20 6 9 17l-5-5" />
                </svg>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
