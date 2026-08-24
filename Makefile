SHELL := /bin/bash

.PHONY: \
	build \
	ui \
	install \
	push \
	checkpoint \
	rollback \
	current_dir \
	tag

current_dir:
	@echo -----------------------------
	@echo "current dir:"
	@echo "    $$(pwd)"
	@echo -----------------------------
	@echo

push: build
	@scripts/push.sh

build:
	@scripts/build.sh

# 构建Web UI 前端（产物输出到 internal/webui/dist/ 并随 Go embed 打包）
ui:
	@cd internal/webui/frontend && npm install && npm run build

install:
	@echo "🔨 正在安装 jen 到 GOPATH/bin..."
	@go install ./cmd/jen
	@echo "✅ 安装成功！可以直接使用 jen 命令"

checkpoint: current_dir
	@if git diff --quiet && git diff --cached --quiet; then \
		echo "没有需要提交的变更"; \
	else \
		git add . && \
		git commit -m "checkpoint: 临时保存点 [$$(date +%Y%m%d-%H%M%S)]" && \
		echo "=== 已提交 ===" && \
		git --no-pager log -n 1; \
	fi

rollback: current_dir
	@echo "=== 警告：这将丢弃所有未提交的修改 ==="
	@echo "当前未提交的修改："
	@git status --short
	@echo ""
	@read -p "确定要丢弃所有修改吗？: " confirm && \
	if [ -z "$$confirm" ]; then \
		confirm="y"; \
	fi && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		git checkout -- . && \
		echo "=== 回滚完成，工作区已清空 ===" && \
		git status; \
	else \
		echo "操作已取消"; \
	fi

# 打 tag 并推送
# 用法:
#   make tag           -> 自动递增 patch 版本 (v0.1.0 -> v0.1.1)
#   make tag V=v1.0.0  -> 指定版本号
tag:
	@if [ -n "$(V)" ]; then \
		version="$(V)"; \
	else \
		latest=$$(git tag --sort=-v:refname | head -n 1); \
		if [ -z "$$latest" ]; then \
			version="v0.1.0"; \
		else \
			major=$$(echo $$latest | sed 's/^v//' | cut -d. -f1); \
			minor=$$(echo $$latest | sed 's/^v//' | cut -d. -f2); \
			patch=$$(echo $$latest | sed 's/^v//' | cut -d. -f3); \
			patch=$$((patch + 1)); \
			version="v$$major.$$minor.$$patch"; \
		fi; \
	fi; \
	echo "📦 当前最新 tag: $$(git tag --sort=-v:refname | head -n 1 || echo '无')"; \
	echo "🏷️  即将创建 tag: $$version"; \
	read -p "确认打 tag $$version 吗？[y/N] " confirm; \
	if [ -z "$$confirm" ]; then confirm="n"; fi; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		git tag -a $$version -m "release: $$version" && \
		git push origin $$version && \
		echo "✅ tag $$version 已创建并推送到远程"; \
	else \
		echo "操作已取消"; \
	fi