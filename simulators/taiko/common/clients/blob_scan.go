package clients

import (
	"fmt"
	"time"
)

type BlobScanClient struct {
	Mysql       *HiveManagedClient
	BlobApi     *HiveManagedClient
	BlobIndexer *HiveManagedClient
}

func (b *BlobScanClient) BlobAPIURL() string {
	return fmt.Sprintf("http://%s:%d", b.BlobApi.NetworkIP(), BlobscanAPIPort)
}

func (b *BlobScanClient) PostgresURL() string {
	return fmt.Sprintf("postgresql://blobscan:s3cr3t@%s:5432/blobscan_dev?schema=public", b.Mysql.NetworkIP())
}

func (b *BlobScanClient) Start() error {
	if err := b.Mysql.Start(); err != nil {
		return err
	}

	time.Sleep(time.Second)

	if err := b.BlobApi.Start(); err != nil {
		return err
	}

	time.Sleep(time.Second * 3)

	if err := b.BlobIndexer.Start(); err != nil {
		return err
	}
	return nil
}

func (b *BlobScanClient) Shutdown() error {
	b.Mysql.Shutdown()
	b.BlobApi.Shutdown()
	b.BlobIndexer.Shutdown()

	return nil
}

func (b *BlobScanClient) IsRunning() bool {
	return b.Mysql.IsRunning() && b.BlobApi.IsRunning() && b.BlobIndexer.IsRunning()
}
