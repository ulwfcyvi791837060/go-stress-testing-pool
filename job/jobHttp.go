package job

import (
	"go-stress-testing-pool/model"
	"go-stress-testing-pool/server/golink"
	"sync"
)

func RunHttp(param []interface{}) {
	var chanId uint64
	var ch chan *model.RequestResults
	var cycleNumber uint64
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

	chanId = paramMap["chanId"].(uint64)
	ch = paramMap["ch"].(chan *model.RequestResults)
	cycleNumber = paramMap["cycleNumber"].(uint64)
	wg = paramMap["wg"].(*sync.WaitGroup)
	request = paramMap["request"].(*model.Request)

	cTest(chanId, ch, cycleNumber, wg, request)

	if resultChan != nil {
		//resultChan <- c
	} else {
		//fmt.Println(c)
	}

	//time.Sleep(time.Millisecond*10);
}
func cTest( chanId uint64, ch chan *model.RequestResults, cycleNumber uint64, wg *sync.WaitGroup, request *model.Request)  {
	golink.Http(chanId, ch, cycleNumber, wg, request)
	return
}
