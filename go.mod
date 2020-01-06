module github.com/davyxu/cellmesh

go 1.12

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/davyxu/cellnet v0.0.0-20190628065413-a644d2409b6d
	github.com/davyxu/golog v0.1.0
	github.com/davyxu/protoplus v0.1.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/satori/go.uuid v1.2.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

// 本地修改cellmesh时使用
replace github.com/davyxu/protoplus => ../protoplus

replace github.com/davyxu/cellnet => ../cellnet
