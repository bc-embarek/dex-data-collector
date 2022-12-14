package dex

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"testing"
)

var collector Collector

func init() {
	if rpcClient, err := rpc.Dial("https://eth.public-rpc.com"); err != nil {
		log.Fatalf("connect rpc endpoint failed with error message:%s", err.Error())
	} else {
		collector = Collector{
			GethClient: gethclient.New(rpcClient),
		}
	}

}

func TestGetPairs(t *testing.T) {
	// uniswap v2 factory
	pairs := collector.GetAllPairs(common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"), false)
	fmt.Println(pairs)
}

func TestGetTokens(t *testing.T) {
	tokens := collector.getTokens([]common.Address{common.HexToAddress("0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599")})
	fmt.Println(tokens)
}
