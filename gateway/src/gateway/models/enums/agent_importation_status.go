package enums

type AgentImportationStatus string
type AgentImportationStatusList []AgentImportationStatus

const (
	AgentImportationStatusFirstStep  AgentImportationStatus = "first-step"
	AgentImportationStatusSecondStep AgentImportationStatus = "second-step"
	AgentImportationStatusThirdStep  AgentImportationStatus = "third-step"
	AgentImportationStatusFinished   AgentImportationStatus = "finished"
)

func GetAllAgentStatus() AgentImportationStatusList {
	return AgentImportationStatusList{
		AgentImportationStatusFirstStep,
		AgentImportationStatusSecondStep,
		AgentImportationStatusThirdStep,
		AgentImportationStatusFinished,
	}
}

func (a AgentImportationStatus) ToPointer() *AgentImportationStatus {
	return &a
}

func (l AgentImportationStatusList) Contains(status AgentImportationStatus) bool {
	for _, v := range l {
		if v == status {
			return true
		}
	}
	return false
}
