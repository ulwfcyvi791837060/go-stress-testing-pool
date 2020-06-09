/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"go-stress-testing-pool/heper"
	"go-stress-testing-pool/model"
	"go-stress-testing-pool/server/client"
	"sync"
	"time"
)

// http go link
func Http(chanId uint64, requestResultCh chan<- *model.RequestResults, cycleNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < cycleNumber; i++ {

		var (
			startTime = time.Now()
			//成功数加1
			isSucceed = false
			errCode   = model.HttpOk
		)

		resp, err := client.HttpRequest(request.Method, request.Url, request.GetBody(), request.Headers, request.Timeout)

		//结束时间- 开始时间
		requestTime := uint64(heper.DiffNano(startTime))
		//fmt.Printf("requestTime:%05d \n", requestTime)
		// resp, err := server.HttpGetResp(request.Url)
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			// 验证请求是否成功
			errCode, isSucceed = request.VerifyHttp(request, resp)
		}

		requestResults := &model.RequestResults{
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}

		requestResults.SetId(chanId, i)

		requestResultCh <- requestResults
	}

	return
}
