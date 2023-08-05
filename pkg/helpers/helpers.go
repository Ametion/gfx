package helpers

import (
	"github.com/Ametion/GFX/gfx"
)

func Contains(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func HandlerToMiddleware(h gfx.HandlerFunc) gfx.MiddlewareFunc {
	return func(c *gfx.Context) {
		h(c)
	}
}
