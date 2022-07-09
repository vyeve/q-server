package repository

import (
	"context"

	"github.com/vyeve/q-server/models"
)

type Repository interface {
	UploadTransfers(ctx context.Context, r *models.Receipt) error
}
