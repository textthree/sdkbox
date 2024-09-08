package alioss

import "time"

type OssDirList struct {
	Path           string
	DirName        string
	LastModifyTime time.Time
}

type Objects struct {
	Path       string
	ObjectName string
	AccessUrl  string
}

const SortTypeModifyTimeAsc = 1
const SortTypeModifyTimeDesc = 2
const SortTypeNumber = 3 // 按第一个数字。如：1.xxx.jpg 取第一个点之前的数字
