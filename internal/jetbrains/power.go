package jetbrains

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"math/big"
)

// GeneratePowerResult computes the EQUAL,x,y,z->r tuple used by power.conf to
// trick the IDE into accepting our signed certificate. x = cert.Signature as a
// big int, (y, z) = root CA RSA (E, N), r = x^y mod cert.PublicKey.N.
func GeneratePowerResult(cert, rootCA *x509.Certificate) string {
	x := (&big.Int{}).SetBytes(cert.Signature)
	z := rootCA.PublicKey.(*rsa.PublicKey).N
	y := rootCA.PublicKey.(*rsa.PublicKey).E
	r := &big.Int{}
	r.Exp(x, big.NewInt(int64(y)), cert.PublicKey.(*rsa.PublicKey).N)

	return fmt.Sprintf("EQUAL,%d,%d,%d->%d", x, y, z, r)
}
