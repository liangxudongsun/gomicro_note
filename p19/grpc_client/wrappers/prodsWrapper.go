package wrappers

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
	"gomicro_note/p19/grpc_client/models"
	"strconv"
)

func newProd(id int32, pname string) *models.ProdModel {
	return &models.ProdModel{ProdId: id, ProdName: pname}
}

func defaultProds(rsp interface{}) (*models.ProdListResponse, error) {
	model := make([]*models.ProdModel, 0)
	var i int32
	for i = 0; i < 5; i++ {
		model = append(model, newProd(20+i, "prodName"+strconv.Itoa(20+int(i))))
	}
	res := rsp.(*models.ProdListResponse)
	res.Data = model
	return res, nil
}

// ProdsWrapper 商品装饰器
type ProdsWrapper struct {
	client.Client
}

// Call 调用方法
func (p *ProdsWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	configA := hystrix.CommandConfig{Timeout: 1000}
	hystrix.ConfigureCommand(cmdName, configA)
	return hystrix.Do(cmdName, func() error {
		return p.Client.Call(ctx, req, rsp)
	}, func(e error) error {
		defaultProds(rsp)
		return nil
	})
	return p.Client.Call(ctx, req, rsp)
}

// NewProdsWrapper 初始化一个商品装饰器
func NewProdsWrapper(c client.Client) client.Client {
	return &ProdsWrapper{c}
}
