package api

import (
	"context"
	"github.com/goiiot/libmqtt/benchmark/ga"
	"github.com/goiiot/libmqtt/common"
	"sync"
	"testing"
	"time"
)

func BenchmarkApi(b *testing.B) {
	group := &sync.WaitGroup{}

	for i := 0; i < 600; i++ {
		time.Sleep(time.Second)
		for i := 0; i < 20; i++ {
			group.Add(4)
			go requestProductList(group)
			go requestDeviceList(group)
			go requestLog(group)
			go requestHistory(group)
		}
	}

	group.Wait()
}

func requestHistory(group *sync.WaitGroup) {
	reqHistory := ga.TestAddress + "/device/instance/1800843366925897728/properties/_query?currentPage=0&pageIndex=0&pageSize=60&sorts%5B0%5D.name=timestamp&sorts%5B0%5D.order=desc&terms%5B0%5D.column=property&terms%5B0%5D.value=a2&terms%5B1%5D.column=timestamp&terms%5B1%5D.value=1711514130000&terms%5B1%5D.termType=gte&terms%5B2%5D.column=timestamp&terms%5B2%5D.value=1719462930999&terms%5B2%5D.termType=lte"
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer func() {
		group.Done()
		cancelFunc()
	}()
	common.DoRequestWithTimeout("GET", reqHistory, "", ga.TestToken, ctx)
}

func requestLog(group *sync.WaitGroup) {
	reqLog := ga.TestAddress + "/logger/access/_query?currentPage=0&pageIndex=0&pageSize=30&terms%5B0%5D.terms%5B0%5D.column=requestTime&terms%5B0%5D.terms%5B0%5D.value=1718899200000&terms%5B0%5D.terms%5B0%5D.type=and&terms%5B0%5D.terms%5B0%5D.termType=gt&terms%5B0%5D.terms%5B1%5D.column=requestTime&terms%5B0%5D.terms%5B1%5D.value=1719503999999&terms%5B0%5D.terms%5B1%5D.type=and&terms%5B0%5D.terms%5B1%5D.termType=lt&sorts%5B0%5D.name=requestTime&sorts%5B0%5D.order=desc&terms%5B0%5D.column=requestTime&terms%5B0%5D.value=1719458154285&terms%5B0%5D.termType=lt"
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer func() {
		group.Done()
		cancelFunc()
	}()
	common.DoRequestWithTimeout("GET", reqLog, "", ga.TestToken, ctx)
}

func requestDeviceList(group *sync.WaitGroup) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer func() {
		group.Done()
		cancelFunc()
	}()

	productList := ga.TestAddress + "/v3/device/product/queryPager"
	queryBody := "{\n  \"terms\": [\n    {\n      \"column\": \"productId\",\n      \"value\": \"0\",\n      \"termType\": \"dev-product-type\"\n    }\n  ],\n  \"paging\": false\n}"
	common.DoRequestWithTimeout("POST", productList, queryBody, ga.TestToken, ctx)
}

func requestProductList(group *sync.WaitGroup) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer func() {
		group.Done()
		cancelFunc()
	}()

	productList := ga.TestAddress + "/v3/device/instance/_query"
	queryBody := "{\n  \"pageSize\": 15,\n  \"currentPage\": 1,\n  \"sorts\": [\n    {\n      \"name\": \"createTime\",\n      \"order\": \"desc\"\n    }\n  ],\n  \"terms\": [\n    {\n      \"column\": \"id\",\n      \"value\": \"all\",\n      \"termType\": \"dev-group\"\n    },\n    {\n      \"column\": \"productId\",\n      \"value\": \"0\",\n      \"termType\": \"dev-product-type\"\n    }\n  ]\n}"
	common.DoRequestWithTimeout("POST", productList, queryBody, ga.TestToken, ctx)

}
