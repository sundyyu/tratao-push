package check

import (
	"time"
)

/**
 * @desc 定期器循环运行预警信息检查
 * @author 于朝鹏
 * @date 2018年10月8日 15:20
 */
func TickLoop(check Check, tickTime int) {
	go func() {
		c := check
		t := tickTime
		ticker := time.NewTicker(time.Second * time.Duration(t))
		for {
			select {
			case <-ticker.C:
				c.Update()
			}
		}
	}()
}
