package blob

import (
	"context"
	"fmt"
	"os"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobContainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	azblobService "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"go.uber.org/zap"
)

type BlobService struct {
	ctx    context.Context
	logger *zap.Logger
	client *azblobService.Client
}

func NewBlobService(ctx context.Context, logger *zap.Logger) (*BlobService, error) {
	logger.Info("Initializing blob service...")
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	cred, err := azblobService.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		logger.Error("Invalid credentials with error: " + err.Error())
		return nil, err
	}
	client, err := azblobService.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		logger.Error("Invalid client with error: " + err.Error())
		return nil, err
	}

	return &BlobService{
		ctx:    ctx,
		logger: logger,
		client: client,
	}, nil
}

func (b *BlobService) ListContainers(ctx context.Context) ([]string, error) {
	containers := make([]string, 0)
	containersPager := b.client.NewListContainersPager(&azblob.ListContainersOptions{})
	for containersPager.More() {
		resp, err := containersPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, container := range resp.ContainerItems {
			containers = append(containers, *(container.Name))
		}
	}

	return containers, nil
}

func (b *BlobService) ListBlobsHierarchy(ctx context.Context, containerName string, prefix string) ([]string, []string, error) {
	containerClient := b.client.NewContainerClient(containerName)

	blobs := make([]string, 0)
	prefixes := make([]string, 0)
	blobsPager := containerClient.NewListBlobsHierarchyPager("/", &azblobContainer.ListBlobsHierarchyOptions{})

	fmt.Printf("Listing blobs for container %s", containerName)
	for blobsPager.More() {
		resp, err := blobsPager.NextPage(ctx)
		if err != nil {
			return nil, nil, err
		}
		for _, blob := range resp.Segment.BlobItems {
			blobs = append(blobs, *(blob.Name))
		}
		for _, blobPrefix := range resp.Segment.BlobPrefixes {
			prefixes = append(prefixes, *(blobPrefix.Name))
		}
	}

	return blobs, prefixes, nil
}

func (b *BlobService) ListBlobs(ctx context.Context, containerName string, prefix string) ([]string, error) {
	containerClient := b.client.NewContainerClient(containerName)

	blobs := make([]string, 0)
	blobsPager := containerClient.NewListBlobsHierarchyPager("/", &azblobContainer.ListBlobsHierarchyOptions{
		Include: azblobContainer.ListBlobsInclude{Metadata: true},
		Prefix:  &prefix,
	})

	for blobsPager.More() {
		resp, err := blobsPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, blob := range resp.Segment.BlobItems {
			blobs = append(blobs, *(blob.Name))
		}
	}

	return blobs, nil
}

func (b *BlobService) ListPrefixes(ctx context.Context, containerName string, prefix string) ([]string, error) {
	containerClient := b.client.NewContainerClient(containerName)

	prefixes := make([]string, 0)
	blobsPager := containerClient.NewListBlobsHierarchyPager("/", &azblobContainer.ListBlobsHierarchyOptions{Prefix: &prefix})

	for blobsPager.More() {
		resp, err := blobsPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, blobPrefix := range resp.Segment.BlobPrefixes {
			prefixes = append(prefixes, *(blobPrefix.Name))
		}
	}

	return prefixes, nil
}

// FLAT
// func (b *BlobService) ListBlobs(ctx context.Context, container, prefix string) ([]string, []string, error) {
// 	blobs := make([]string, 0)
// 	prefixes := make([]string, 0)
// 	blobsPager := b.client.ListBlobsPager(container, &azblob.ListBlobsFlatOptions{
// 		Prefix: &prefix,
// 	})

// 	for blobsPager.More() {
// 		resp, err := blobsPager.NextPage(ctx)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		for _, blob := range resp.Segment.BlobItems {
// 			blobs = append(blobs, *(blob.Name))
// 		}
// 		// for _, prefix := range resp.Segment.BlobPrefixes {
// 		// 	prefixes = append(prefixes, *(prefix.Name))
// 		// }
// 	}

// 	// if blobsPager.Err() != nil {
// 	// 	b.logger.Error("Error listing blobs: " + blobsPager.Err().Error())
// 	// 	return blobs, blobsPager.Err()
// 	// }
// 	return blobs, prefixes, nil
// }
