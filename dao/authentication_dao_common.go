package dao

import "github.com/RedHatInsights/sources-api-go/config"

/*
	Common DAO code for the 1..n authentication DAOs we support.
	Currently, that list is:
	1. DB Dao, using postgres and an encrypted column as a backend
	2. Hashicorp vault, as the name implies storing authentication things in the vault
	3. Amazon Secrets Manager, a WIP that uses Amazon Secrets Manager to store authentication things
*/

// init sets the default DAO implementation so that other packages can request it easily.
func init() {
	GetAuthenticationDao = getDefaultAuthenticationDao
}

// GetAuthenticationDao is a function definition that can be replaced in runtime in case some other DAO provider is
// needed.
var GetAuthenticationDao func(daoParams *RequestParams) AuthenticationDao

// getDefaultAuthenticationDao gets the default DAO implementation which will have the given tenant ID.
func getDefaultAuthenticationDao(daoParams *RequestParams) AuthenticationDao {
	switch config.Get().SecretStore {
	case config.DatabaseStore:
		return &authenticationDaoDbImpl{RequestParams: daoParams}
	case config.VaultStore:
		return &authenticationDaoVaultImpl{RequestParams: daoParams}
	case config.SecretsManagerStore:
		return &authenticationSecretsManagerDaoImpl{
			RequestParams:           daoParams,
			authenticationDaoDbImpl: authenticationDaoDbImpl{RequestParams: daoParams},
		}
	default:
		return &noSecretStoreAuthenticationDao{}
	}
}
