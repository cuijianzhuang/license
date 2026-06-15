package gitlab

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"license/internal/config"
	"license/internal/crypto"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	gorsa "github.com/Lyafei/go-rsa"
	"github.com/gin-gonic/gin"
)

// License quantity granted by the generated license.
const licenseQuantity = 1000

// RSA keys are loaded lazily from disk on first use. Both vars are written
// exactly once inside loadKeys.Do — sync.Once gives us the happens-before
// guarantee, so no additional mutex is needed.
var (
	loadKeys      sync.Once
	loadKeysErr   error
	privateKeyPEM []byte
	publicKeyPEM  []byte
)

func getKeys() ([]byte, []byte, error) {
	loadKeys.Do(func() {
		dir := config.GetConfig().DataDir
		pub, err := os.ReadFile(dir + "/.license_encryption_key.pub")
		if err != nil {
			loadKeysErr = fmt.Errorf("failed to read public key: %v", err)
			return
		}
		priv, err := os.ReadFile(dir + "/.license_decryption_key.pri")
		if err != nil {
			loadKeysErr = fmt.Errorf("failed to read private key: %v", err)
			return
		}
		publicKeyPEM = pub
		privateKeyPEM = priv
	})
	return privateKeyPEM, publicKeyPEM, loadKeysErr
}

// LoadKeys eagerly loads the GitLab keypair so missing files surface during
// startup instead of on the first request.
func LoadKeys() error {
	_, _, err := getKeys()
	return err
}

// newLicenseIdentifiers returns a freshly randomized set of identifiers (restriction id,
// subscription id/name and per-add-on purchase XIDs) so every generated license carries
// a unique identity. A single rand.Read call backs all values.
func newLicenseIdentifiers() (restrictionID, subID, subName, duoXID, dapXID string, err error) {
	var buf [40]byte
	if _, err = rand.Read(buf[:]); err != nil {
		return "", "", "", "", "", err
	}
	restrictionID = fmt.Sprintf("%x", buf[0:8])
	subID = fmt.Sprintf("offline-subscription-%x", buf[8:16])
	subName = fmt.Sprintf("ultimate-duo-self-hosted-%x", buf[16:24])
	duoXID = fmt.Sprintf("duo-enterprise-%x", buf[24:32])
	dapXID = fmt.Sprintf("self-hosted-dap-%x", buf[32:40])
	return
}

// createLicenseJson creates a JSON representation of the license with optimized buffer usage.
// The payload follows GitLab's cloud / offline-cloud licensing schema (GitLab 17+):
// it carries activation timestamps, subscription metadata and add_on_products instead of
// the legacy top-level features whitelist.
func createLicenseJson(licenseInfo LicenseInfo, expireTime string) ([]byte, error) {

	var expirationDate time.Time
	var err error
	if len(expireTime) == 0 {
		// Default expiration time is 2 years
		expirationDate = time.Now().AddDate(2, 0, 0)
	} else {
		expirationDate, err = time.Parse(time.DateTime, expireTime)
		if err != nil {
			log.Printf("Failed to parse expiration time: %v", err)
			return nil, err
		}
	}

	restrictionID, subID, subName, duoXID, dapXID, err := newLicenseIdentifiers()
	if err != nil {
		log.Printf("Failed to generate license identifiers: %v", err)
		return nil, err
	}

	now := time.Now()
	// activated_at is conventionally the day after issuance, normalized to 00:00 UTC.
	activatedAt := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)

	addOnProducts := map[string][]AddOnPurchase{
		"duo_enterprise": {{
			Quantity:    licenseQuantity,
			StartedOn:   CustomTime{Time: now},
			ExpiresOn:   CustomTime{Time: expirationDate},
			PurchaseXID: duoXID,
		}},
		"self_hosted_dap": {{
			Quantity:    1,
			StartedOn:   CustomTime{Time: now},
			ExpiresOn:   CustomTime{Time: expirationDate},
			PurchaseXID: dapXID,
		}},
	}

	license := License{
		Version: 1,
		License: licenseInfo,

		IssuedAt:       CustomTime{Time: now},
		ExpiresAt:      CustomTime{Time: expirationDate},
		NotifyAdminsAt: CustomTime{Time: expirationDate.AddDate(0, 0, -30)},
		NotifyUsersAt:  CustomTime{Time: expirationDate.AddDate(0, 0, -7)},
		BlockChangesAt: CustomTime{Time: expirationDate.AddDate(0, 0, 180)},

		ActivatedAt:  CustomDateTime{Time: activatedAt},
		LastSyncedAt: CustomDateTime{Time: activatedAt},
		NextSyncAt:   CustomDateTime{Time: activatedAt.AddDate(0, 0, 90)},

		CloudLicensingEnabled:        true,
		OfflineCloudLicensingEnabled: true,
		AutoRenewEnabled:             false,
		SeatReconciliationEnabled:    true,
		OperationalMetricsEnabled:    true,
		GeneratedFromCustomersDot:    true,
		GeneratedFromCancellation:    false,
		TemporaryExtension:           false,
		ContractOveragesAllowed:      true,

		Restrictions: Restriction{
			ID:                      restrictionID,
			Plan:                    "ultimate",
			ActiveUserCount:         licenseQuantity,
			SubscriptionID:          subID,
			SubscriptionName:        subName,
			ReconciliationCompleted: true,
			AddOnProducts:           addOnProducts,
		},
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	// Match Ruby's JSON.dump output: don't escape <, >, & as \u003c, \u003e, \u0026.
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(license); err != nil {
		return nil, err
	}

	// Remove trailing newline added by encoder
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	return data, nil
}

