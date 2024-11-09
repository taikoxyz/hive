package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"taiko/bindings/devnettierrouter"
	"taiko/bindings/taikol1"
	"taiko/common/utils"
	"taiko/params"
)

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

	tierRouter, err := devnettierrouter.NewDevnetTierRouter(tierRouterAddress, client)
	if err != nil {
		return err
	}

	providerAddress, err := tierRouter.GetProvider(&bind.CallOpts{Context: ctx}, common.Big0)
	if err != nil {
		return err
	}

	tierProvider, err := devnettierrouter.NewDevnetTierRouter(providerAddress, client)
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

	a.tiers = make(map[uint16]*devnettierrouter.ITierProviderTier)
	for _, id := range ids {
		tier, err := tierProvider.GetTier(&bind.CallOpts{Context: ctx}, id)
		if err != nil {
			return err
		}
		a.tiers[id] = &tier
	}

	return nil
}
