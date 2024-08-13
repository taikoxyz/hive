package clients

type ValidatorClient struct {
	*HiveManagedClient
}

type ValidatorClients []*ValidatorClient

// Return subset of clients that are currently running
func (all ValidatorClients) Running() ValidatorClients {
	res := make(ValidatorClients, 0)
	for _, vc := range all {
		if vc.IsRunning() {
			res = append(res, vc)
		}
	}
	return res
}
