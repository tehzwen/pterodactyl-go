package websocket

type SocketEventType string

var (
	SocketEventTypeAuth          SocketEventType = "auth"
	SocketEventTypeAuthSucess    SocketEventType = "auth success"
	SocketEventTypeSendCommand   SocketEventType = "send command"
	SocketEventTypeSetState      SocketEventType = "set state"
	SocketEventTypeConsoleOutput SocketEventType = "console output"
	SocketEventTypeInstallLogs   SocketEventType = "install logs"
	SocketEventTypeTransferLogs  SocketEventType = "transfer logs"
	SocketEventTypeStatus        SocketEventType = "status"
	SocketEventTypeStats         SocketEventType = "stats"
	SocketEventTypeJWTError      SocketEventType = "jwt error"
	SocketEventTypeDaemonMessage SocketEventType = "daemon message"
	SocketEventTypeTokenExpiring SocketEventType = "token expiring"
)

type SocketMessage struct {
	Event SocketEventType `json:"event"`
	Args  []string        `json:"args"`
}
