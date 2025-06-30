package google_drive_storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveStorage struct {
	StorageID    uuid.UUID `json:"storageId"    gorm:"primaryKey;type:uuid;column:storage_id"`
	ClientID     string    `json:"clientId"     gorm:"not null;type:text;column:client_id"`
	ClientSecret string    `json:"clientSecret" gorm:"not null;type:text;column:client_secret"`
	TokenJSON    string    `json:"tokenJson"    gorm:"not null;type:text;column:token_json"`
}

func (s *GoogleDriveStorage) TableName() string {
	return "google_drive_storages"
}

func (s *GoogleDriveStorage) SaveFile(
	logger *slog.Logger,
	fileID uuid.UUID,
	file io.Reader,
) error {
	driveService, err := s.getDriveService()
	if err != nil {
		return err
	}

	ctx := context.Background()
	filename := fileID.String()

	// Delete any previous copy so we keep at most one object per logical file.
	_ = s.deleteByName(ctx, driveService, filename) // ignore "not found"

	fileMeta := &drive.File{Name: filename}

	_, err = driveService.Files.Create(fileMeta).Media(file).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to upload file to Google Drive: %w", err)
	}

	logger.Info("file uploaded to Google Drive", "name", filename)

	return nil
}

func (s *GoogleDriveStorage) GetFile(fileID uuid.UUID) (io.ReadCloser, error) {
	driveService, err := s.getDriveService()
	if err != nil {
		return nil, err
	}

	fileIDGoogle, err := s.lookupFileID(driveService, fileID.String())
	if err != nil {
		return nil, err
	}

	resp, err := driveService.Files.Get(fileIDGoogle).Download()
	if err != nil {
		return nil, fmt.Errorf("failed to download file from Google Drive: %w", err)
	}

	return resp.Body, nil
}

func (s *GoogleDriveStorage) DeleteFile(fileID uuid.UUID) error {
	driveService, err := s.getDriveService()
	if err != nil {
		return err
	}

	ctx := context.Background()
	return s.deleteByName(ctx, driveService, fileID.String())
}

func (s *GoogleDriveStorage) Validate() error {
	switch {
	case s.ClientID == "":
		return errors.New("client ID is required")
	case s.ClientSecret == "":
		return errors.New("client secret is required")
	case s.TokenJSON == "":
		return errors.New("token JSON is required")
	}

	return nil
}

func (s *GoogleDriveStorage) TestConnection() error {
	driveService, err := s.getDriveService()
	if err != nil {
		return err
	}

	ctx := context.Background()
	testFilename := "test-connection-" + uuid.New().String()
	testData := []byte("test")

	// Test write operation
	fileMeta := &drive.File{Name: testFilename}
	file, err := driveService.Files.Create(fileMeta).
		Media(strings.NewReader(string(testData))).
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to write test file to Google Drive: %w", err)
	}

	// Test read operation
	resp, err := driveService.Files.Get(file.Id).Download()
	if err != nil {
		// Clean up test file before returning error
		_ = driveService.Files.Delete(file.Id).Context(ctx).Do()
		return fmt.Errorf("failed to read test file from Google Drive: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	readData, err := io.ReadAll(resp.Body)
	if err != nil {
		// Clean up test file before returning error
		_ = driveService.Files.Delete(file.Id).Context(ctx).Do()
		return fmt.Errorf("failed to read test file data: %w", err)
	}

	// Clean up test file
	if err := driveService.Files.Delete(file.Id).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to clean up test file: %w", err)
	}

	// Verify data matches
	if string(readData) != string(testData) {
		return fmt.Errorf(
			"test file data mismatch: expected %q, got %q",
			string(testData),
			string(readData),
		)
	}

	return nil
}

func (s *GoogleDriveStorage) getDriveService() (*drive.Service, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(s.TokenJSON), &token); err != nil {
		return nil, fmt.Errorf("invalid token JSON: %w", err)
	}

	ctx := context.Background()

	cfg := &oauth2.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/drive.file"},
	}

	tokenSource := cfg.TokenSource(ctx, &token)

	// Try to get a fresh token
	_, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get fresh token: %w", err)
	}

	driveService, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive client: %w", err)
	}

	return driveService, nil
}

func (s *GoogleDriveStorage) lookupFileID(
	driveService *drive.Service,
	name string,
) (string, error) {
	query := fmt.Sprintf("name = '%s' and trashed = false", escapeForQuery(name))

	results, err := driveService.Files.List().
		Q(query).
		Fields("files(id)").
		PageSize(1).
		Do()
	if err != nil {
		return "", fmt.Errorf("file lookup failed: %w", err)
	}

	if len(results.Files) == 0 {
		return "", fmt.Errorf("file %q not found in Google Drive", name)
	}

	return results.Files[0].Id, nil
}

func (s *GoogleDriveStorage) deleteByName(
	ctx context.Context,
	driveService *drive.Service,
	name string,
) error {
	query := fmt.Sprintf("name = '%s' and trashed = false", escapeForQuery(name))

	err := driveService.
		Files.
		List().
		Q(query).
		Fields("files(id)").
		Pages(ctx, func(p *drive.FileList) error {
			for _, file := range p.Files {
				if err := driveService.Files.Delete(file.Id).Context(ctx).Do(); err != nil {
					return err
				}
			}

			return nil
		})

	if err != nil {
		return fmt.Errorf("failed to delete %q: %w", name, err)
	}

	return nil
}

func escapeForQuery(s string) string {
	return strings.ReplaceAll(s, `'`, `\'`)
}
