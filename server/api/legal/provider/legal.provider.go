package provider

import (
	"context"
	"fmt"
	"modular_monolith/server/api/legal/models"
	"modular_monolith/server/config/db"
	"time"

	"github.com/jackc/pgx/v4"
)

func InsertFile(file *models.UploadFileRequest, metadata *models.File) (string, error) {
	ctx := context.Background()

	checkQuery := `SELECT id FROM file WHERE id = $1`
	var existingID string
	err := db.DB.QueryRow(ctx, checkQuery, metadata.Id).Scan(&existingID)
	if err == nil {
		return "", fmt.Errorf("file with ID %s already exists", metadata.Id)
	}

	if err != nil && err.Error() != "no rows in result set" {
		return "", fmt.Errorf("error checking existing file: %w", err)
	}

	insertQuery := `
		INSERT INTO file (
			id, name, download, size, type, path, folder_id,
			title, nomor, tahun, tgl_penetapan, tgl_berlaku, uploaded_by, is_trash, index_text, "updatedAt"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
		        $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id
	`

	var insertedID string
	err = db.DB.QueryRow(
		ctx,
		insertQuery,
		metadata.Id,
		metadata.FileName,
		metadata.Download,
		metadata.Size,
		metadata.Type,
		metadata.Path,
		file.FolderID,
		file.Title,
		file.Nomor,
		file.Tahun,
		file.TglPenetapan,
		file.TglBerlaku,
		metadata.UploadedBy,
		metadata.IsTrash,
		metadata.IndexText,
		time.Now(),
	).Scan(&insertedID)
	if err != nil {
		return "", fmt.Errorf("failed to insert file: %w", err)
	}

	return insertedID, nil
}

func InsertRelation(ctx context.Context, tx pgx.Tx, fileID string, relationType string, relations []models.RelatedFile) error {
	query := `
		INSERT INTO relation (impacter_file_id, impacted_file_id, type, reason)
		VALUES ($1, $2, $3, $4)
	`

	for _, rel := range relations {
		if rel.TargetFileID == "" {
			return fmt.Errorf("target_file_id is required for %s relation", relationType)
		}

		_, err := tx.Exec(ctx, query, fileID, rel.TargetFileID, relationType, rel.Reason)
		if err != nil {
			return fmt.Errorf("failed to insert %s relation: %w", relationType, err)
		}
	}
	return nil
}

func CheckFolderExists(folderID int) error {
	ctx := context.Background()

	fmt.Println(folderID)
	query := `SELECT 1 FROM folder WHERE id = $1`
	var exists int
	err := db.DB.QueryRow(ctx, query, folderID).Scan(&exists)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("folder with ID %s does not exist", folderID)
		}
		return fmt.Errorf("error checking folder existence: %w", err)
	}

	return nil
}
