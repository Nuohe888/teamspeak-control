package vo

type TsInfoRes struct {
	Domain           string `json:"domain"`
	DefaultVoicePort string `json:"defaultVoicePort"`
	QueryPort        string `json:"queryPort"`
	FiletransferPort string `json:"filetransferPort"`
}

type TsStatusRes struct {
	Status string `json:"status"`
}

type TsUuidReq struct {
	Uuid string `json:"uuid"`
}
