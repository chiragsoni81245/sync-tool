package provider

import (
	"fmt"
	"sync-tool/internal/db"
	"sync-tool/internal/gdrive"
	"sync-tool/internal/github"
)

type Provider interface {
    Sync(target db.SyncTarget) error
}

func GetProviderViaName(providerName db.SyncProvider) (Provider, error) {
    switch providerName {
    case db.ProviderGitHub: return github.New(), nil
    case db.ProviderGDrive: return gdrive.New(), nil
    }
    return nil, fmt.Errorf("Invalid provider name")
}
