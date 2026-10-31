package gdrive

import (
	"context"
	"errors"
	"io"

	"kiroku/backend/internal/backup/provider"
)

type GDriveProvider struct {
	clientID     string
	clientSecret string
	refreshToken string
	folderID     string
}

func New(clientID, clientSecret, refreshToken, folderID string) *GDriveProvider {
	return &GDriveProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		folderID:     folderID,
	}
}

func (g *GDriveProvider) Name() string {
	return "gdrive"
}

func (g *GDriveProvider) Upload(ctx context.Context, key string, r io.Reader, meta provider.Metadata) error {
	// Implementation placeholder for Google Drive API upload
	return errors.New("gdrive upload not yet fully configured")
}

func (g *GDriveProvider) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, errors.New("gdrive download not implemented")
}

func (g *GDriveProvider) List(ctx context.Context) ([]provider.Metadata, error) {
	return nil, errors.New("gdrive list not implemented")
}

func (g *GDriveProvider) Delete(ctx context.Context, key string) error {
	return errors.New("gdrive delete not implemented")
}

func (g *GDriveProvider) HealthCheck(ctx context.Context) error {
	if g.refreshToken == "" {
		return errors.New("gdrive refresh token missing")
	}
	return nil
}
