package generator

import "embed"

//go:embed template/*
var fs embed.FS

const templatePathPrefix = "template/"
