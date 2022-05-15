package serverConsts

const (
	ConfigFileName = "server.yaml"
	CasbinFileName = "rbac_model.conf" // casbin 规则文件名称

	WebCheckInterval = 60 * 60

	WsDefaultNameSpace = "default"
	WsDefaultRoom      = "default"
	WsEvent            = "OnChat"
	WsMsgEvent         = "eventWebSocketMsg"

	ApiPath = "/api/v1"
	WsPath  = ApiPath + "/ws"

	AdminUserName     = "admin"
	AdminUserPassword = "password"
	AdminRoleName     = "admin"
)

var (
	SortMap = map[string]string{
		"ascend":  "ASC",
		"descend": "DESC",
		"":        "ASC",
	}
)
