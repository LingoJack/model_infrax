package generator

import "embed"

// templateFS 嵌入所有模板文件到二进制中
// 使用 embed.FS 可以保持目录结构，方便按路径读取
// 路径相对于当前 go 文件所在目录（generator 目录）
// 由于模板文件在 ../assert/template 目录下，需要使用相对路径引用
//
//go:embed ../assert/template/*.template
//go:embed ../assert/template/itea-go/*.template
//go:embed ../assert/template/tools/*.template
var templateFS embed.FS

// 模板文件路径前缀，用于从嵌入的文件系统中读取文件
// 由于 embed 是从 generator 目录开始的，所以需要加上 ../assert/template/ 前缀
const templatePathPrefix = "../assert/template/"