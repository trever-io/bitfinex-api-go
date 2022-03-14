package currency

import (
	"strings"
)

type Conf struct {
	Currency  string
	Label     string
	Symbol    string
	Pairs     []string
	Pools     []string
	Explorers ExplorerConf
	Unit      string
	Methods   []string
	Fees      Fees
	Infos     Infos
}

type ExplorerConf struct {
	BaseUri        string
	AddressUri     string
	TransactionUri string
}

type Fees struct {
	WithdrawFee float64
	DepositFee  float64
}

type Infos struct {
	MinOrderSize  string
	MaxOrderSize  string
	InitialMargin float64
	MinMargin     float64
}

type ConfigMapping string

const (
	LabelMap    ConfigMapping = "pub:map:currency:label"
	SymbolMap   ConfigMapping = "pub:map:currency:sym"
	UnitMap     ConfigMapping = "pub:map:currency:unit"
	ExplorerMap ConfigMapping = "pub:map:currency:explorer"
	ExchangeMap ConfigMapping = "pub:list:pair:exchange"
	MethodMap   ConfigMapping = "pub:map:tx:method"
	FeesMap     ConfigMapping = "pub:map:currency:tx:fee"
	InfoMap     ConfigMapping = "pub:info:pair"
)

type RawConf struct {
	Mapping string
	Data    interface{}
}

func parseLabelMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Label = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := Conf{}
			cfg.Label = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseMethodMap(config map[string]Conf, raw []interface{}) {
	for _, rawMethod := range raw {
		data := rawMethod.([]interface{})
		curs := data[1].([]interface{})
		for _, curInt := range curs {
			cur := curInt.(string)
			if val, ok := config[cur]; ok {
				if val.Methods == nil {
					val.Methods = make([]string, 0)
				}

				val.Methods = append(val.Methods, data[0].(string))
				config[cur] = val
			} else {
				cfg := Conf{}
				cfg.Methods = []string{data[0].(string)}
				cfg.Currency = cur
				config[cur] = cfg
			}
		}
	}
}

func parseFeesMap(config map[string]Conf, raw []interface{}) {
	for _, rawFees := range raw {
		data := rawFees.([]interface{})
		cur := data[0].(string)
		fees := data[1].([]interface{})
		depositFee := fees[0].(float64)
		withdrawFee := fees[1].(float64)
		if val, ok := config[cur]; ok {
			val.Fees = Fees{
				WithdrawFee: withdrawFee,
				DepositFee:  depositFee,
			}
			config[cur] = val
		} else {
			cfg := Conf{}
			cfg.Fees = Fees{
				WithdrawFee: withdrawFee,
				DepositFee:  depositFee,
			}
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseSymbMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Symbol = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := Conf{}
			cfg.Symbol = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseUnitMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Unit = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := Conf{}
			cfg.Unit = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseExplorerMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		explorers := data[1].([]interface{})
		var cfg Conf
		if val, ok := config[cur]; ok {
			cfg = val
		} else {
			// create new empty config instance
			cc := Conf{}
			cc.Currency = cur
			cfg = cc
		}
		ec := ExplorerConf{
			explorers[0].(string),
			explorers[1].(string),
			explorers[2].(string),
		}
		cfg.Explorers = ec
		config[cur] = cfg
	}
}

func parseExchangeMap(config map[string]Conf, raw []interface{}) {
	for _, rs := range raw {
		symbol := rs.(string)
		var base, quote string

		if len(symbol) > 6 {
			base = strings.Split(symbol, ":")[0]
			quote = strings.Split(symbol, ":")[1]
		} else {
			base = symbol[3:]
			quote = symbol[:3]
		}

		// append if base exists in configs
		if val, ok := config[base]; ok {
			val.Pairs = append(val.Pairs, symbol)
			config[base] = val
		}

		// append if quote exists in configs
		if val, ok := config[quote]; ok {
			val.Pairs = append(val.Pairs, symbol)
			config[quote] = val
		}
	}
}

func parseInfoMap(config map[string]Conf, raw []interface{}) {
	for _, rawInfo := range raw {
		data := rawInfo.([]interface{})
		symbol := data[0].(string)
		infos := data[1].([]interface{})

		minOrderSize := infos[3].(string)
		maxOrderSize := infos[4].(string)
		var initialMargin float64
		var minMargin float64

		if infos[8] != nil {
			initialMargin = infos[8].(float64)
		}
		if infos[9] != nil {
			minMargin = infos[9].(float64)
		}

		if val, ok := config[symbol]; ok {
			val.Infos = Infos{
				MinOrderSize:  minOrderSize,
				MaxOrderSize:  maxOrderSize,
				InitialMargin: initialMargin,
				MinMargin:     minMargin,
			}
		} else {
			cfg := Conf{}
			cfg.Infos = Infos{
				MinOrderSize:  minOrderSize,
				MaxOrderSize:  maxOrderSize,
				InitialMargin: initialMargin,
				MinMargin:     minMargin,
			}
			cfg.Symbol = symbol
			config[symbol] = cfg
		}
	}
}

func FromRaw(raw []RawConf) ([]Conf, error) {
	configMap := make(map[string]Conf)
	for _, r := range raw {
		switch ConfigMapping(r.Mapping) {
		case LabelMap:
			data := r.Data.([]interface{})
			parseLabelMap(configMap, data)
		case SymbolMap:
			data := r.Data.([]interface{})
			parseSymbMap(configMap, data)
		case UnitMap:
			data := r.Data.([]interface{})
			parseUnitMap(configMap, data)
		case ExplorerMap:
			data := r.Data.([]interface{})
			parseExplorerMap(configMap, data)
		case ExchangeMap:
			data := r.Data.([]interface{})
			parseExchangeMap(configMap, data)
		case MethodMap:
			data := r.Data.([]interface{})
			parseMethodMap(configMap, data)
		case FeesMap:
			data := r.Data.([]interface{})
			parseFeesMap(configMap, data)
		case InfoMap:
			data := r.Data.([]interface{})
			parseInfoMap(configMap, data)
		}
	}

	// convert map to array
	configs := make([]Conf, 0)
	for _, v := range configMap {
		configs = append(configs, v)
	}

	return configs, nil
}
