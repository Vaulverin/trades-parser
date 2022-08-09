package logparser_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vaulverin/trades-parser/logparser"
	"testing"
)

func TestLogParser_ParseSwaps(t *testing.T) {
	rpcUrl := "https://rpc.ankr.com/eth"
	pairAddr := common.HexToAddress("0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852")
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		t.Fatal(err)
	}
	parser, err := logparser.New(client, pairAddr, [2]uint8{18, 6})
	if err != nil {
		t.Fatal(err)
	}
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Wrong parameters", func(t *testing.T) {
		_, err := parser.ParseSwaps(context.Background(), currentBlock+10, nil)
		if err == nil {
			t.Error("Must return an error")
		}
	})

	t.Run("Wrong parameters", func(t *testing.T) {
		from := uint64(10)
		to := uint64(1)
		_, err := parser.ParseSwaps(context.Background(), from, &to)
		if err == nil {
			t.Error("Must return an error")
		}
	})

	t.Run("Wrong parameters", func(t *testing.T) {
		from := currentBlock - 10
		to := currentBlock + 10
		_, err := parser.ParseSwaps(context.Background(), from, &to)
		if err == nil {
			t.Error("Must return an error")
		}
	})

	t.Run("Current block", func(t *testing.T) {
		swaps, err := parser.ParseSwaps(context.Background(), currentBlock, nil)
		if err != nil {
			t.Error("Failed to fetch swaps:", err)
		}
		t.Log(swaps)
	})

	t.Run("From to", func(t *testing.T) {
		swaps, err := parser.ParseSwaps(context.Background(), currentBlock-100, &currentBlock)
		if err != nil {
			t.Error("Failed to fetch swaps:", err)
		}
		t.Log(swaps)
	})

	//t.Run("From to with long range", func(t *testing.T) {
	//	to := currentBlock - 5000
	//	swaps, err := parser.ParseSwaps(context.Background(), currentBlock-10000, &to)
	//	if err != nil {
	//		t.Error("Failed to fetch swaps:", err)
	//	}
	//	t.Log(swaps)
	//})
	//
	//t.Run("All swaps", func(t *testing.T) {
	//	swaps, err := parser.ParseSwaps(context.Background(), 0, nil)
	//	if err != nil {
	//		t.Error("Failed to fetch swaps:", err)
	//	}
	//	t.Log("Count:", len(swaps))
	//})
}
