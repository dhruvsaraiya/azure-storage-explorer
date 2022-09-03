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
	return &BlobService{
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (b *BlobService) Init() error {
	b.logger.Info("Initializing blob service...")
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		b.logger.Error("Invalid credentials with error: " + err.Error())
		return err
	}
	serviceClient, err := azblob.NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		b.logger.Error("Invalid service client with error: " + err.Error())
		return err
	}
	b.client = serviceClient
	return nil
}
