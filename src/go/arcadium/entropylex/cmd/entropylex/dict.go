package main

import (
	"embed"
)

//go:embed dict/*
var dictFS embed.FS
