package discovery

// 只过滤出需要的结果
func FilterByTag(sdList []*ServiceDesc, tags ...string) (ret []*ServiceDesc) {

	if len(sdList) == 0 {
		return
	}

nextSD:
	for _, sd := range sdList {

		for _, sdTag := range sd.Tags {
			for _, tag := range tags {
				if sdTag == tag {
					ret = append(ret, sd)
					continue nextSD
				}
			}
		}

	}

	return
}
