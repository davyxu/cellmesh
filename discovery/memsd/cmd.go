package main

import (
	"fmt"
	"os"
	"sort"
)

func viewSvc() {

	sd := initSD()

	list := sd.QueryAll()

	sort.Slice(list, func(i, j int) bool {

		a := list[i]
		b := list[j]

		if a.GetMeta("SvcGroup") != b.GetMeta("SvcGroup") {
			return a.GetMeta("SvcGroup") < b.GetMeta("SvcGroup")
		}

		if a.Port != b.Port {
			return a.Port < b.Port
		}

		if a.Host != b.Host {
			return a.Host < b.Host
		}

		return a.ID < b.ID
	})

	for _, desc := range list {

		fmt.Println(desc.FormatString())
	}
}

func viewKey() {
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

func getValue(key string) {
	sd := initSD()
	var value string
	err := sd.GetValue(key, &value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(value)
}

func setValue(key, value string) {
	sd := initSD()
	err := sd.SetValue(key, value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func clearSvc() {

	sd := initSD()
	sd.ClearService()
}

func clearKey() {

	sd := initSD()
	sd.ClearKey()
}
