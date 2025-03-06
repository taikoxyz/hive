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

#### Run l2-full-sync test:

* how it tests
    * Start one cluster at first;
    * Start the second and the third clusters and then wait for the l2 finalized header;
    * Verify the latest l2 finalized header in the three cluster;

```shell
make test_full_sync
```

#### Run pacaya preconf test:

* how it tests
    * Start one cluster at first;
    * Wait for pacaya hardfork block high;
    * Stop proposer;
    * Call preconf api and create preconf blocks;
    * Verify preconf status;
    * Create propose block and verify propose status;

```shell
make test_preconf_preconf
```

#### Run pacaya preconf reorg test:

* how it tests
    * Start one cluster at first
    * Wait for pacaya hardfork block high
    * Stop proposer
    * reorg propose blocks
        * create preconf blocks
        * set reorg point
        * create propose block
        * reorg l1 chain
        * wait for propose block recovered and verify the status same to origin.
    * reorg preconf blocks
        * create preconf blocks
        * create a reorg propose block then will reorg all the preconf blocks
        * verify reorg status

```shell
make test_preconf_reorg
```

#### Run reorg-reorg test:

* how it tests
    * Normally start one cluster and wait for a random reorg l2 block
    * Reorg l2chain to the latest finalized high
    * Waiting for the l2chain be reorged to the latest finalized high+1

```shell
make test_reorg_reorg
```

#### Run blob-l1-beacon test:

* how it tests
    * Start one cluster
    * Let proposer use blobs to propose blocks
    * Waiting taiko-geth tough the target number

```shell
make test_blob_beacon
```

#### Run blob-scan test:

* how it tests
    * Start one cluster
    * Let proposer use blobs to propose blocks
    * Waiting taiko-geth tough the target number

```shell
make test_blob_server
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