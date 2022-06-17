package accounthistory

import (
	table "github.com/zecrey-labs/zecrey-legend/common/model/account"
	"github.com/zecrey-labs/zecrey-legend/service/api/explorer/internal/svc"
)

type AccountHistory interface {
	GetAccountByAccountName(accountName string) (account *table.AccountHistory, err error)
	GetAccountByAccountIndex(accountIndex int64) (account *table.AccountHistory, err error)
	GetAccountByPk(pk string) (account *table.AccountHistory, err error)
	GetAccountAssetsByIndex(accountIndex int64) (accountAssets []*table.AccountHistory, err error)
}

func New(svcCtx *svc.ServiceContext) AccountHistory {
	return &accountHistory{
		table: `account_history`,
		db:    svcCtx.GormPointer,
		cache: svcCtx.Cache,
	}
}