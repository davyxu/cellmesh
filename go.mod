module github.com/davyxu/cellmesh

go 1.12

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/davyxu/cellnet v0.0.0-00010101000000-000000000000
	github.com/davyxu/golexer v0.0.0-20180314091252-f048a86ae200
	github.com/davyxu/golog v0.1.0
	github.com/davyxu/protoplus v0.1.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/satori/go.uuid v1.2.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/davyxu/protoplus => ../protoplus

replace github.com/davyxu/golog => ../golog

replace github.com/davyxu/golexer => ../golexer

replace github.com/davyxu/cellnet => ../cellnet
