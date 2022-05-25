package dnr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const BlockBuffer = 100
const LogDarknodeRegistered = "0x7c56cb7f63b6922d24414bf7c2b2c40c7ea1ea637c3f400efa766a85ecf2f093"
const LogDarknodeDeregistered = "0xf73268ea792d9dbf3e21a95ec9711f0b535c5f6c99f6b4f54f6766838086b842"
const LogNewEpoch = "0xaf2fc4796f2932ce294c3684deffe5098d3ef65dc2dd64efa80ef94eed88b01e"

type DNR struct {
	config         *params.CliqueConfig
	LastEpochBlock uint64                  `json:"epoch"`
	Validators     map[common.Address]bool `json:"validators"`
	synced         bool
	syncLock       sync.RWMutex
}

func NewDNR(config *params.CliqueConfig, db ethdb.Database) *DNR {
	validators := map[common.Address]bool{}
	for _, validator := range config.InitialValidators {
		validators[validator] = true
	}
	defaultDNR := &DNR{
		config:         config,
		LastEpochBlock: config.EpochBlock,
		Validators:     validators,
	}
	dnr, err := GetLatestDNR(db)
	if err != nil {
		dnr.config = config
		log.Info("loaded latest dnr from disk", "last_epoch", dnr.LastEpochBlock)
		return dnr
	}
	log.Warn("failed to load latest dnr from disk, creating new drn entry with genesis data", "error", err)
	if err = dnr.store(db); err != nil {
		log.Warn("failed to store default epoch in db", "epoch_number", config.EpochBlock)
	}
	return defaultDNR
}

func (d *DNR) Watch(ctx context.Context, db ethdb.Database) {
	d.syncLock.Lock()
	d.synced = false
	d.syncLock.Unlock()
	client, err := ethclient.DialContext(ctx, d.config.API)
	if err != nil {
		panic(err)
	}
	lastBlock := d.LastEpochBlock
	for {
		startBlockNumber := new(big.Int).SetUint64(lastBlock + 1)
		lastBlockNumber := new(big.Int).Add(startBlockNumber, big.NewInt(1000))

		latestBlock, err := client.BlockByNumber(ctx, nil)
		if err != nil {
			panic(err)
		}

		// a 100 block delay to handle reorgs
		latestBlockNum := new(big.Int).Sub(latestBlock.Number(), big.NewInt(100))

		if latestBlockNum.Cmp(startBlockNumber) < 0 {
			// sleep for some time before retrying as no new block were created
			time.Sleep(time.Minute)
			continue
		}

		if latestBlockNum.Cmp(lastBlockNumber) < 0 {
			lastBlockNumber = latestBlockNum
		}

		logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: startBlockNumber,
			ToBlock:   lastBlockNumber,
			Addresses: []common.Address{d.config.DNR},
			Topics: [][]common.Hash{
				{common.HexToHash(LogDarknodeRegistered)},
				{common.HexToHash(LogDarknodeDeregistered)},
				{common.HexToHash(LogNewEpoch)},
			},
		})
		if err != nil {
			panic(err)
		}

		for _, eventLog := range logs {
			switch eventLog.Topics[0].Hex() {
			case LogDarknodeRegistered:
				d.Validators[common.BytesToAddress(eventLog.Topics[2].Bytes())] = true
			case LogDarknodeDeregistered:
				if _, ok := d.Validators[common.BytesToAddress(eventLog.Topics[2].Bytes())]; ok {
					delete(d.Validators, common.BytesToAddress(eventLog.Topics[2].Bytes()))
				}
			case LogNewEpoch:
				d.LastEpochBlock = eventLog.BlockNumber
				if err = d.store(db); err != nil {
					log.Warn("failed to store epoch in db", "epoch_number", eventLog.BlockNumber)
				}
			}
		}

		// if caught up set synced to true
		if latestBlockNum.Cmp(lastBlockNumber) == 0 {
			d.syncLock.Lock()
			d.synced = true
			d.syncLock.Unlock()
		}

		lastBlock = lastBlockNumber.Uint64()
	}
}

func (d *DNR) store(db ethdb.Database) error {
	blob, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if err = db.Put([]byte(fmt.Sprintf("dnr-%v", d.LastEpochBlock)), blob); err != nil {
		return err
	}
	return db.Put([]byte("dnr-latest"), blob)
}

func (d *DNR) WaitSynced() {
	for {
		d.syncLock.RLock()
		if d.synced {
			return
		}
		d.syncLock.RUnlock()
		time.Sleep(10 * time.Second)
	}
}

func GetLatestDNR(db ethdb.Database) (*DNR, error) {
	blob, err := db.Get([]byte(fmt.Sprintf("dnr-latest")))
	if err != nil {
		return nil, err
	}
	dnr := new(DNR)
	if err := json.Unmarshal(blob, dnr); err != nil {
		return nil, err
	}
	return dnr, nil
}

func GetDNR(db ethdb.Database, epoch uint64) (*DNR, error) {
	blob, err := db.Get([]byte(fmt.Sprintf("dnr-%v", epoch)))
	if err != nil {
		return nil, err
	}
	dnr := new(DNR)
	if err := json.Unmarshal(blob, dnr); err != nil {
		return nil, err
	}
	return dnr, nil
}
