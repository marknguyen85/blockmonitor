package helpers

import (
	"errors"
	"math/big"
	"time"

	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"telegramalert/node"
)

var (
	timeTickerFlag = cli.DurationFlag{
		Name:  "duration",
		Usage: "The duration delay for each times to checks the situation of block (seconds)",
		Value: 60,
	}
)

type Blockchain struct {
	Client      *rpc.Client
	LatestBlock *big.Int
	Duration    time.Duration
}

// NewBlcClientFlag returns flags for block-chain
func NewBlcClientFlag() []cli.Flag {
	return []cli.Flag{timeTickerFlag}
}

func NewBlcClientFromFlags(ctx *cli.Context) (*Blockchain, error) {
	var (
		delay = ctx.Duration(timeTickerFlag.Name)
	)

	client, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	blcClient := &Blockchain{
		Client:      client,
		LatestBlock: new(big.Int).SetUint64(0),
		Duration:    delay,
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
