package api

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"license/jrebel/constant"
	"license/jrebel/vo"
	"log"
	"net/http"
	"strconv"
)

// LeasesController defines the controller structure
type LeasesController struct {
}

// NewLeasesController creates a new controller instance
func NewLeasesController() *LeasesController {
	return &LeasesController{}
}

// sign creates a digital signature using RSA-SHA1 algorithm.
func sign(clientRandomness, guid string, offline bool, validFrom, validUntil int64) string {
	signatureBase := clientRandomness + ";" + constant.ServerRandomness + ";" + guid + ";" + strconv.FormatBool(offline)
	if offline {
		// signatureBase += ";" + strconv.FormatInt(*validFrom, 10) + ";" + strconv.FormatInt(*validUntil, 10)
		signatureBase += ";" + strconv.FormatInt(validFrom, 10) + ";" + strconv.FormatInt(validUntil, 10)
	}

	log.Printf("signature: %s", signatureBase)

	block, _ := pem.Decode([]byte(constant.LeasesPrivateKey))
	if block == nil {
		log.Println("Failed to decode PEM block containing the private key")
		return ""
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("Failed to parse RSA private key: %v", err)
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(signatureBase))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)
	if err != nil {
		log.Printf("Failed to sign data: %v", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}

// LeasesHandler handles the "/leases" endpoint.
func (controller *LeasesController) LeasesHandler(c *gin.Context) {
	clientRandomness := c.PostForm("randomness")
	username := c.PostForm("username")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)

	var validFrom, validUntil int64
	if offline {
		// Calculate time after 180 days, note that the unit is milliseconds
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}

	signature := sign(clientRandomness, guid, offline, validFrom, validUntil)

	leasesHandlerVO := vo.LeasesHandlerVO{
		ServerVersion:         constant.ServerVersion,
		ServerProtocolVersion: constant.ServerProtocolVersion,
		ServerGuid:            constant.ServerGuid,
		GroupType:             constant.GroupType,
		ID:                    1,
		LicenseType:           1,
		EvaluationLicense:     false,
		Signature:             signature,
		ServerRandomness:      constant.ServerRandomness,
		SeatPoolType:          constant.SeatPoolType,
		StatusCode:            constant.StatusCode,
		Offline:               offline,
		ValidFrom:             validFrom,
		ValidUntil:            validUntil,
		Company:               username,
		OrderId:               uuid.NewString(),
		ZeroIds:               make([]string, 0),
		LicenseValidFrom:      1490544001000,
		LicenseValidUntil:     4102415999000,
	}

	c.JSON(http.StatusOK, leasesHandlerVO)
}

// Leases1Handler handles the "/leases/1" endpoint.
func (controller *LeasesController) Leases1Handler(c *gin.Context) {
	username := c.DefaultQuery("username", "")

	clientRandomness := c.PostForm("randomness")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)

	var validFrom, validUntil int64
	if offline {
		// Calculate time after 180 days, note that the unit is milliseconds
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}

	signature := sign(clientRandomness, guid, offline, validFrom, validUntil)

	leasesOneHandlerVO := vo.LeasesOneHandlerVO{
		ServerVersion:         constant.ServerVersion,
		ServerProtocolVersion: constant.ServerProtocolVersion,
		ServerGuid:            constant.ServerGuid,
		Signature:             signature,
		ServerRandomness:      constant.ServerRandomness,
		Features:              "{}",
		GroupType:             constant.GroupType,
		StatusCode:            constant.StatusCode,
		Company:               username,
		Msg:                   "",
		StatusMessage:         "",
	}

	c.JSON(http.StatusOK, leasesOneHandlerVO)
}

// ValidateHandler handles the "/validate-connection" endpoint.
func (controller *LeasesController) ValidateHandler(c *gin.Context) {

	clientRandomness := c.PostForm("randomness")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)

	var validFrom, validUntil int64
	if offline {
		// Calculate time after 180 days, note that the unit is milliseconds
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}

	signature := sign(clientRandomness, guid, offline, validFrom, validUntil)

	validateHandlerVO := vo.ValidateHandlerVO{
		ServerVersion:         constant.ServerVersion,
		ServerProtocolVersion: constant.ServerProtocolVersion,
		ServerGuid:            constant.ServerGuid,
		Signature:             signature,
		ServerRandomness:      constant.ServerRandomness,
		Features:              "{}",
		GroupType:             constant.GroupType,
		StatusCode:            constant.StatusCode,
		Company:               constant.Company,
		CanGetLease:           true,
		LicenseType:           "1",
		EvaluationLicense:     false,
		SeatPoolType:          constant.SeatPoolType,
	}

	c.JSON(http.StatusOK, validateHandlerVO)
}
