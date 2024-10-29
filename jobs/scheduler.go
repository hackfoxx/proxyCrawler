package jobs

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
)

func GetScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler()
	// 初始化
	fmt.Println("初始运行扫描器")
	CrawlersJob()
	fmt.Println("初始运行验证器")
	ValidatorJob()
	err := s.Every(5).Minute().Do(CrawlersJob)
	if err != nil {
		log.Fatalf("Failed to schedule CrawlersJob: %v", err) // 记录错误
	}
	err = s.Every(5).Minute().Do(ValidatorJob)
	return s
}
