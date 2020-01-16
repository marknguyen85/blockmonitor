package helpers

import (
	"errors"
	"math/big"

	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"telegramalert/node"
)

type Blockchain struct {
	Client      *rpc.Client
	LatestBlock *big.Int
}

func NewBlcClientFromFlags(ctx *cli.Context) (*Blockchain, error) {
	client, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	blcClient := &Blockchain{
		Client:      client,
		LatestBlock: new(big.Int).SetUint64(0),
	}
	return blcClient, nil
}

func (blc *Blockchain) GetLastBlock() (*big.Int, error) {
	var header *types.Header
	err := blc.Client.Call(&header, "evr_getBlockByNumber", "latest", true)
	if err != nil {
		return nil, err
	}
	if header == nil {
		return nil, errors.New("can not get latest block")
	}
	return header.Number, nil
}
