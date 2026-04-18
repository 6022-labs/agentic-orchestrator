package enums

import (
	"fmt"
	"strings"
)

type ParticipantRole string
type ParticipantRoleList []ParticipantRole

const (
	ParticipantRoleClone       ParticipantRole = "clone"
	ParticipantRoleHuman       ParticipantRole = "human"
	ParticipantRoleExpert      ParticipantRole = "expert"
	ParticipantRoleFacilitator ParticipantRole = "facilitator"
)

func GetAllParticipantRoles() ParticipantRoleList {
	return ParticipantRoleList{
		ParticipantRoleClone,
		ParticipantRoleHuman,
		ParticipantRoleExpert,
		ParticipantRoleFacilitator,
	}
}

func (k ParticipantRole) String() string {
	return string(k)
}

func (k ParticipantRoleList) String() string {
	var keys []string
	for _, v := range k {
		keys = append(keys, string(v))
	}
	return strings.Join(keys, ", ")
}

func (k *ParticipantRole) Set(value string) error {
	normalizedValue := strings.ToLower(value)
	for _, v := range GetAllParticipantRoles() {
		if normalizedValue == v.String() {
			*k = v
			return nil
		}
	}
	return fmt.Errorf("invalid ParticipantRole: %s", value)
}

func (k ParticipantRoleList) Contains(key ParticipantRole) bool {
	for _, v := range k {
		if v == key {
			return true
		}
	}
	return false
}
