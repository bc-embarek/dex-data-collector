package main

import (
	"context"
	"dex-data-collector/dex"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	rpcEndpoint      string
	factoryAddress   string
	includeTokenInfo bool

	ethClient  *ethclient.Client
	gethClient *gethclient.Client
)

func parseArgs() {
	flag.StringVar(&factoryAddress, "f", "", "factory contract address. eg: -f 0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	flag.StringVar(&rpcEndpoint, "r", "", "rpc endpoint: support wss or http. eg: -r https://eth.public-rpc.com")
	flag.BoolVar(&includeTokenInfo, "i", false, "include token info. eg -i false")
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range []string{"r", "f"} {
		if !seen[req] {
			log.Fatalf("missing required -%s argument", req)
		}
	}
}

func buildEthClient() {
	if rpcClient, err := rpc.Dial(rpcEndpoint); err != nil {
		log.Fatalf("connect rpc endpoint failed with error message:%s", err.Error())
	} else {
		ethClient = ethclient.NewClient(rpcClient)
		gethClient = gethclient.New(rpcClient)
	}
}

func init() {
	parseArgs()
	buildEthClient()
}

func main() {
	// get pair data from rpc
	collector := dex.Collector{GethClient: gethClient}
	pairs := collector.GetAllPairs(common.HexToAddress(factoryAddress), includeTokenInfo)

	// export to csv file
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		log.Fatal("get chain_id failed with error: ", err)
	}
	exportFilename := fmt.Sprintf("chain-%d-factory-%s-%v.csv", chainID, factoryAddress, time.Now().UTC().UnixNano()/1000000)
	file, err := os.Create(exportFilename)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("close file failed with error: ", err)
		}
	}()
	if err != nil {
		log.Fatal("open file failed with error: ", err)
	}

	exportFile := csv.NewWriter(file)
	defer exportFile.Flush()
	if err := exportFile.Write([]string{"pair address", "token0 address", "token1 address", "token0 symbol", "token0 name", "token0 decimals", "token1 symbol", "token1 name", "token1 decimals"}); err != nil {
		log.Fatal("write header to file failed with error: ", err)
	}

	var records [][]string
	for _, pair := range pairs {
		record := []string{pair.Address.Hex(), pair.Token0Address.Hex(), pair.Token1Address.Hex()}

		if pair.Token0 != nil {
			record = append(record, pair.Token0.Symbol, pair.Token0.Name, strconv.Itoa(int(pair.Token0.Decimals)))
		} else {
			record = append(record, "", "", "")
		}
		if pair.Token1 != nil {
			record = append(record, pair.Token1.Symbol, pair.Token1.Name, strconv.Itoa(int(pair.Token1.Decimals)))
		} else {
			record = append(record, "", "", "")
		}
		records = append(records, record)
	}

	if err := exportFile.WriteAll(records); err != nil {
		log.Fatal("export record to file failed with error: ", err)
	}
}
