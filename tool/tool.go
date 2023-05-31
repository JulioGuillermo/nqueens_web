package tool

import (
	"fmt"
	"time"
)

func GetDurStr(d time.Duration) string {
	return fmt.Sprintf(
		"%d:%d:%d:%d",
		int(d.Hours()),
		int(d.Minutes())%60,
		int(d.Seconds())%60,
		int(d.Milliseconds())%1000,
	)
}
