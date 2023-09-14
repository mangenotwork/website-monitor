package dao

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mangenotwork/common/utils"
	"net/http"
	"strconv"
	"strings"
	"website-monitor/master/entity"
)

// GetCertificateInfo 获取SSL证书信息
func GetCertificateInfo(caseUrl string) (entity.SSLCertificateInfo, bool) {
	caseUrl = urlStr(caseUrl)
	var info = entity.SSLCertificateInfo{
		Url: caseUrl,
	}
	var cert *x509.Certificate
	var err error
	client := http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if len(rawCerts) < 1 {
					return nil
				}
				cert, err = x509.ParseCertificate(rawCerts[0])
				return err
			},
		},
	}
	_, err = client.Get(caseUrl)
	if err != nil || cert == nil {
		return info, false
	}
	info.NotBefore = cert.NotBefore.Unix()
	info.NotAfter = cert.NotAfter.Unix()
	info.EffectiveTime = fmt.Sprintf("%s 到 %s", utils.Timestamp2Date(info.NotBefore), utils.Timestamp2Date(info.NotAfter))
	info.DNSName = strings.Join(cert.DNSNames, ";")
	info.OCSPServer = strings.Join(cert.OCSPServer, ";")
	info.CRLDistributionPoints = strings.Join(cert.CRLDistributionPoints, ";")
	info.Issuer = cert.Issuer.String()
	info.IssuingCertificateURL = strings.Join(cert.IssuingCertificateURL, ";")
	info.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()
	info.Subject = cert.Subject.String()
	info.Version = strconv.Itoa(cert.Version)
	info.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	return info, true
}
