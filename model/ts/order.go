package ts

type Order struct {
	Model
	Uuid               *string `gorm:"primarykey" json:"uuid"`
	ExpTime            *int    `json:"expTime"`
	Port               *string `json:"port"`
	Status             *int    `json:"status"` //-1:未知 0:创建 1:启动 2:停止
	ServerId           *uint   `json:"serverId"`
	TsDefaultVoicePort *string `json:"tsDefaultVoicePort"`
	TsQueryPort        *string `json:"tsQueryPort"`
	TsFiletransferPort *string `json:"tsFiletransferPort"`
	TsApikey           *string `json:"tsApikey"`
	TsLoginName        *string `json:"tsLoginName"`
	TsPassword         *string `json:"tsPassword"`
	TsToken            *string `json:"tsToken"`
}

func (Order) TableName() string {
	return "ts_order"
}
