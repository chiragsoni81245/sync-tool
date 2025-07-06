package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"errors"
	"regexp"
)

// compute MD5 hash of a local file
func GetLocalMD5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func ExtractGoogleDriveFolderID(url string) (string, error) {
	// Handle URLs like:
	// https://drive.google.com/drive/folders/1AbCDefGhIjKlMnOp
	// https://drive.google.com/drive/u/0/folders/1AbCDefGhIjKlMnOp
	// https://drive.google.com/open?id=1AbCDefGhIjKlMnOp

	patterns := []string{
		`/folders/([a-zA-Z0-9_-]+)`,        // matches /folders/{id}
		`[?&]id=([a-zA-Z0-9_-]+)`,          // matches ?id={id}
	}

	for _, pat := range patterns {
		re := regexp.MustCompile(pat)
		match := re.FindStringSubmatch(url)
		if len(match) > 1 {
			return match[1], nil
		}
	}

	return "", errors.New("invalid Google Drive folder URL")
}
