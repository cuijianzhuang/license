package service

import "github.com/gin-gonic/gin"

// RpcService defines methods related to remote procedure calls
type RpcService interface {
	Ping(ctx *gin.Context, machineId, salt string)
	ObtainTicket(ctx *gin.Context, username, hostName, machineId, salt string)
	ReleaseTicket(ctx *gin.Context, machineId, salt string)
}
