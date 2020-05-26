module github.com/davyxu/cellmesh

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davyxu/cellnet v0.0.0-20190628065413-a644d2409b6d
	github.com/davyxu/protoplus v0.1.0
	github.com/davyxu/tabtoy v0.0.0-20200515034133-f40af96ceda7
	github.com/davyxu/ulog v1.0.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/kr/pretty v0.2.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/davyxu/ulog => ../ulog

replace github.com/davyxu/protoplus => ../protoplus

replace github.com/davyxu/cellnet => ../cellnet
