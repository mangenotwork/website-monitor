package handler

import (
	"github.com/mangenotwork/common/ginHelper"
)

func Index(c *ginHelper.GinCtx) {
	c.APIOutPut("", "ok")
	return
}
