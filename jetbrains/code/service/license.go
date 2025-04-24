package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"license/jetbrains/code/dto"
	"license/jetbrains/code/mapper"
	"license/jetbrains/util"
	"license/logger"
	"time"
)

var productCodes = []string{"II", "PS", "AC", "DB", "RM", "WS", "RD", "CL", "PC", "GO", "DS", "DC", "DPN", "DM"}

func GenerateLicense(licenseeName, effectiveDate string, codes []string) (string, error) {

	if effectiveDate == "" {
		// Get current time
		now := time.Now()

		// Today at 23:59:59
		endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

		// Same day three years later
		threeYearsLater := endOfToday.AddDate(3, 0, 0)
		// Format date
		effectiveDate = threeYearsLater.Format("2006-01-02 15:04:05")
	}

	if len(codes) == 0 {
		productMapper := mapper.GormProductMapper{}
		products, err := productMapper.List()
		if err != nil {
			return "", err
		}
		for _, product := range products {
			codes = append(codes, product.ProductCode)
		}

		pluginMapper := mapper.GormPluginMapper{}
		plugins, err := pluginMapper.List()
		if err != nil {
			return "", err
		}
		for _, plugin := range plugins {
			codes = append(codes, plugin.PluginCode)
		}

		// productCodes
		for _, item := range productCodes {
			codes = append(codes, item)
		}
	}

	licenseID, err := randomString(10)
	if err != nil {
		logger.Error("Failed to generate license ID:", err)
		return "", err
	}
	var products []dto.Product
	for _, code := range codes {
		products = append(products, dto.Product{
			Code:         code,
			FallbackDate: effectiveDate,
			PaidUpTo:     effectiveDate,
			Extended:     true,
		})
	}
	licensePart := dto.LicensePart{
		LicenseID:         licenseID,
		LicenseeName:      licenseeName,
		Products:          products,
		AssigneeName:      "",
		Metadata:          "0120231110PSAA003008",
		Hash:              "51149839/0:-1370131430",
		GracePeriodDays:   7,
		AutoProlongated:   true,
		IsAutoProlongated: true,
		Trial:             false,
		AiAllowed:         true,
	}

	licensePartJSON, err := json.Marshal(licensePart)
	if err != nil {
		return "", err
	}
	licensePartBase64 := base64.StdEncoding.EncodeToString(licensePartJSON)
	fakeCert := util.GetFake()
	signatureBase64 := signWithRSA(fakeCert.PrivateKey, licensePartJSON)
	cert := fakeCert.CodeCert
	certBase64 := base64.StdEncoding.EncodeToString(cert.Raw)
	println(signatureBase64)
	return licenseID + "-" + licensePartBase64 + "-" + signatureBase64 + "-" + certBase64, nil
}

func signWithRSA(privateKey *rsa.PrivateKey, data []byte) string {
	// Hash the data with SHA1
	hash := sha1.New()
	_, err := hash.Write(data)
	if err != nil {
		logger.Error("Failed to hash: ", err)
		return ""
	}
	hashed := hash.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)
	if err != nil {
		logger.Error("Failed to sign: ", err)
		return ""
	}
	signature := base64.StdEncoding.EncodeToString(sign)
	return signature
}

// randomString generates a random string of length n, containing uppercase letters and numbers
func randomString(n int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	byteArray := make([]byte, n) // Create a byte slice of length n

	// Read random bytes to fill byteArray
	if _, err := rand.Read(byteArray); err != nil {
		return "", err // Return error message
	}

	// Map each byte to a character in charset, ensuring it's an uppercase letter or number
	for i, b := range byteArray {
		byteArray[i] = charset[b%byte(len(charset))] // Use modulo to convert byte to a character in charset
	}

	return string(byteArray), nil // Convert to string and return
}
