
package main 

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

func (v *VDCs) GetAll (session *vCloudSession, org *OrganisationReference) {
    xml_decoder := xml.NewDecoder(session.Get(org.Href))
    xml_decoder.Decode(v)

    // Loop over URLs and reduce the HREFs to URIs
    // We don't need the whole URL
    for k, val := range v.Records {
        if val.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            u, _ := url.Parse(val.Href)
            v.Records[k].Href = u.Path
        }
    }
}