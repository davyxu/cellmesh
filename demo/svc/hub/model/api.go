package model

import "github.com/davyxu/cellnet"

var (
	HubSession cellnet.Session
	OnHubReady func()
)
