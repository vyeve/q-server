package repository

import (
	"context"

	"github.com/vyeve/q-server/models"
)

// Repository defines methods to make deal with persistent storage
type Repository interface {
	UploadTransfers(ctx context.Context, r *models.Receipt) error
}
