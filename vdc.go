
package main 

import (
    "encoding/xml"
    "net/url"

    // "io"
    // "io/ioutil"
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
    r := session.Get(org.Href)
    defer r.Close()

    _ = xml.NewDecoder(r).Decode(v)
    // xml_decoder.Decode(v)

    // Loop over URLs and reduce the HREFs to URIs
    // We don't need the whole URL
    for k, val := range v.Records {
        if val.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            u, _ := url.Parse(val.Href)
            v.Records[k].Href = u.Path
        }
    }

    // r := session.Get(org.Href)
    // defer r.Close()

    // io.Copy(ioutil.Discard, r)
}