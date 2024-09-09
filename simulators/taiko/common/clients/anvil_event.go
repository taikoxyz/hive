package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"taiko/bindings/devnettierprovider"
	"taiko/bindings/taikol1"
	"taiko/common/utils"
	"taiko/params"
)

func (a *AnvilClient) HandleProposedEvent(start uint64) chan struct{} {
	a.Logf("%s: start watch proposed event, start number: %d", a.ClientType(), start)
	stopCh := make(chan struct{})

	a.proposedEvents = make([]*BlockProposed, 0)

	go func() {
		sink := make(chan *taikol1.TaikoL1BlockProposed, 10)
		sub, err := a.taikoL1.WatchBlockProposed(&bind.WatchOpts{Start: &start}, sink, nil, nil)
		if err != nil {
			a.Fatalf("%s, failed to watch BlockProposed event, err: %v", a.ClientType(), err)
		}
		defer sub.Unsubscribe()

		sink2 := make(chan *taikol1.TaikoL1BlockProposedV2, 10)
		sub2, err := a.taikoL1.WatchBlockProposedV2(&bind.WatchOpts{Start: &start}, sink2, nil)
		if err != nil {
			a.Fatalf("%s, failed to watch BlockProposed event, err: %v", a.ClientType(), err)
		}
		defer sub2.Unsubscribe()

		for {
			select {
			case <-stopCh:
				return
			case blockProposed := <-sink:
				a.proposedEvents = append(a.proposedEvents, &BlockProposed{
					BlockId: blockProposed.BlockId.Uint64(),
					MinTier: blockProposed.Meta.MinTier,
				})
			case blockProposed := <-sink2:
				a.proposedEvents = append(a.proposedEvents, &BlockProposed{
					BlockId: blockProposed.BlockId.Uint64(),
					MinTier: blockProposed.Meta.MinTier,
				})
			}
		}
	}()

	return stopCh
}

func (a *AnvilClient) depProposerEvent() {
	var pv1 *BlockProposed
	if len(a.proposedEvents) == 0 {
		return
	}
	pv1, a.proposedEvents = a.proposedEvents[0], a.proposedEvents[1:]
	a.IncreaseTime(a.tiers[pv1.MinTier].CooldownWindow.Uint64() * 60)
}

func (a *AnvilClient) FillTiers(ctx context.Context) error {
	client, err := ethclient.DialContext(ctx, a.WSURL())
	if err != nil {
		//return err
		a.Fatalf("%s, failed to dial ws client, err: %v", a.ClientType(), err)
	}

	if a.taikoL1 == nil {
		var (
			envs = params.EnvParams()
			err  error
		)
		a.taikoL1, err = taikol1.NewTaikoL1(common.HexToAddress(envs["TAIKO_L1"]), client)
		if err != nil {
			//return err
			a.Fatalf("%s, failed to new taikoL1, err: %v", a.ClientType(), err)
		}
	}

	tierRouterAddress, err := a.taikoL1.Resolve0(&bind.CallOpts{Context: ctx}, utils.StringToBytes32("tier_router"), false)

	tierRouter, err := devnettierprovider.NewDevnetTierProvider(tierRouterAddress, client)
	if err != nil {
		return err
	}

	providerAddress, err := tierRouter.GetProvider(&bind.CallOpts{Context: ctx}, common.Big0)
	if err != nil {
		return err
	}

	tierProvider, err := devnettierprovider.NewDevnetTierProvider(providerAddress, client)
	if err != nil {
		return err
	}

	ids, err := tierProvider.GetTierIds(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return errEmptyTiersList
	}

	a.tiers = make(map[uint16]*devnettierprovider.ITierProviderTier)
	for _, id := range ids {
		tier, err := tierProvider.GetTier(&bind.CallOpts{Context: ctx}, id)
		if err != nil {
			return err
		}
		a.tiers[id] = &tier
	}

	return nil
}
