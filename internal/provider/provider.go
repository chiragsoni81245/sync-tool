package provider

import "sync-tool/internal/db"

type Provider interface {
    PullSync(target db.SyncTarget) error
    PushSync(target db.SyncTarget) error
}
