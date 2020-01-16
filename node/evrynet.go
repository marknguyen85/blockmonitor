package node

import (
	"errors"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"
)

const (
	rpcEndpointFlag = "rpcendpoint"
	evrynetEndpoint = "http://0.0.0.0:8545"
)

// EvrynetEndpoint returns configured Evrynet node endpoint.
func EvrynetEndpoint() string {
	return evrynetEndpoint
}

// NewEvrynetNodeFlags return flags to EvrynetNode
func NewEvrynetNodeFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  rpcEndpointFlag,
			Usage: "RPC endpoint to send request",
			Value: EvrynetEndpoint(),
		}}
}

// NewEvrynetClientFromFlags returns Evrynet client from flag variable, or error if occurs
func NewEvrynetClientFromFlags(ctx *cli.Context) (*rpc.Client, error) {
	evrynetClientURL := ctx.String(rpcEndpointFlag)
	client, err := rpc.Dial(evrynetClientURL)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("can not connects to: " + rpcEndpointFlag)
	}
	return client, nil
}
