package input

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Cfg struct {
	Rpc   string   `yaml:"rpc"`
	Dexes []DexCfg `yaml:"dexes"`
}

type DexCfg struct {
	Name     string         `yaml:"name"`
	PairAddr common.Address `yaml:"pair"`
	Decimals [2]uint8       `yaml:"decimals"`
}

var (
	Config Cfg
	Start  uint64
	End    *uint64
)

func Parse() {
	var (
		configPath string
		start      int64
		end        int64
	)
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Int64Var(&start, "start", -1, "start block to parse logs")
	flag.Int64Var(&end, "end", -1, "end block to parse logs (optional)")
	flag.Parse()

	if len(configPath) == 0 || start == -1 {
		fmt.Println("Usage: trades-parser")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if start < 0 || (end != -1 && end < 0) {
		fmt.Println("'start' and 'end' parameters must be positive")
		os.Exit(1)
	}

	Start = uint64(start)

	if end >= 0 {
		uend := uint64(end)
		End = &uend
	}

	if err := cleanenv.ReadConfig(configPath, &Config); err != nil {
		log.Fatal(err)
	}
}