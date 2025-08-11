package impl

import (
	"license/jetbrains/server/entity"
	"license/jetbrains/util"
	"license/rpc/service"

	"github.com/gin-gonic/gin"
)

// JetbrainsRpcService Service represents the JRebel RPC service.
type JetbrainsRpcService struct {
}

// Ensure Service implements the common.RpcService interface
var _ service.RpcService = &JetbrainsRpcService{}

// Ping handles the ping XML request and returns a signed response.
func (s *JetbrainsRpcService) Ping(ctx *gin.Context, machineId, salt string) {
	pingReq := &entity.BaseRequest{
		Salt:      salt,
		MachineId: machineId,
	}
	pingResponse := entity.NewPingResponse(pingReq, util.GetFake())
	ctx.Render(200, entity.NewXMLTicket(pingResponse))
}

// ObtainTicket handles the ticket obtaining XML request and returns a signed response.
func (s *JetbrainsRpcService) ObtainTicket(ctx *gin.Context, username, hostName, machineId, salt string) {
	if len(username) == 0 {
		username = hostName
	}
	obtainReq := &entity.BaseRequest{
		Salt:      salt,
		UserName:  username,
		MachineId: machineId,
	}
	ticketResponse := entity.NewObtainTicketResponse(obtainReq, util.GetFake())
	ctx.Render(200, entity.NewXMLTicket(ticketResponse))
}

// ReleaseTicket handles the ticket release XML request and returns a signed response.
func (s *JetbrainsRpcService) ReleaseTicket(ctx *gin.Context, machineId, salt string) {
	releaseReq := &entity.BaseRequest{
		Salt:      salt,
		MachineId: machineId,
	}
	releaseTicketResponse := entity.NewReleaseTicketResponse(releaseReq, util.GetFake())
	ctx.Render(200, entity.NewXMLTicket(releaseTicketResponse))
}
