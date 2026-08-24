interface Props {
  busy: boolean
  onSave: () => void
  onSaveAndGenerate: () => void
}

// 底部操作栏
export function Toolbar({ busy, onSave, onSaveAndGenerate }: Props) {
  return (
    <div className="toolbar">
      <button type="button" disabled={busy} onClick={onSave}>
        保存
      </button>
      <button type="button" disabled={busy} onClick={onSaveAndGenerate}>
        保存并生成
      </button>
    </div>
  )
}
