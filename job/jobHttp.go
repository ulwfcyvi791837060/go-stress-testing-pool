package job

import (
	"go-stress-testing/model"
	"go-stress-testing/server/golink"
	"sync"
)

func RunHttp(param []interface{}) {
	var i uint64
	var ch chan *model.RequestResults
	var totalNumber uint64
	var wg *sync.WaitGroup
	var request *model.Request


	var paramMap map[string]interface{}
	var resultChan chan interface{}
	for _, val := range param {
		switch v := val.(type) {
		case map[string]interface{}:
			paramMap = v
		case chan interface{}:
			resultChan = v
		}
	}

	i = paramMap["i"].(uint64)
	ch = paramMap["ch"].(chan *model.RequestResults)
	totalNumber = paramMap["totalNumber"].(uint64)
	wg = paramMap["wg"].(*sync.WaitGroup)
	request = paramMap["request"].(*model.Request)

	cTest(i, ch, totalNumber, wg, request)

	if resultChan != nil {
		//resultChan <- c
	} else {
		//fmt.Println(c)
	}

	//time.Sleep(time.Millisecond*10);
}
func cTest( i uint64, ch chan *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request)  {
	golink.Http(i, ch, totalNumber, wg, request)
	return
}
