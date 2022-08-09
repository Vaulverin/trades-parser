package logparser

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronjan/hunch"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/vaulverin/trades-parser/contracts"
	"github.com/vaulverin/trades-parser/units"
	"math/big"
)

type LogParser struct {
	client   *ethclient.Client
	address  common.Address
	contract *contracts.UniswapV2Pair
	decimals [2]uint8
}

func New(client *ethclient.Client, address common.Address, decimals [2]uint8) (*LogParser, error) {
	contract, err := contracts.NewUniswapV2Pair(address, client)
	if err != nil {
		return nil, err
	}
	return &LogParser{client: client, address: address, contract: contract, decimals: decimals}, nil
}

const (
	Buy = iota
	Sell
)

type Swap struct {
	Block     uint64
	Timestamp uint64
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Side      int
}

func (lp *LogParser) ParseSwaps(ctx context.Context, from uint64, to *uint64) ([]Swap, error) {
	currentBlock, err := lp.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	if from > currentBlock {
		return nil, errors.New(fmt.Sprintf("Block number 'from':%d is grater than current block %d", from, currentBlock))
	}
	if to != nil && *to > currentBlock {
		return nil, errors.New(fmt.Sprintf("Block number 'to':%d is grater than current block %d", to, currentBlock))
	}
	if to != nil && *to < from {
		return nil, errors.New(fmt.Sprintf("Block number 'from':%d is grater than 'to' block %d", from, to))
	}
	var swaps []Swap
	endBlock := currentBlock
	if to != nil {
		endBlock = *to
	}
	step := uint64(1000)
	for start := from; start < endBlock; start += step {
		end := start + step
		if end > endBlock {
			end = endBlock
		}
		parsed, err := hunch.Retry(ctx, 3, func(ctx context.Context) (interface{}, error) {
			return lp.parse(ctx, start, &end)
		})
		if err != nil {
			return nil, err
		}
		swaps = append(swaps, parsed.([]Swap)...)
	}

	return swaps, nil
}

func (lp *LogParser) parse(ctx context.Context, from uint64, to *uint64) ([]Swap, error) {
	iterator, err := lp.contract.FilterSwap(&bind.FilterOpts{
		Start:   from,
		End:     to,
		Context: ctx,
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	var swaps []Swap
	cache := map[uint64]uint64{}
	for iterator.Next() {
		event := iterator.Event
		timestamp, ok := cache[event.Raw.BlockNumber]
		if !ok {
			blockInfo, err := lp.client.BlockByNumber(ctx, big.NewInt(int64(event.Raw.BlockNumber)))
			if err != nil {
				return nil, err
			}
			timestamp = blockInfo.Time()
			cache[event.Raw.BlockNumber] = timestamp
		}
		swap := Swap{Block: event.Raw.BlockNumber, Timestamp: timestamp}
		if event.Amount0In.Cmp(big.NewInt(0)) == 0 {
			swap.Side = Buy
			swap.Amount = units.ToDecimal(event.Amount0Out, lp.decimals[0])
			swap.Price = units.ToDecimal(event.Amount1In, lp.decimals[1]).Div(swap.Amount)
		} else {
			swap.Side = Sell
			swap.Amount = units.ToDecimal(event.Amount0In, lp.decimals[0])
			swap.Price = units.ToDecimal(event.Amount1Out, lp.decimals[1]).Div(swap.Amount)
		}
		swaps = append(swaps, swap)
	}

	err = iterator.Close()
	if err != nil {
		return nil, err
	}

	return swaps, err
}
