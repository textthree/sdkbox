package alioss

import "time"

type OssDirList struct {
	Path           string
	DirName        string
	LastModifyTime time.Time
}

type Objects struct {
	Path           string
	ObjectName     string
	AccessUrl     string
}