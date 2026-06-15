package jrebel

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	httpScheme  = "http"
	httpPort    = "80"
	httpsScheme = "https"
	httpsPort   = "443"
	endStr      = "/"
)

// IndexController defines the controller structure
type IndexController struct {
}

// NewIndexController creates a new controller instance
func NewIndexController() *IndexController {
	return &IndexController{}
}

// IndexHandler using Gin to handle the root endpoint and display license information.
func (controller *IndexController) IndexHandler(c *gin.Context) {
	scheme := httpScheme
	if c.Request.TLS != nil {
		scheme = httpsScheme
	}

	host := c.Request.Host
	if pos := strings.Index(host, ":"); pos != -1 {
		host = host[:pos]
	}
	if host == "" {
		host = "127.0.0.1"
	}

	port := c.Request.URL.Port()
	if port == "" {
		// default to HTTP port, this should be configured or detected
		port = "80"
	}

	requestURI := strings.TrimSuffix(c.Request.RequestURI, endStr)

	licenseURL := fmt.Sprintf("%s://%s%s", scheme, host, requestURI)
	if (scheme == httpScheme && port != httpPort) || (scheme == httpsScheme && port != httpsPort) {
		licenseURL = fmt.Sprintf("%s://%s:%s%s", scheme, host, port, requestURI)
	}

	html := fmt.Sprintf(`
		<h3>Instructions for use</h3>
		<hr/>
		<h1>Hello, This is a Jrebel & JetBrains License Server!</h1>
		<p>License Server started at %s</p>
		<p>JetBrains Activation address was: <span style='color:red'>%s</span></p>
		<p>JRebel 7.1 and earlier version Activation address was: <span style='color:red'>%s/{tokenname}</span>, with any email.</p>
		<p>JRebel 2018.1 and later version Activation address was: %s/{guid} (eg:<span style='color:red'>%s/%s</span>), with any email.</p>
		<hr/>
		<h1>Hello, this is a Jrebel & JetBrains License Server!</h1>
		<p>JetBrains license server activation address %s</p>
		<p>JetBrains activation address is: <span style='color:red'>%s</span></p>
		<p>JRebel 7.1 and older versions activation address: <span style='color:red'>%s/{tokenname}</span>, with any email address.</p>
		<p>JRebel 2018.1+ version activation address: %s/{guid} (example: <span style='color:red'>%s/%s</span>), with any email address.</p>
	`, licenseURL, licenseURL, licenseURL, licenseURL, licenseURL, c.Query("guid"), licenseURL, licenseURL, licenseURL, licenseURL, licenseURL, c.Query("guid"))

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		return
	}
}
