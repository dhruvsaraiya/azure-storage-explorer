package blob

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"
)

type BlobService struct {
	ctx    context.Context
	logger *zap.Logger
	client *azblob.ServiceClient
}

func NewBlobService(ctx context.Context, logger *zap.Logger) (*BlobService, error) {
	logger.Info("Initializing blob service...")
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		logger.Error("Invalid credentials with error: " + err.Error())
		return nil, err
	}
	serviceClient, err := azblob.NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		logger.Error("Invalid service client with error: " + err.Error())
		return nil, err
	}

	return &BlobService{
		ctx:    ctx,
		logger: logger,
		client: serviceClient,
	}, nil
}

func (b *BlobService) ListContainers() error {
	containerClient, err := b.client.NewContainerClient("github")
	if err != nil {
		b.logger.Error("Invalid container client with error: " + err.Error())
	}

	pager := containerClient.ListBlobsHierarchy("/", &azblob.ContainerListBlobsHierarchyOptions{
		Include: []azblob.ListBlobsIncludeItem{
			azblob.ListBlobsIncludeItemMetadata,
			azblob.ListBlobsIncludeItemTags,
		},
	})

	for pager.NextPage(context.TODO()) {
		resp := pager.PageResponse()
		for _, blob := range resp.ListBlobsHierarchySegmentResponse.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
		for _, blob := range resp.ListBlobsHierarchySegmentResponse.Segment.BlobPrefixes {
			fmt.Println(*blob.Name)
		}
	}

	if pager.Err() != nil {
		b.logger.Error("Error listing blobs: " + pager.Err().Error())
	}
	return nil
}
