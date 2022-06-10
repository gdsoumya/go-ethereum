package rpc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"io"
	"net/http"
)

func GetAddress(ip, port string) (*common.Address, error) {
	client, err := rpc.Dial(fmt.Sprintf("http://%v:%v", ip, port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node, error=%w", err)
	}
	resp := []common.Address{}
	if err = client.Call(&resp, "personal_listAccounts"); err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		return &resp[0], nil
	}
	return nil, fmt.Errorf("failed to get validator address from node")
}

func GetLatestBlock(ctx context.Context, ip, port string) (uint64, error) {
	client, err := ethclient.Dial(fmt.Sprintf("http://%v:%v", ip, port))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to node, error=%w", err)
	}
	return client.BlockNumber(ctx)
}

func GenericRpcCall(ip, port string, buff []byte) ([]byte, error) {
	resp, err := http.DefaultClient.Post(fmt.Sprintf("http://%v:%v", ip, port), "application/json", bytes.NewBuffer(buff))
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
