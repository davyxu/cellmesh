module github.com/davyxu/cellmesh

go 1.12

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/davyxu/cellmesh_demo v0.0.0-20190820092800-96af1dea1c93 // indirect
	github.com/davyxu/cellnet v0.0.0-20190628065413-a644d2409b6d
	github.com/davyxu/golog v0.1.0
	github.com/davyxu/protoplus v0.1.0
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df // indirect
	github.com/ouqiang/supervisor-event-listener v0.0.0-20180320124003-031edd705fcd // indirect
	gopkg.in/ini.v1 v1.49.0 // indirect
)

// 本地修改cellmesh时使用
replace github.com/davyxu/protoplus => ../protoplus

replace github.com/davyxu/cellnet => ../cellnet
