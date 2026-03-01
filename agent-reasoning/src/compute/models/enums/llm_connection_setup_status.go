package enums

import "strings"

type LLMConnectionSetupStatus string
type LLMConnectionSetupStatusList []LLMConnectionSetupStatus

const (
	LLMConnectionSetupStatusEstablished LLMConnectionSetupStatus = "established"
)

func GetAllLLMConnectionSetupStatuses() LLMConnectionSetupStatusList {
	return LLMConnectionSetupStatusList{
		LLMConnectionSetupStatusEstablished,
	}
}

func (k *LLMConnectionSetupStatus) ToPointer() *LLMConnectionSetupStatus {
	if k == nil {
		return nil
	}
	return k
}

func (k LLMConnectionSetupStatus) String() string {
	return string(k)
}

func (k LLMConnectionSetupStatusList) String() string {
	var statuses []string
	for _, v := range k {
		statuses = append(statuses, string(v))
	}
	return strings.Join(statuses, ", ")
}

func (k LLMConnectionSetupStatusList) Contains(status LLMConnectionSetupStatus) bool {
	for _, v := range k {
		if v == status {
			return true
		}
	}
	return false
}
