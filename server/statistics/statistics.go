/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 18:14
 */

package statistics

import (
	"fmt"
	"go-stress-testing-pool/model"
	"strings"
	"sync"
	"time"
)

var (
	// 输出统计数据的时间
	exportStatisticsTime = 1 * time.Second
)

// 接收结果并处理
// 统计的时间都是纳秒，显示的时间 都是毫秒
// concurrent 并发数
func ReceivingResults(concurrent uint64, requestResultCh <-chan *model.RequestResults, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	var (
		stopChan = make(chan bool)
	)

	// 时间
	var (
		processingTime uint64 // 处理总时间
		requestTime    uint64 // 请求总时间
		maxTime        uint64 // 最大时长
		minTime        uint64 // 最小时长
		successNum     uint64 // 成功处理数，code为0
		failureNum     uint64 // 处理失败数，code不为0
		chanIdLen      int    // 并发数 以id去重
		chanIds        = make(map[uint64]bool)
	)

	statTime := uint64(time.Now().UnixNano())

	// 错误码/错误个数
	var errCode = make(map[int]int)

	// 定时输出一次计算结果
	ticker := time.NewTicker(exportStatisticsTime)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				requestTime = endTime - statTime
				//chanIdLen 并发数
				go calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIdLen, errCode)
			case <-stopChan:
				// 处理完成

				return
			}
		}
	}()

	header()

	for data := range requestResultCh {
		// fmt.Println("处理一条数据", data.Id, data.Time, data.IsSucceed, data.ErrCode)
		processingTime = processingTime + data.Time

		if maxTime <= data.Time {
			maxTime = data.Time
		}

		if minTime == 0 {
			minTime = data.Time
		} else if minTime > data.Time {
			minTime = data.Time
		}

		// 是否请求成功
		if data.IsSucceed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}

		// 统计错误码
		if value, ok := errCode[data.ErrCode]; ok {
			errCode[data.ErrCode] = value + 1
		} else {
			errCode[data.ErrCode] = 1
		}

		if _, ok := chanIds[data.ChanId]; !ok {
			chanIds[data.ChanId] = true
			chanIdLen = len(chanIds)
		}
	}

	// 数据全部接受完成，停止定时输出统计数据
	stopChan <- true

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - statTime

	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIdLen, errCode)

	fmt.Printf("\n\n")

	fmt.Println("*************************  结果 stat  ****************************")
	fmt.Println("处理协程数量:", concurrent)
	// fmt.Println("处理协程数量:", concurrent, "程序处理总时长:", fmt.Sprintf("%.3f", float64(processingTime/concurrent)/1e9), "秒")
	fmt.Println("请求总数:", successNum+failureNum, "总请求时间:", fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"秒", "successNum:", successNum, "failureNum:", failureNum)

	fmt.Println("*************************  结果 end   ****************************")

	fmt.Printf("\n\n")
}

// 计算数据
func calculateData(concurrentToal, processingTime, requestTime, maxTime, minTime, successNum, failureNum uint64, chanIdLen int, errCode map[int]int) {
	if processingTime == 0 {
		processingTime = 1
	}

	var (
		qps              float64
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)

	//QPS = req/sec = 请求数/秒 1e9即是数10的9次方
	// 平均 每个协程成功数*总协程数据/总耗时 (每秒)
	//fmt.Printf("concurrent %7d",concurrent)
	//fmt.Printf("processingTime %7d ",processingTime)
	//fmt.Printf("请求数 %8.2f ",float64(successNum*1e9*successNum))
	//fmt.Printf("processingTime %8.2f ",float64(processingTime))

	/*chanIdLenStr := strconv.Itoa(chanIdLen)
	chanIdLenIntNum, _ := strconv.Atoi(chanIdLenStr)
	chanIdLenInt64 := uint64(chanIdLenIntNum)*/
	if processingTime != 0 {
		//qps = float64(successNum*1e9*concurrent) / float64(processingTime)
		qps = float64(successNum*1e9*successNum)/ float64(processingTime)
		//qps = float64(successNum+failureNum) / float64(processingTime/1e9)
	}


	// 平均时长 总耗时/总请求数/并发数 纳秒=>毫秒
	/*if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6*concurrent)
	}*/
	if successNum != 0 && (successNum+failureNum) != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6)
	}

	// 纳秒=>毫秒
	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	requestTimeFloat = float64(requestTime) / 1e9
	processingTimeFloat := 0.0
	if successNum+failureNum != 0  {
		processingTimeFloat = float64(processingTime) / float64(successNum*1e9)
	}

	// 打印的时长都为毫秒
	// result := fmt.Sprintf("请求总数:%8d|successNum:%8d|failureNum:%8d|qps:%9.3f|maxTime:%9.3f|minTime:%9.3f|平均时长:%9.3f|errCode:%v", successNum+failureNum, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, errCode)
	// fmt.Println(result)
	table(successNum, failureNum, errCode, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat,processingTimeFloat, chanIdLen)
}

// 打印表头信息
func header() {
	fmt.Printf("\n\n")
	// 打印的时长都为毫秒 总请数
	fmt.Println("─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────")
	result := fmt.Sprintf(" 请求耗时│ 并发数│ 成功数│ 失败数│   qps  │最长耗时│最短耗时│平均耗时│ 错误码")
	fmt.Println(result)
	// result = fmt.Sprintf("耗时(s)  │总请求数│成功数│失败数│QPS│最长耗时│最短耗时│平均耗时│错误码")
	// fmt.Println(result)
	fmt.Println("─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────")

	return
}

// 打印表格
func table(successNum, failureNum uint64, errCode map[int]int, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64,processingTimeFloat float64, chanIdLen int) {
	// 打印的时长都为毫秒
	result := fmt.Sprintf("请求耗时 %4.0fs│处理耗时 %4.0fs│完成数 %7d│成功数 %7d│失败数 %7d│qps %8.2f│最长耗时ms %8.2f│最短耗时ms %8.2f│平均耗时ms %8.2f│错误码 %v",
		requestTimeFloat, processingTimeFloat,chanIdLen, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, printMap(errCode))
	fmt.Println(result)

	return
}

// 输出错误码、次数 节约字符(终端一行字符大小有限)
func printMap(errCode map[int]int) (mapStr string) {

	var (
		mapArr []string
	)
	for key, value := range errCode {
		mapArr = append(mapArr, fmt.Sprintf("%d:%d", key, value))
	}

	mapStr = strings.Join(mapArr, ";")

	return
}
