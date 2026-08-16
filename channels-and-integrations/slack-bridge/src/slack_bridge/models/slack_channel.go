package models

type SlackChannel struct {
	Id          string       `gorm:"primaryKey;type:string;not null"                                              json:"id"`
	Name        string       `gorm:"type:string;not null"                                                         json:"name"`
	Description string       `gorm:"type:string"                                                                  json:"description"`
	Agents      []SlackAgent `gorm:"many2many:slack_channel_agents;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"agents,omitempty"`
	Users       []SlackUser  `gorm:"many2many:slack_channel_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"users,omitempty"`
}
