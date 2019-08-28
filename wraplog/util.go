package wraplog

import (
	"net/url"
)

//SetPrefix implements encode prefix to log
func GetCtx(encodedCtx, newUnEscapedCtx string) string {
	newEncodeCtx := url.QueryEscape("->" + newUnEscapedCtx)
	return encodedCtx + newEncodeCtx
}
