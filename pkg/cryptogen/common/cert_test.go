package common

import (
	"testing"
)

func TestNewCert(t *testing.T) {
	certPem := `-----BEGIN CERTIFICATE-----
MIICKTCCAc+gAwIBAgIRALdRir6uy+WuoKMbiVpFOFEwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzEuZXhhbXBsZS5jb20wHhcNMjAxMjE1MDcxNTAwWhcNMzAxMjEzMDcxNTAw
WjBqMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzENMAsGA1UECxMEcGVlcjEfMB0GA1UEAxMWcGVlcjAub3Jn
MS5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABAmj4G9864Au
qwUIPC79Knx30qdzTasb7FltKuaA+/TRlaPrxCLshIv10O1bMLeC64LHCcQslYIP
bGff1FvmAeqjTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1Ud
IwQkMCKAIOSIMM5bV/hv41orevOUfvHzZxFqYi6mdo0kFoeM5gX6MAoGCCqGSM49
BAMCA0gAMEUCIQDH4ULVm1YTi43cmbHf/8wXLQ/F02gYkmI16o6KXbABNAIgYHPt
/65Gb4NDnC5pwuS1URyi7CXv0dywArWinizUkEA=
-----END CERTIFICATE-----`

	cert, err := NewCertFromPem(CertTypeX509, []byte(certPem))
	if err != nil {
		t.Error(err)
	}
	xcert, err := cert.ToX509()
	if err != nil {
		t.Error(err)
	}
	ncert, _ := NewCertFromX509Cert(xcert)
	pem, err := ncert.ToPem()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(pem))

}
