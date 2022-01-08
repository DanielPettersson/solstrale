package trace

import (
	"math/rand"
	"time"
)

func RayTrace(width int, height int) []byte {

	ret := make([]byte, width*height*4)

	rand.Seed(time.Now().UnixNano())
	rand.Read(ret)

	return ret

}
