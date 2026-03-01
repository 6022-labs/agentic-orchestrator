package enums

import "strings"

type DatabaseConnectionSetupStatus string
type DatabaseConnectionSetupStatusList []DatabaseConnectionSetupStatus

const (
	DatabaseConnectionSetupStatusConfigured   DatabaseConnectionSetupStatus = "configured"
	DatabaseConnectionSetupStatusEstablished  DatabaseConnectionSetupStatus = "established"
	DatabaseConnectionSetupStatusEstablishing DatabaseConnectionSetupStatus = "establishing"
)

func GetAllDatabaseConnectionSetupStatuses() DatabaseConnectionSetupStatusList {
	return DatabaseConnectionSetupStatusList{
		DatabaseConnectionSetupStatusConfigured,
		DatabaseConnectionSetupStatusEstablished,
		DatabaseConnectionSetupStatusEstablishing,
	}
}

func (k *DatabaseConnectionSetupStatus) ToPointer() *DatabaseConnectionSetupStatus {
	if k == nil {
		return nil
	}
	return k
}

func (k DatabaseConnectionSetupStatus) String() string {
	return string(k)
}

func (k DatabaseConnectionSetupStatusList) String() string {
	var statuses []string
	for _, v := range k {
		statuses = append(statuses, string(v))
	}
	return strings.Join(statuses, ", ")
}

func (k DatabaseConnectionSetupStatusList) Contains(status DatabaseConnectionSetupStatus) bool {
	for _, v := range k {
		if v == status {
			return true
		}
	}
	return false
}
