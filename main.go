package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shopspring/decimal"
	"github.com/vaulverin/trades-parser/input"
	"github.com/vaulverin/trades-parser/logparser"
	"log"
	"os"
)

type SwapsIntersection struct {
	DexesWithBuys  map[string]struct{}
	DexesWithSells map[string]struct{}
	Buys           []SwapView
	Sells          []SwapView
}

type SwapView struct {
	DexName   string
	Price     decimal.Decimal
	Size      decimal.Decimal
	Side      string
	Timestamp uint64
}

func main() {
	input.Parse()

	client, err := ethclient.Dial(input.Config.Rpc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loading logs...")
	intersections := map[uint64]*SwapsIntersection{}
	for _, dex := range input.Config.Dexes {
		parser, err := logparser.New(client, dex.PairAddr, dex.Decimals)
		if err != nil {
			log.Println(err)
			continue
		}
		swaps, err := parser.ParseSwaps(context.Background(), input.Start, input.End)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, swap := range swaps {
			value, ok := intersections[swap.Block]
			if !ok {
				value = &SwapsIntersection{
					DexesWithBuys:  map[string]struct{}{},
					DexesWithSells: map[string]struct{}{},
					Buys:           []SwapView{},
					Sells:          []SwapView{},
				}
				intersections[swap.Block] = value
			}
			view := SwapView{
				DexName:   dex.Name,
				Price:     swap.Price,
				Size:      swap.Amount,
				Timestamp: swap.Timestamp,
			}
			if swap.Side == logparser.Buy {
				value.DexesWithBuys[dex.Name] = struct{}{}
				view.Side = "Buy"
				value.Buys = append(value.Buys, view)
			} else {
				value.DexesWithSells[dex.Name] = struct{}{}
				view.Side = "Sell"
				value.Sells = append(value.Sells, view)
			}
		}
	}
	fmt.Println("Done")

	lenDex := len(input.Config.Dexes)
	for block, intersection := range intersections {
		var views []SwapView
		if len(intersection.DexesWithBuys) == lenDex {
			views = append(views, intersection.Buys...)
		}
		if len(intersection.DexesWithSells) == lenDex {
			views = append(views, intersection.Sells...)
		}
		if len(views) > 0 {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetTitle(fmt.Sprintf("Block #%d", block))
			t.AppendHeader(table.Row{"Dex Name", "Price", "Size", "Side", "Timestamp"})
			for _, view := range views {
				t.AppendRows([]table.Row{
					{view.DexName, view.Price.StringFixed(2), view.Size.StringFixed(6), view.Side, view.Timestamp},
				})
			}
			t.Render()
		}
	}
}
