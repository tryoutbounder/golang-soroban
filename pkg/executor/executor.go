package executor

import (
	"net/http"

	"github.com/tryoutbounder/golang-soroban/pkg/rpc"
)

type Executor struct {
	rpc *rpc.RpcClient
}

func NewExecutor(rpcUrl string, httpClient *http.Client) *Executor {
	return &Executor{
		rpc: rpc.NewClient(rpcUrl, httpClient),
	}
}

func (e *Executor) GetRpc() *rpc.RpcClient {
	return e.rpc
}
