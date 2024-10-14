# hive - Ethereum end-to-end test harness

Hive is a system for running integration tests against Ethereum clients.

Ethereum Foundation maintains two public Hive instances to check for consensus, p2p and
blockchain compatibility:

- eth1 consensus, graphql and p2p tests are on <https://hivetests.ethdevops.io>
- Engine API integration and rpc tests are on <https://hivetests2.ethdevops.io>

**To read more about hive, please check [the documentation][doc].**

### Run hiveview

> Open `http://localhost:8080` on explorer, and then we can see all the test results.

```shell
./build/bin/hiveview --serve --logdir ./workspace/logs
```

### Update contract transactions

* Get the taiko-mono repo and turn to const_contracts branch

```
git clone git@github.com:taikoxyz/taiko-mono.git
git checkout const_contracts
```

* Update contract txs list(all the txs in simulators/params/l1contract_txs.txt file)

```shell
git clone git@github.com:taikoxyz/hive.git
cd simulators/taiko
TAIKO_MONO_DIR=taiko-mono_path ./scripts/deploy_l1_contract.sh
```

### Build binaries

```shell
git clone git@github.com:taikoxyz/hive.git
make hive
make hiveview
```

### Run hive test cases

* Run l2-full-sync test:

```shell
./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,taiko/taiko-geth,taiko/driver,taiko/taiko-geth,taiko/driver --sim taiko --sim.limit "taiko-genesis/l2-full-sync"
```

* Run l2-snap-sync test:

```shell
./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,taiko/taiko-geth,taiko/driver,taiko/taiko-geth,taiko/driver --sim taiko --sim.limit "taiko-genesis/l2-snap-sync"
```

* Run taiko-reorg test:

```shell
./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "taiko-reorg/taiko-reorg"
```

* Run blob-l1-beacon test:

```shell
./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "taiko-blob/blob-l1-beacon"
```

* Run blob-server test:

```shell
./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,storage/redis,storage/postgres,blobscan/blobscan-api,blobscan/blobscan-indexer --sim taiko --sim.limit "taiko-blob/blob-server"
```

### Hive framework

* `./clients`: contains multi-docker images these are used for testing.
* `./simulators/taiko`
    * `binding`: taiko contract SDK.
    * `scripts`: some shell scripts used to get txs and update envs and update SDK files.
    * `docker`: docker: docker-compose folder used to deploy taiko contracts and get txs.
    * `common/clients`: They are docker images that relate to clients.
    * `common/config`: Most of them are related to the configuration of the beacon geth node.
    * `common/spoofing`: Used to genesis l1 execution node.
    * `common/testnet`
    * `suites`: tests instances, new test case plz add in this folder.
        * `base`:
            * taiko-base/l2-full-sync
            * taiko-base/l2-snap-sync
        * `blob`:
            * taiko-blob/blob-l1-beacon
            * taiko-blob/blob-server
        * `reorg`:
            * taiko-reorg/taiko-reorg

### add new test case

* Create a new folder named `client` in `simulators/taiko/suites`

* Add `execution.go`

```go
type ClientTestSpec struct {
    suite_base.BaseTestSpec
}

func (r ClientTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
    panic("Plz add test content in this function.")
}

```

* Add `tests.go`
```go
var testSuite = hivesim.Suite{
    Name:        "taiko-blob",
    DisplayName: "driver blob client test",
    Location:    "suites/blob",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
    Tests = append(Tests,
    BlobTestSpec{
        TestL1Beacon: true,
        BaseTestSpec: suite_base.BaseTestSpec{
            Name: "blob-l1-beacon",
        },
    },
    BlobTestSpec{
        TestBlobServer: true,
        BaseTestSpec: suite_base.BaseTestSpec{
            Name: "blob-server",
        },
    })
}

func Suite(clients clients.ClientGroups) hivesim.Suite {
    // Load params.yml
    beaconConfig, err := params.UnmarshalConfig(taparams.ConfigContent, nil)
    if err != nil {
        panic(err)
    }
    
    var genesis core.Genesis
    // Load genesis.json
    if err = json.Unmarshal(taparams.GenesisContent, &genesis); err != nil {
        panic(err)
    }
	
    suites.SuiteHydrate(&testSuite, clients, Tests, &execution_config.GenesisState{
    BeaconConfig:     beaconConfig,
    Genesis:          &genesis,
    })
	return testSuite
}
```