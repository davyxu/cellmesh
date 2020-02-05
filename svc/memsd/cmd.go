package main

import (
	"fmt"
	"os"
	"sort"
)

func ViewSvc() {

	sd := initSD()

	list := sd.QueryAll()

	sort.Slice(list, func(i, j int) bool {

		a := list[i]
		b := list[j]

		if a.Port != b.Port {
			return a.Port < b.Port
		}

		if a.Host != b.Host {
			return a.Host < b.Host
		}

		return a.ID < b.ID
	})

	for _, desc := range list {

		fmt.Println(desc.VisualString())
	}
}

func ViewKey() {
	sd := initSD()
	list := sd.GetRawValueList("")
	sort.Slice(list, func(i, j int) bool {

		a := list[i]
		b := list[j]

		return a.Key < b.Key
	})

	for _, meta := range list {
		fmt.Printf("  %s = (size %d)\n", meta.Key, len(meta.Value))
	}
}

func GetValue(key string) {
	sd := initSD()
	var value string
	err := sd.GetValue(key, &value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(value)
}

func SetValue(key, value string) {
	sd := initSD()
	err := sd.SetValue(key, value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ClearSvc() {

	sd := initSD()
	sd.ClearService()
}

func ClearValue() {

	sd := initSD()
	sd.ClearKey()
}

func DeleteValue(key string) {

	sd := initSD()
	sd.DeleteValue(key)
}
