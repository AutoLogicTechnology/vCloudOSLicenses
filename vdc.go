
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"
)

type VdcLinkRecord struct {
    Name    string `xml:"name,attr"`
    Href    string `xml:"href,attr"`
    Type    string `xml:"type,attr"`
}

type VDCs struct {
    Records []*VdcLinkRecord `xml:"Link"`
}

func (v *VDCs) GetAll (session *VCloudSession, org *OrganisationReference) {
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