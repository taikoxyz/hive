package clients

import (
	"context"
	"testing"
)

func TestVerifyL1Origin(t *testing.T) {
	VerifyL1Origin(context.Background(), 15, rpccli.L1, rpccli.L2)
}
