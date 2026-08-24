package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LingoJack/model_infrax/assets"
	"github.com/LingoJack/model_infrax/internal/logger"
)

const (
	// .model_infrax 目录及其下的文件
	modelInfraxDir = ".model_infrax"
	configFileName = "config.yml"
	schemaFileName = "schema.sql"
)

// confirmOverwrite 询问用户是否覆盖已有文件，返回 true 表示覆盖
func confirmOverwrite(reader *bufio.Reader, filePath string) bool {
	logger.ColorPrintf(logger.ColorHiYellow, "⚠ 文件已存在: %s，是否覆盖？[y/N] ", filePath)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// generateConfigTemplate 在当前目录初始化 .model_infrax 配置目录
// 生成 config.yml 和 schema.sql 模板文件
func generateConfigTemplate() error {
	logger.Infof("[generateConfigTemplate] 开始初始化 .model_infrax 配置目录...")

	stdinReader := bufio.NewReader(os.Stdin)
	configFilePath := filepath.Join(modelInfraxDir, configFileName)
	schemaFilePath := filepath.Join(modelInfraxDir, schemaFileName)

	// 创建 .model_infrax 目录
	if err := os.MkdirAll(modelInfraxDir, 0755); err != nil {
		logger.Errorf("[generateConfigTemplate] 创建目录失败: %v, 目录路径: %s", err, modelInfraxDir)
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入 config.yml（已存在则询问是否覆盖）
	writeConfig := true
	if _, err := os.Stat(configFilePath); err == nil {
		writeConfig = confirmOverwrite(stdinReader, configFilePath)
	}
	if writeConfig {
		if err := os.WriteFile(configFilePath, []byte(assets.DefaultConfigYml), 0644); err != nil {
			logger.Errorf("[generateConfigTemplate] 写入配置文件失败: %v, 文件路径: %s", err, configFilePath)
			return fmt.Errorf("写入配置文件失败: %w", err)
		}
		logger.ColorPrintf(logger.ColorHiGreen, "  ✓ 已写入 %s\n", configFilePath)
	} else {
		logger.ColorPrintf(logger.ColorWhite, "  ⏭ 跳过 %s\n", configFilePath)
	}

	// 写入 schema.sql（已存在则询问是否覆盖）
	writeSchema := true
	if _, err := os.Stat(schemaFilePath); err == nil {
		writeSchema = confirmOverwrite(stdinReader, schemaFilePath)
	}
	if writeSchema {
		if err := os.WriteFile(schemaFilePath, []byte(assets.DefaultSchemaSql), 0644); err != nil {
			logger.Errorf("[generateConfigTemplate] 写入 schema 文件失败: %v, 文件路径: %s", err, schemaFilePath)
			return fmt.Errorf("写入 schema 文件失败: %w", err)
		}
		logger.ColorPrintf(logger.ColorHiGreen, "  ✓ 已写入 %s\n", schemaFilePath)
	} else {
		logger.ColorPrintf(logger.ColorWhite, "  ⏭ 跳过 %s\n", schemaFilePath)
	}

	// 处理 .gitignore
	gitignorePath := ".gitignore"
	gitignoreEntries := []string{".model_infrax/", "/target/jen"}
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// .gitignore 不存在，创建并写入
		content := strings.Join(gitignoreEntries, "\n") + "\n"
		if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
			logger.Errorf("[generateConfigTemplate] 创建 .gitignore 失败: %v", err)
			return fmt.Errorf("创建 .gitignore 失败: %w", err)
		}
		logger.ColorPrintf(logger.ColorHiGreen, "  ✓ 已创建 %s\n", gitignorePath)
	} else {
		// .gitignore 已存在，检查并追加缺失的条目
		existingBytes, err := os.ReadFile(gitignorePath)
		if err != nil {
			logger.Errorf("[generateConfigTemplate] 读取 .gitignore 失败: %v", err)
			return fmt.Errorf("读取 .gitignore 失败: %w", err)
		}
		existingContent := string(existingBytes)
		existingLines := strings.Split(existingContent, "\n")

		var toAppend []string
		for _, entry := range gitignoreEntries {
			found := false
			for _, line := range existingLines {
				if strings.TrimSpace(line) == entry {
					found = true
					break
				}
			}
			if !found {
				toAppend = append(toAppend, entry)
			}
		}

		if len(toAppend) > 0 {
			appendContent := ""
			// 确保已有内容末尾有换行
			if len(existingContent) > 0 && !strings.HasSuffix(existingContent, "\n") {
				appendContent += "\n"
			}
			appendContent += strings.Join(toAppend, "\n") + "\n"
			f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				logger.Errorf("[generateConfigTemplate] 打开 .gitignore 失败: %v", err)
				return fmt.Errorf("打开 .gitignore 失败: %w", err)
			}
			defer f.Close()
			if _, err := f.WriteString(appendContent); err != nil {
				logger.Errorf("[generateConfigTemplate] 写入 .gitignore 失败: %v", err)
				return fmt.Errorf("写入 .gitignore 失败: %w", err)
			}
			logger.ColorPrintf(logger.ColorHiGreen, "  ✓ 已向 %s 追加: %s\n", gitignorePath, strings.Join(toAppend, ", "))
		} else {
			logger.ColorPrintf(logger.ColorWhite, "  ⏭ %s 已包含所需条目，跳过\n", gitignorePath)
		}
	}

	fmt.Println()
	logger.ColorPrintf(logger.ColorHiGreen, "✓ 初始化完成！\n")
	fmt.Println()
	logger.ColorPrintf(logger.ColorHiCyan, "文件列表:\n")
	logger.ColorPrintf(logger.ColorWhite, "  📄 %s  — 配置文件\n", configFilePath)
	logger.ColorPrintf(logger.ColorWhite, "  📄 %s  — SQL 建表语句\n", schemaFilePath)
	logger.ColorPrintf(logger.ColorWhite, "  📄 %s  — Git 忽略规则\n", gitignorePath)
	fmt.Println()
	logger.ColorPrintf(logger.ColorHiCyan, "使用说明:\n")
	logger.ColorPrintf(logger.ColorWhite, "  1. 在 %s 中编写你的 CREATE TABLE 语句\n", schemaFilePath)
	logger.ColorPrintf(logger.ColorWhite, "  2. 编辑 %s 调整生成选项（可选）\n", configFilePath)
	logger.ColorPrintf(logger.ColorWhite, "  3. 运行 jen 命令生成代码，结果输出到 target/jen/ 目录\n")

	return nil
}
