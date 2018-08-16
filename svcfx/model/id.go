package fxmodel

var (
	IDTail string
)

func GetSvcID(svcName string) string {
	return svcName + "_" + IDTail
}
