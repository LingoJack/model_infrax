import { useEffect, useMemo, useRef, useState } from 'react'
import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import { fetchFiles, fetchFileContent } from '../api'

hljs.registerLanguage('go', go)

interface Props {
  reloadKey: number // 生成完成后自增，触发重新加载
  onError: (msg: string) => void
}

// 产物文件浏览：左侧文件树，右侧代码视图（语法高亮 + 复制）
export function OutputView({ reloadKey, onError }: Props) {
  const [outputPath, setOutputPath] = useState('')
  const [files, setFiles] = useState<string[]>([])
  const [active, setActive] = useState<string>('')
  const [content, setContent] = useState('')
  const [copied, setCopied] = useState(false)
  const codeRef = useRef<HTMLElement>(null)

  useEffect(() => {
    fetchFiles()
      .then((r) => {
        setOutputPath(r.output_path)
        setFiles(r.files)
        if (r.files.length > 0 && !r.files.includes(active)) {
          setActive(r.files[0])
        }
      })
      .catch((e) => onError(String(e.message ?? e)))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [reloadKey])

  useEffect(() => {
    if (!active) return
    fetchFileContent(active)
      .then(setContent)
      .catch((e) => {
        setContent('')
        onError(String(e.message ?? e))
      })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [active])

  useEffect(() => {
    if (codeRef.current) {
      codeRef.current.textContent = content
      hljs.highlightElement(codeRef.current)
    }
  }, [content])

  // 按目录分组
  const groups = useMemo(() => {
    const m = new Map<string, string[]>()
    for (const f of files) {
      const dir = f.includes('/') ? f.slice(0, f.lastIndexOf('/')) : '.'
      const list = m.get(dir) ?? []
      list.push(f)
      m.set(dir, list)
    }
    return [...m.entries()]
  }, [files])

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(content)
    } catch {
      // 兼容非 HTTPS 场景
      const ta = document.createElement('textarea')
      ta.value = content
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  if (files.length === 0) {
    return (
      <div className="card">
        <h2>生成产物</h2>
        <p className="empty">
          {outputPath
            ? `输出目录 ${outputPath} 下暂无文件，点击「保存并生成」后在这里查看`
            : '尚未生成代码，点击「保存并生成」后在这里查看'}
        </p>
      </div>
    )
  }

  return (
    <div className="card">
      <h2>
        生成产物
        <span className="sub">
          {outputPath}（{files.length} 个文件）
        </span>
      </h2>
      <div className="outview">
        <div className="ftree">
          {groups.map(([dir, items]) => (
            <div key={dir} className="fgroup">
              <div className="fdir">{dir}/</div>
              {items.map((f) => (
                <div
                  key={f}
                  className={`fitem${f === active ? ' active' : ''}`}
                  onClick={() => setActive(f)}
                >
                  {f.slice(f.lastIndexOf('/') + 1)}
                </div>
              ))}
            </div>
          ))}
        </div>
        <div className="fcode">
          <div className="fhead">
            <span className="fname">{active}</span>
            <button type="button" className="copybtn" onClick={copy}>
              {copied ? '✓ 已复制' : '复制'}
            </button>
          </div>
          <pre>
            <code ref={codeRef} className="language-go" />
          </pre>
        </div>
      </div>
    </div>
  )
}
