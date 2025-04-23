package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ebfe/keccak"
)

func md5Hash(msg string) string {
	hash := md5.New()
	hash.Write([]byte(msg))
	str := hex.EncodeToString(hash.Sum(nil))
	return str[8:24]
}

func keccak384Hash(msg string) string {
	enc := keccak.New384()
	enc.Write([]byte(msg))
	str := hex.EncodeToString(enc.Sum(nil))
	return str[12:28]
}

// GenerateLicense generates license codes for different FinalShell versions and editions
// Returns an array of license codes in the following order:
// [0] - License for Version < 3.9.6 Advanced Edition
// [1] - License for Version < 3.9.6 Professional Edition
// [2] - License for 3.9.6 <= Version < 4.5 Advanced Edition
// [3] - License for 3.9.6 <= Version < 4.5 Professional Edition
// [4] - License for Version >= 4.5 Advanced Edition
// [5] - License for Version >= 4.5 Professional Edition
// Note: Version prefixes are no longer included in the returned strings.
// Frontend must handle the display of version information.
func GenerateLicense(machineCode string) []string {
	var result []string
	// License for Version < 3.9.6 Advanced Edition
	result = append(result, md5Hash("61305"+machineCode+"8552"))
	// License for Version < 3.9.6 Professional Edition
	result = append(result, md5Hash("2356"+machineCode+"13593"))
	// License for 3.9.6 <= Version < 4.5 Advanced Edition
	result = append(result, keccak384Hash(machineCode+"hSf(78cvVlS5E"))
	// License for3.9.6 <= Version < 4.5 Professional Edition
	result = append(result, keccak384Hash(machineCode+"FF3Go(*Xvbb5s2"))
	// License for Version >= 4.5 Advanced Edition
	result = append(result, keccak384Hash(machineCode+"wcegS3gzA$"))
	// License for Version >= 4.5 Professional Edition
	result = append(result, keccak384Hash(machineCode+"b(xxkHn%z);x"))
	return result
}
