
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"
)

type Link struct {
    XMLName string `xml:"Link"`

    Name    string `xml:"name,attr"`
    Href    string `xml:"href,attr"`
    Type    string `xml:"type,attr"`
}

type VDCs struct {
    XMLName string `xml:"Org"`

    Records []*Link `xml:"Link"`
}

func (v *VDCs) GetAll (session *VCloudSession, org *OrgReference) {
    r := session.Get(org.Href)
    defer r.Body.Close()

    _ = xml.NewDecoder(r.Body).Decode(v)

    for k, val := range v.Records {
        if val.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            u, _ := url.Parse(val.Href)
            v.Records[k].Href = u.Path
        }
    }
}