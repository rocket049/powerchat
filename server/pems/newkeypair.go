package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

//证书模板，通过该模板默认设置一些证书需要的字段，比如序列号，组织信息，有效期等等
func CertTemplate() (*x509.Certificate, error) {
	//生成随机的序列号 (不同组织可以有不同的序列号生成方式)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}
	tmpl := x509.Certificate{SerialNumber: serialNumber,
		Subject:               pkix.Name{Organization: []string{hostname}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 20), //20年的有效期
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func createRootKey() (*rsa.PrivateKey, *x509.Certificate) {
	//生成一对新的公私钥
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating random key: %v", err)
	}
	rootCertTmpl, err := CertTemplate()
	if err != nil {
		log.Fatalf("creating cert template: %v", err)
	} //在模板的基础上增加一些新的证书信息
	rootCertTmpl.IsCA = true //是否是CA
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	return rootKey, rootCertTmpl
}

//CreateCert 参数：
//    template   证书申请者的证书模板
//    parent     父证书
//    pub        证书申请者的公钥
//    parentPriv 父证书对应的私钥

func CreateCert(template, parent *x509.Certificate,
	pub interface{}, parentPriv interface{}) (cert *x509.Certificate, certPEM []byte, err error) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	} // parse the resulting certificate so we can use it again cert,
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	} //将 certDER 用 pem 编码，生成 certPEM 证书
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

func GetParentCert(name string) (cert *x509.Certificate, privKey *rsa.PrivateKey) {
	certPem, err := ioutil.ReadFile(name + "-cert.pem")
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(certPem)
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("parseCertificate:%v\n", err)
	}
	keyPem, err := ioutil.ReadFile(name + "-key.pem")
	if err != nil {
		panic(err)
	}
	block, _ = pem.Decode(keyPem)
	privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("parsePrivateKey:%v\n", err)
	}
	return
}

func saveKeyPair(name string, prik *rsa.PrivateKey, savePub bool) error {
	var pubk rsa.PublicKey
	pubk.N = prik.N
	pubk.E = prik.E
	// save private key
	var blk pem.Block
	blk.Type = "RSA PRIVATE KEY"
	file1 := name + "-key.pem"
	fp, err := os.Create(file1)
	if err != nil {
		panic(err)
	}
	blk.Bytes = x509.MarshalPKCS1PrivateKey(prik)
	err = pem.Encode(fp, &blk)
	if err != nil {
		panic(err)
	}
	fp.Close()
	//save public key
	if savePub {
		blk.Type = "RSA PUBLIC KEY"
		file1 = name + "-pub.pem"
		fp, err = os.Create(file1)
		if err != nil {
			panic(err)
		}
		blk.Bytes = x509.MarshalPKCS1PublicKey(&pubk)
		err = pem.Encode(fp, &blk)
		if err != nil {
			panic(err)
		}
		fp.Close()
	}
	return nil
}

//CreateTlsKeyPair name:输出文件前缀；parent：签名用的根证书前缀
func CreateTlsKeyPair(name, parent string) {
	var rootCertPEM []byte
	var err error
	rootKey, rootCertTmpl := createRootKey()
	if len(parent) == 0 {
		_, rootCertPEM, err = CreateCert(rootCertTmpl, rootCertTmpl,
			&rootKey.PublicKey, rootKey)
	} else {
		cert, privKey := GetParentCert(parent)
		_, rootCertPEM, err = CreateCert(rootCertTmpl, cert,
			&rootKey.PublicKey, privKey)
	}
	//rootCert, rootCertPEM, err := CreateCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		log.Fatalf("error creating cert: %v", err)
	}
	fmt.Printf("%s\n", rootCertPEM)
	ioutil.WriteFile(name+"-cert.pem", []byte(rootCertPEM), 0644)
	//fmt.Printf("%#x\n", rootCert.Signature) // 证书的签名信息
	saveKeyPair(name, rootKey, false)
}

func main() {
	var name = flag.String("name", "new", "[-name ident] default: new")
	var parent = flag.String("parent", "",
		"[-parent parentKeyFilePrefix] default: null")
	flag.Parse()
	CreateTlsKeyPair(*name, *parent)
}