// generateRandomIV generates a random initialization vector (IV)
func generateRandomIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize) // AES block size is fixed at 16 bytes
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// Encrypt wraps the Encrypt method, using AES-CBC encryption and PKCS7 padding
func Encrypt(data, key, iv []byte) ([]byte, error) {
	aesTool := crypto.AesCbcPkcs7{Key: key, Iv: iv}
	enc, err := aesTool.Encrypt(data)
	if err != nil {
		log.Println("Encrypt error:", err)
		return nil, err
	}
	return enc, err
}

// Uses RSA private key to "encrypt" data with cached keys
func encryptWithPrivateKey(data string) (string, error) {
	privateKey, _, err := getKeys()
	if err != nil {
		return "", err
	}
	encrypt, err := gorsa.PriKeyEncrypt(data, string(privateKey))
	if err != nil {
		log.Printf("Failed to encrypt data with RSA private key: %v", err)
		return "", err
	}
	return encrypt, nil
}

// encryptLicense encrypts license data using AES and RSA.
func encryptLicense(data []byte) (string, error) {
	key := make([]byte, 16)
	iv := make([]byte, aes.BlockSize)

	if _, err := rand.Read(key); err != nil {
		log.Printf("Failed to generate AES key: %v", err)
		return "", err
	}
	if _, err := rand.Read(iv); err != nil {
		log.Printf("Failed to generate AES IV: %v", err)
		return "", err
	}

	encryptedData, err := Encrypt(data, key, iv)
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		return "", err
	}

	// Note: RSA encryption is typically done with a public key, but technically can be done with a private key (although not recommended)
	encryptedKey, err := encryptWithPrivateKey(string(key))
	if err != nil {
		return "", err
	}

	// Encode encrypted data as Base64
	encryptedDataStr := base64.StdEncoding.EncodeToString(encryptedData)
	ivStr := base64.StdEncoding.EncodeToString(iv)

	// Package as JSON format
	result := map[string]string{
		"data": encryptedDataStr,
		"key":  encryptedKey,
		"iv":   ivStr,
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to package JSON data: %v", err)
		return "", err
	}

	// Encode JSON as Base64
	encodedFinal := base64.StdEncoding.EncodeToString(jsonData)
	return encodedFinal, nil
}

// Generate generates a license and sends it via HTTP response
func Generate(ctx *gin.Context, licenseInfo LicenseInfo, expireTime string) {
	createLicense(ctx, licenseInfo, expireTime)
}

// createLicense creates and sends a license. Responses are intentionally not cached:
// the payload carries a freshly randomized subscription_id / subscription_name on every
// call, so each download must be regenerated end-to-end.
func createLicense(ctx *gin.Context, licenseInfo LicenseInfo, expireTime string) {
	// Create license JSON data
	licenseJson, err := createLicenseJson(licenseInfo, expireTime)
	if err != nil {
		log.Printf("Failed to create license JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Encrypt the license data
	encryptedLicense, err := encryptLicense(licenseJson)
	if err != nil {
		log.Printf("Failed to encrypt license: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create ZIP file in memory first for caching
	buf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buf)

	// Add public key file to ZIP
	if err := addFileToZipOptimized(zipWriter, config.GetConfig().DataDir+"/.license_encryption_key.pub", "license/.license_encryption_key.pub"); err != nil {
		log.Printf("Failed to add public key to ZIP: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Add encrypted license data to ZIP
	if err := addLicenseToZip(zipWriter, encryptedLicense, "license/license.gitlab-license"); err != nil {
		log.Printf("Failed to add license to ZIP: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := zipWriter.Close(); err != nil {
		log.Printf("Failed to close ZIP writer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	zipData := buf.Bytes()

	// Send response
	ctx.Header("Content-Disposition", "attachment; filename=license.zip")
	ctx.Header("Content-Type", "application/zip")
	ctx.Data(http.StatusOK, "application/zip", zipData)
}

// addFileToZipOptimized reads a file from the filesystem and adds it to the ZIP with buffered I/O
func addFileToZipOptimized(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = zipPath
	header.Method = zip.Deflate
	header.Modified = fileInfo.ModTime()

	zipFile, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, file)
	return err
}

// addLicenseToZip directly writes string data to a ZIP entry
func addLicenseToZip(zipWriter *zip.Writer, data, zipPath string) error {
	// Create a new zip.FileHeader, set filename and modification time
	header := &zip.FileHeader{
		Name:     zipPath,
		Method:   zip.Deflate, // Use compression to reduce file size
		Modified: time.Now(),  // Set current time as file modification time
	}

	// Create ZIP entry
	zipFile, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write data to ZIP entry
	_, err = zipFile.Write([]byte(data))
	return err
}
