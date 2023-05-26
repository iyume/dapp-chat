package config

import (
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/server"
	"gopkg.in/ini.v1"
)

// --mine --datadir nodes/signer1 --networkid 12345 --port 30304 --authrpc.port 8552 \
// --unlock 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
// --miner.etherbase 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
// --miner.gasprice 1 --syncmode full \
// --bootnodes 'enode://9b91d8e38cf58aaf6eca6d05627b8d266baabf445e5980194932ed1827c7ff66527aeb4dfe5e62d1e4b4a51f9e5e53c3222966641c0eb9c4dbd1bea946d44d80@127.0.0.1:30303' \
// --nat none --netrestrict '127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16'

// nat: p2p.nat.Parse
// netrestrict: p2p.netutil.ParseNetlist

type Config struct {
	DataDir string
	Http    server.HTTPConfig
	Backend api.BackendConfig
}

var DefaultConfig = Config{
	DataDir: "chatdata",
	Http:    server.DefaultHTTPConfig,
	Backend: api.DefaultBackendConfig,
}

func LoadINIConfig(file string) Config {
	inifile, err := ini.InsensitiveLoad(file)
	if err != nil {
		panic(err)
	}
	cfg := DefaultConfig
	cfghttp := cfg.Http
	cfgbackend := cfg.Backend

	// Main config (no section)
	inimain := inifile.Section("")
	cfg.DataDir = inimain.Key("data_dir").MustString(cfg.DataDir)
	// Http config
	if inihttp := inifile.Section("http"); inifile.HasSection("http") {
		cfghttp.Address = inihttp.Key("address").MustString(cfghttp.Address)
		cfghttp.Token = inihttp.Key("token").MustString(cfghttp.Token)
	}
	// Backend config
	if inibackend := inifile.Section("backend"); inifile.HasSection("backend") {
		cfgbackend.Key = MustHexToECDSA(inibackend.Key("private_key").String())
		bootnodes_slice := inibackend.Key("bootnodes").Strings(",")
		bootnodes := []*enode.Node{}
		for _, val := range bootnodes_slice {
			bootnodes = append(bootnodes, enode.MustParse(val))
		}
		cfgbackend.BootstrapNodes = bootnodes
		cfgbackend.Address = inibackend.Key("address").MustString(cfgbackend.Address)
	}

	cfg.Http = cfghttp
	cfg.Backend = cfgbackend
	return cfg
}
