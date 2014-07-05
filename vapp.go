
package main 

import (
    "encoding/xml"
    "net/url"

    // "io"
    // "io/ioutil"
)

type VAppLinkRecord struct {
    Name    string `xml:"name,attr"`
    Href    string `xml:"href,attr"`
    Type    string `xml:"type,atrr"`
}

type vApps struct {
    Records []*VAppLinkRecord `xml:"ResourceEntity"`
}

func (a *vApps) GetAll (session *vCloudSession, vdc *VdcLinkRecord) {
    r := session.Get(vdc.Href)
    defer r.Close()

    _ = xml.NewDecoder(r).Decode(a)
    // xml_decoder.Decode(a)

    // Loop over URLs and reduce the HREFs to URIs
    // We don't need the whole URL
    for k, v := range a.Records {
        u, _ := url.Parse(v.Href)
        a.Records[k].Href = u.Path 
    }

    // r := session.Get(vdc.Href)
    // defer r.Close()

    // io.Copy(ioutil.Discard, r)
}