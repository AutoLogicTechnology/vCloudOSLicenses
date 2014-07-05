
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"

    "log"
)

type VMRecordOSType struct {
    Name string `xml:"osType,attr"`

    Href string `xml:"vcloud:href,attr"`
    Type string `xml:"vcloud:type,attr"`
}

type VMRecord struct {
    OSType VMRecordOSType `xml:"OperatingSystemSection"`
}

type Child struct {
    Server []*VMRecord `xml:"Vm"`
}

type VMs struct {
    Records Child `xml:"Children"`
}

func (v *VMs) GetAll (session *VCloudSession, vapp *VAppLinkRecord) {
    r := session.Get(vapp.Href)
    defer r.Body.Close()

    _ = xml.NewDecoder(r.Body).Decode(v)
    for _, v1 := range v.Records.Server {
        u, _ := url.Parse(v1.OSType.Href)
        v1.OSType.Href = u.Path 
    }

    log.Printf("vm = %+v", v.Records.Server)
}