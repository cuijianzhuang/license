package impl

import (
	"license/internal/jetbrains/server/protocol"
	"license/internal/jetbrains/util"
	"license/internal/rpc/service"

	"github.com/gin-gonic/gin"
)

// JetbrainsRpcService Service represents the JRebel RPC service.
type JetbrainsRpcService struct {
}

// Ensure Service implements the common.RpcService interface
var _ service.RpcService = &JetbrainsRpcService{}

// Ping handles the ping XML request and returns a signed response.
func (s *JetbrainsRpcService) Ping(ctx *gin.Context, machineId, salt string) {
	pingReq := &protocol.BaseRequest{
		Salt:      salt,
		MachineId: machineId,
	}
	pingResponse := protocol.NewPingResponse(pingReq, util.GetFake())
	ctx.Render(200, protocol.NewXMLTicket(pingResponse))
}

// ObtainTicket handles the ticket obtaining XML request and returns a signed response.
func (s *JetbrainsRpcService) ObtainTicket(ctx *gin.Context, username, hostName, machineId, salt string) {
	if len(username) == 0 {
		username = hostName
	}
	obtainReq := &protocol.BaseRequest{
		Salt:      salt,
		UserName:  username,
		MachineId: machineId,
	}
	ticketResponse := protocol.NewObtainTicketResponse(obtainReq, util.GetFake())
	ctx.Render(200, protocol.NewXMLTicket(ticketResponse))
}

// ReleaseTicket handles the ticket release XML request and returns a signed response.
func (s *JetbrainsRpcService) ReleaseTicket(ctx *gin.Context, machineId, salt string) {
	releaseReq := &protocol.BaseRequest{
		Salt:      salt,
		MachineId: machineId,
	}
	releaseTicketResponse := protocol.NewReleaseTicketResponse(releaseReq, util.GetFake())
	ctx.Render(200, protocol.NewXMLTicket(releaseTicketResponse))
}
