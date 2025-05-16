package ts

type Server struct {
	Model
	Id        *uint   `gorm:"primarykey" json:"id"`
	Name      *string `json:"name"`
	Host      *string `json:"host"`
	Port      *int    `json:"port"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	PortRange *string `json:"portRange"` //如30000-60000
	ImageName *string `json:"imageName"`
	Domain    *string `json:"domain"` //Domain域名.隔开
	Status    *int    `json:"status"`
}

func (Server) TableName() string {
	return "ts_server"
}
