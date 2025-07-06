// internal/gdrive/gdrive.go
package gdrive

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/logger"
	"sync-tool/internal/utils"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Gdrive struct {
}

func New() *Gdrive{
    return &Gdrive{}
}

func (p *Gdrive) PullSync(target db.SyncTarget) error {
	timeNow := time.Now()
	target.LastSyncedAt = &timeNow

    ctx := context.Background()

	// Load service account credentials
	credBytes, err := os.ReadFile(config.App.GoogleDriveCredentialsFilepath)
	if err != nil {
		return fmt.Errorf("Failed to read credentials.json: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credBytes, drive.DriveReadonlyScope)
	if err != nil {
		return fmt.Errorf("Failed to parse credentials: %v", err)
	}

	client := config.Client(ctx)

	// Create Drive service
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Failed to create Drive service: %v", err)
	}

	folderID, _ := utils.ExtractGoogleDriveFolderID(target.RemoteRef)
	destPath := target.LocalPath

	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create destination folder: %v", err)
	}

	err = downloadFolder(srv, folderID, destPath)
	if err != nil {
		logger.Log.Errorf("Download failed: %v", err)
        return nil
	}

	target.LastSyncStatus = db.StatusSuccess
	target.StatusMessage = "Sync successful"
	db.Save(&target)
	return nil
}

func (p *Gdrive) PushSync(target db.SyncTarget) error {
    logger.Log.Infof("Push operation for Gdrive provider is not built yet!")
    return nil
}

func downloadFile(srv *drive.Service, file *drive.File, destFolder string) error {
	if file.Md5Checksum == "" {
		logger.Log.Infof("Skipping file without md5Checksum: %s (%s)\n", file.Name, file.Id)
		return nil
	}

	outPath := filepath.Join(destFolder, file.Name)

	// If file exists, compare hash
	if _, err := os.Stat(outPath); err == nil {
		localHash, err := utils.GetLocalMD5(outPath)
		if err == nil && localHash == file.Md5Checksum {
			logger.Log.Infof("Skipping unchanged file: %s\n", file.Name)
			return nil
		}
        if err != nil {
            return err
        }
	}

	// Otherwise, download
	logger.Log.Infof("Downloading: %s (%s)\n", file.Name, file.Id)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	resp, err := srv.Files.Get(file.Id).Download()
	if err != nil {
        return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

// downloadFolder recursively downloads all non-Google-Docs files in a Drive folder
func downloadFolder(srv *drive.Service, folderID, destPath string) error {
	q := fmt.Sprintf("'%s' in parents and trashed = false", folderID)

	fileList, err := srv.Files.List().
		Q(q).
		Fields("files(id, name, mimeType, md5Checksum)").
		Do()
	if err != nil {
		return err
	}

	for _, f := range fileList.Files {
		if f.MimeType == "application/vnd.google-apps.folder" {
			subDir := filepath.Join(destPath, f.Name)
			os.MkdirAll(subDir, os.ModePerm)
			err := downloadFolder(srv, f.Id, subDir)
			if err != nil {
                return err
			}
		} else if f.MimeType[:20] == "application/vnd.google" {
			fmt.Printf("Skipping Google Doc: %s\n", f.Name)
		} else {
			err := downloadFile(srv, f, destPath)
			if err != nil {
                return err
			}
		}
	}

	return nil
}
