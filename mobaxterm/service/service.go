package service

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"license/config"
	"license/logger"
	"license/utils/useragent"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"os"
)

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// FetchVersions fetches MobaXterm version numbers from the website
func FetchVersions() ([]string, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Create request
	req, err := http.NewRequest("GET", "https://mobaxterm.mobatek.net/download-home-edition.html", nil)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to create request: %v", err))
		return nil, err
	}

	// Get and set random User-Agent along with standard headers
	headers := useragent.GetRandomWithAcceptHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to fetch MobaXterm website: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("", fmt.Errorf("bad status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to parse HTML: %v", err))
		return nil, err
	}

	// Create a set to store unique versions
	versionsMap := make(map[string]bool)

	// Strategy 1: Try to get current version from the download links
	// Use a more specific regex to match "Version X.Y (date)" pattern
	versionRegex := regexp.MustCompile(`Version (\d+\.\d+)`)

	// Strategy 2: Look for version_titre class - this is our primary source
	doc.Find("p.version_titre").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		matches := versionRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				versionsMap[match[1]] = true
			}
		}
	})

	// Convert map to sorted slice
	versions := make([]string, 0, len(versionsMap))
	for version := range versionsMap {
		versions = append(versions, version)
	}

	// If we found no versions, return default list
	if len(versions) == 0 {
		logger.Error("", fmt.Errorf("no versions found, returning default list"))
		return []string{"25.1", "25.0", "24.4", "24.3", "23.6", "23.5", "23.0"}, nil
	}

	// Sort versions in descending order (newer versions first)
	sort.Slice(versions, func(i, j int) bool {
		// Split version into major and minor components
		partsI := strings.Split(versions[i], ".")
		partsJ := strings.Split(versions[j], ".")

		// Parse major version
		majorI, _ := strconv.Atoi(partsI[0])
		majorJ, _ := strconv.Atoi(partsJ[0])

		if majorI != majorJ {
			return majorI > majorJ
		}

		// Parse minor version
		minorI, _ := strconv.Atoi(partsI[1])
		minorJ, _ := strconv.Atoi(partsJ[1])

		return minorI > minorJ
	})

	return versions, nil
}

// GenerateLicense generates a license for MobaXterm
func GenerateLicense(count int, username, version string, c *gin.Context) {
	// Check if parameters are empty or invalid
	if username == "" || version == "" || count <= 0 {
		logger.Error("", fmt.Errorf("parameter error: %s, %s, %d", username, version, count))
	}

	// Split version number
	versionArr := strings.Split(version, ".")
	if len(versionArr) != 2 {
		logger.Error("", fmt.Errorf("version format error: %s", version))
	}

	// Check if version numbers are digits
	if _, err := strconv.Atoi(versionArr[0]); err != nil {
		// panic("版本号格式错误")
		log.Printf("version format error: %s", version)
	}
	if _, err := strconv.Atoi(versionArr[1]); err != nil {
		// panic("版本号格式错误")
		log.Printf("version format error: %s", version)
	}

	// Extract major and minor version numbers
	major, _ := strconv.ParseInt(versionArr[0], 10, 64)
	minor, _ := strconv.ParseInt(versionArr[1], 10, 64)

	license := generateLicense(1, count, username, major, minor)
	// First write to file, solving the issue where direct output of zip file has incorrect size and cannot be used
	toFile(license)
	// Read file, output to browser
	c.FileAttachment(config.GetConfig().DataDir+"/Custom.mxtpro", "Custom.mxtpro")
	// Delete file
	err := os.Remove(config.GetConfig().DataDir + "/Custom.mxtpro")
	if err != nil {
		logger.Error("", fmt.Errorf("failed to delete file: %v", err))
	}
}

func toFile(license []byte) {
	fileName := config.GetConfig().DataDir + "/Custom.mxtpro"
	_ = os.Remove(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to create file: %v", err))
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error("", fmt.Errorf("failed to close file: %v", err))
		}
	}(f)

	zipFile := zip.NewWriter(f)
	defer func(zipFile *zip.Writer) {
		err := zipFile.Close()
		if err != nil {
			logger.Error("", fmt.Errorf("failed to close ZIP writer: %v", err))
		}
	}(zipFile)
	header := &zip.FileHeader{
		Name:               "Pro.key",
		Method:             zip.Store,
		CompressedSize64:   38,
		UncompressedSize64: 38,
	}
	// TODO Key point, putting it in FileHeader will cause incorrect file size
	header.SetModTime(time.Now())
	proFile, err := zipFile.CreateRaw(header)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to create ZIP file: %v", err))
		return
	}
	_, err = proFile.Write(license)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to write to ZIP file: %v", err))
		return
	}
	return
}

func generateLicense(userType, count int, username string, major, minor int64) []byte {
	licenseString := fmt.Sprintf("%d#%s|%d%d#%d#%d3%d6%d#%d#%d#%d#", userType, username, major, minor, count, major, minor, minor, 0, 0, 0)
	return variantBase64Encode(encryptBytes(0x787, []byte(licenseString)))
}

func encryptBytes(key int, bs []byte) []byte {
	var result []byte
	for _, b := range bs {
		encryptedByte := b ^ byte((key>>8)&0xff)
		result = append(result, encryptedByte)
		key = (int(result[len(result)-1]) & key) | 0x482D
	}
	return result
}

var (
	variantBase64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	variantBase64Map   = func() map[int]byte {
		result := make(map[int]byte)
		for i, v := range variantBase64Table {
			result[i] = byte(v)
		}
		return result
	}()
)

func variantBase64Encode(bs []byte) []byte {
	var result []byte
	blocksCount := len(bs) / 3
	leftBytes := len(bs) % 3
	for i := 0; i < blocksCount; i++ {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*i : 3*i+3])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>12)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>18)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
	}
	if leftBytes == 0 {
		return result
	} else if leftBytes == 1 {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*blocksCount:])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
		return result
	} else {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*blocksCount:])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>12)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
		return result
	}
}

func littleEndianBytes(bs []byte) int {
	var result = int(bs[0])
	for i := 1; i < len(bs); i++ {
		result = result | int(bs[i])<<(8*i)
	}
	return result
}
