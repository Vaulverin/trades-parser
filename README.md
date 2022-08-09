# trades-parser
Parses intersections of trades from UniswapV2 compatible DEXes

## Usage
`./trades-parser -config ./config.yaml -start 15306222`

### Output example
```
+----------------------------------------------------+
| Block #15306578                                    |
+-----------+---------+----------+------+------------+
| DEX NAME  | PRICE   | SIZE     | SIDE |  TIMESTAMP |
+-----------+---------+----------+------+------------+
| UniswapV2 | 1781.52 | 3.283707 | Buy  | 1660028730 |
| Sushiswap | 1781.25 | 4.778318 | Buy  | 1660028730 |
| Sushiswap | 1782.38 | 0.073998 | Buy  | 1660028730 |
+-----------+---------+----------+------+------------+

```