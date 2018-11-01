package model

import "github.com/davyxu/cellnet"

var (
	chanByName = map[string][]cellnet.Session{}
)

func AddSubscriber(name string, ses cellnet.Session) {

	list, _ := chanByName[name]
	list = append(list, ses)
	chanByName[name] = list
}

func RemoveSubscriber(ses cellnet.Session, callback func(chanName string)) {

	var found bool
	for {

		found = false

	Refound:
		for name, list := range chanByName {

			for index, libSes := range list {
				if libSes == ses {

					callback(name)
					list = append(list[:index], list[index+1:]...)

					if len(list) == 0 {
						delete(chanByName, name)
					} else {
						chanByName[name] = list
					}

					found = true
					goto Refound // 避免循环删除造成的故障,再查一次
				}
			}
		}

		// 直到没有删的了
		if !found {
			break
		}

	}

	return
}

func VisitSubscriber(name string, callback func(ses cellnet.Session) bool) (count int) {
	if list, ok := chanByName[name]; ok {

		for _, ses := range list {
			count++
			if !callback(ses) {
				return
			}
		}
	}

	return

}
