// Package rpc implements the /rpc/*.action endpoints used by JetBrains IDEs
// (and historically by JRebel) to ping the license server and obtain/release
// tickets.
//
// The handlers dispatch on the presence of `machineId`: JetBrains clients
// always include it, so an empty value means the caller is a JRebel agent.
package rpc

import (
	jbserver "license/internal/jetbrains/server"
	"license/internal/jetbrains/util"

	"github.com/gin-gonic/gin"
)

// Controller wires the rpc handlers into Gin. It holds no per-request state;
// the struct exists so the router can keep its current `NewController()` shape.
type Controller struct{}

// NewController creates a new rpc controller.
func NewController() *Controller {
	return &Controller{}
}

// Ping handles /rpc/ping.action.
func (c *Controller) Ping(ctx *gin.Context) {
	machineID := ctx.DefaultQuery("machineId", "")
	salt := ctx.Query("salt")

	if machineID != "" {
		jetbrainsPing(ctx, machineID, salt)
		return
	}
	jrebelPing(ctx, salt)
}

// ObtainTicket handles /rpc/obtainTicket.action.
func (c *Controller) ObtainTicket(ctx *gin.Context) {
	username := ctx.DefaultQuery("userName", "")
	if username == "" {
		username = ctx.DefaultQuery("username", "")
	}
	hostName := ctx.DefaultQuery("hostName", "")
	machineID := ctx.DefaultQuery("machineId", "")
	salt := ctx.Query("salt")

	if machineID != "" {
		jetbrainsObtainTicket(ctx, username, hostName, machineID, salt)
		return
	}
	jrebelObtainTicket(ctx, username, salt)
}

// ReleaseTicket handles /rpc/releaseTicket.action.
func (c *Controller) ReleaseTicket(ctx *gin.Context) {
	machineID := ctx.DefaultQuery("machineId", "")
	salt := ctx.Query("salt")

	if machineID != "" {
		jetbrainsReleaseTicket(ctx, machineID, salt)
		return
	}
	jrebelReleaseTicket(ctx, salt)
}

func jetbrainsPing(ctx *gin.Context, machineID, salt string) {
	req := &jbserver.BaseRequest{Salt: salt, MachineId: machineID}
	ctx.Render(200, jbserver.NewXMLTicket(jbserver.NewPingResponse(req, util.GetFake())))
}

func jetbrainsObtainTicket(ctx *gin.Context, username, hostName, machineID, salt string) {
	if username == "" {
		username = hostName
	}
	req := &jbserver.BaseRequest{Salt: salt, UserName: username, MachineId: machineID}
	ctx.Render(200, jbserver.NewXMLTicket(jbserver.NewObtainTicketResponse(req, util.GetFake())))
}

func jetbrainsReleaseTicket(ctx *gin.Context, machineID, salt string) {
	req := &jbserver.BaseRequest{Salt: salt, MachineId: machineID}
	ctx.Render(200, jbserver.NewXMLTicket(jbserver.NewReleaseTicketResponse(req, util.GetFake())))
}
