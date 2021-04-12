package v1

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/tinywell/baas/internal/module"
	"github.com/tinywell/baas/internal/service/v1/metadata/common"
	"github.com/tinywell/baas/pkg/runtime"
	"github.com/tinywell/baas/pkg/runtime/docker"
)

func TestService_RunPeer(t *testing.T) {
	type fields struct {
		runner      runtime.ServiceRunner
		runtimeType int
	}
	type args struct {
		ctx   context.Context
		peers []*common.PeerData
	}
	runner, err := docker.NewClient()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "testrunpeer",
			fields: fields{
				runner:      runner,
				runtimeType: 0,
			},
			args: args{
				ctx:   context.Background(),
				peers: []*common.PeerData{peerData()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				runner:      tt.fields.runner,
				runtimeType: tt.fields.runtimeType,
			}
			if err := s.RunPeer(tt.args.ctx, tt.args.peers); (err != nil) != tt.wantErr {
				t.Errorf("Service.RunPeer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func peerData() *common.PeerData {
	dc := &module.DataCenterDocker{Name: "北京"}
	dcraw, err := json.Marshal(dc)
	if err != nil {
		//TODO:
	}

	vmsvc := &module.VMService{
		BaaSData: module.BaaSData{
			ID:        0,
			TenantID:  0,
			NetworkID: 0,
		},
		MSPID:    "Org1MSP",
		Name:     "peer0.org1.example.com",
		Runtime:  0,
		LinkType: 0,
		LinkID:   0,
		CFGRaw:   nil,
		CFG: map[string]interface{}{
			"": nil,
		},
		Status:     0,
		DataCenter: "BJ",
		DCMetadata: dcraw,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Creator:    "test",
	}
	TLSCert := `-----BEGIN CERTIFICATE-----
MIICcTCCAhigAwIBAgIQUR6GTlcmj807OtUhrHKvCTAKBggqhkjOPQQDAjB2MQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz
Y2Eub3JnMS5leGFtcGxlLmNvbTAeFw0yMTA0MDgwNzQxMDBaFw0zMTA0MDYwNzQx
MDBaMFsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH
Ew1TYW4gRnJhbmNpc2NvMR8wHQYDVQQDExZwZWVyMC5vcmcxLmV4YW1wbGUuY29t
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAExyv4Fr8P5QfPyb8/uKLzBFlWyLnQ
dmoZiHLZfsA30iBuutUaRfwhwFavBwSbYfhh2QE5WBTanXOHkr6YWI6Dm6OBojCB
nzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
MAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgVDMloq/C2xQC1NJWWqxuvfqsQYUT
uOjiXkayFYEC7r0wMwYDVR0RBCwwKoIWcGVlcjAub3JnMS5leGFtcGxlLmNvbYIF
cGVlcjCCCWxvY2FsaG9zdDAKBggqhkjOPQQDAgNHADBEAiA7+X7aEOYd6OiAiqgZ
HMAm4WrDwYtT9Gc8+XwbnxeZlwIgNzNRftZqHeqUOjGLuuWoXdbrzC+wrFGBerdT
iFClto8=
-----END CERTIFICATE-----`
	TLSKey := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgH0XcwhXf7JCaLlIf
J7COLmQi0w1O2FO2uqFJFfNE0oyhRANCAATHK/gWvw/lB8/Jvz+4ovMEWVbIudB2
ahmIctl+wDfSIG661RpF/CHAVq8HBJth+GHZATlYFNqdc4eSvphYjoOb
-----END PRIVATE KEY-----`
	MSPCert := `-----BEGIN CERTIFICATE-----
MIICKDCCAc6gAwIBAgIQb+YDQS5mh/JZKV+d3xcz7zAKBggqhkjOPQQDAjBzMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
b3JnMS5leGFtcGxlLmNvbTAeFw0yMTA0MDgwNzQxMDBaFw0zMTA0MDYwNzQxMDBa
MGoxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
YW4gRnJhbmNpc2NvMQ0wCwYDVQQLEwRwZWVyMR8wHQYDVQQDExZwZWVyMC5vcmcx
LmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECSp3GGIJsakr
MBrzQPUh2VWVYca343yOCGe7Iv3Vq8qKXJri6nOlM/DCjT46VDE22YOQInRVFCrT
DUvqk7dCFaNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0j
BCQwIoAgKAOvhoN7WsTNhUpOMALwR8s/I4jpDGNsTXFjVpKkJfkwCgYIKoZIzj0E
AwIDSAAwRQIhAJXoepzRZv3Ua9dvp2fsw/PDSALZWraI4jDaHaOI3pf6AiAGxwe0
ek90raEiiApoMho0t8SJHhz6YvxbCySixJkZJw==
-----END CERTIFICATE-----`
	MSPKey := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgP+g8mwsb/Wq+a1Iu
oRHzZiscS05CyTbreENq1prW4aqhRANCAAQJKncYYgmxqSswGvNA9SHZVZVhxrfj
fI4IZ7si/dWryopcmuLqc6Uz8MKNPjpUMTbZg5AidFUUKtMNS+qTt0IV
-----END PRIVATE KEY-----`
	OUConfig := `NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.org1.example.com-cert.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca.org1.example.com-cert.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/ca.org1.example.com-cert.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/ca.org1.example.com-cert.pem
    OrganizationalUnitIdentifier: orderer
`
	peer := &module.Peer{
		BaaSData: module.BaaSData{
			ID:        0,
			TenantID:  0,
			NetworkID: 0,
		},
		HFNode: module.HFNode{
			MSPID:    "Org1MSP",
			MSPKey:   MSPKey,
			MSPCert:  MSPCert,
			TLSKey:   TLSKey,
			TLSCert:  TLSCert,
			OUConfig: OUConfig,
		},
		Name:       "peer0.org1.example.com",
		DomainName: "peer0.org1.example.com",
		Endpoint:   "peer0.org1.example.com:7051",
		Port:       7051,
		EXTPort:    7051,
		Image:      "hyperledger/fabric-peer:2.2.2",
		StateDB:    "levelDB",
	}
	CACert := `-----BEGIN CERTIFICATE-----
MIICUjCCAfigAwIBAgIRAI9kcuRvw8gDfaoGAjfwfWEwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzEuZXhhbXBsZS5jb20wHhcNMjEwNDA4MDc0MTAwWhcNMzEwNDA2MDc0MTAw
WjBzMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UE
AxMTY2Eub3JnMS5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA
BE5xGqepmgHG9dT1QLD7uOMq40o00OW566jKmRwWdpF15swA2aAmWoOcxzfkjXSV
4QHzper5g86TmJWP7ErBBrWjbTBrMA4GA1UdDwEB/wQEAwIBpjAdBgNVHSUEFjAU
BggrBgEFBQcDAgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zApBgNVHQ4EIgQg
KAOvhoN7WsTNhUpOMALwR8s/I4jpDGNsTXFjVpKkJfkwCgYIKoZIzj0EAwIDSAAw
RQIhAOV8QdfxJDS90NG/ZjxHrh3F+XaJuLB9MBoEgqW8i1TiAiA2WdWhsQWmx88D
WIfjDTLGb1pY5DXRHry5DCki7KbTxg==
-----END CERTIFICATE-----`
	CAKey := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgrML+xaAKDVEx1R8/
JVl47evKuWtKtMHc/1zK0+LiVW2hRANCAAROcRqnqZoBxvXU9UCw+7jjKuNKNNDl
ueuoypkcFnaRdebMANmgJlqDnMc35I10leEB86Xq+YPOk5iVj+xKwQa1
-----END PRIVATE KEY-----`
	TLSCACert := `-----BEGIN CERTIFICATE-----
MIICVzCCAf2gAwIBAgIQaQHOKj7XYaG9bEv2AUrEMTAKBggqhkjOPQQDAjB2MQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz
Y2Eub3JnMS5leGFtcGxlLmNvbTAeFw0yMTA0MDgwNzQxMDBaFw0zMTA0MDYwNzQx
MDBaMHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH
Ew1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcxLmV4YW1wbGUuY29tMR8wHQYD
VQQDExZ0bHNjYS5vcmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEF12AB3qYeqgyNoqQdWfBsRv4Sx4jFM23wLAcE56wbT0DeHZ23dTjF4tZ
5GInTRAKpdpSEnXsY+i4FcXeuMBsdKNtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1Ud
JQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1Ud
DgQiBCBUMyWir8LbFALU0lZarG69+qxBhRO46OJeRrIVgQLuvTAKBggqhkjOPQQD
AgNIADBFAiEA99PnyaTcNoy9sYskktYtScj3SL/qb1Ccq5FdyfPr2egCIC4UoCTt
WgwsJqJGC4m5vgfAgQsFb3r+bDSSZHkyFmus
-----END CERTIFICATE-----`
	TLSCAKey := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg07wx8r/XmmOqS1s+
2Opjx8ndgP0LAt4LzkkY5e1QsSihRANCAAQXXYAHeph6qDI2ipB1Z8GxG/hLHiMU
zbfAsBwTnrBtPQN4dnbd1OMXi1nkYidNEAql2lISdexj6LgVxd64wGx0
-----END PRIVATE KEY-----`
	Org := &module.FOrganization{
		BaaSData: module.BaaSData{
			ID:        0,
			TenantID:  0,
			NetworkID: 0,
		},
		Name:      "org1",
		MSPID:     "Org1MSP",
		CACert:    CACert,
		CAKey:     CAKey,
		TLSCACert: TLSCACert,
		TLSCAKey:  TLSCAKey,
		AdminCert: "",
		Domian:    "org1.example.com",
	}
	return &common.PeerData{
		Service:     vmsvc,
		Extra:       peer,
		Org:         Org,
		NetworkName: "testnet",
	}
}