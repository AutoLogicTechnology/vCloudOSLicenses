
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"
)

type VAppLinkRecord struct {
    Name    string `xml:"name,attr"`
    Href    string `xml:"href,attr"`
    Type    string `xml:"type,attr"`
}

type Resources struct {
    Entities []*VAppLinkRecord `xml:"ResourceEntity"`
}

type vApps struct {
    Records Resources `xml:"ResourceEntities"`
}

func (a *vApps) GetAll (session *VCloudSession, vdc *VdcLinkRecord) {
    r := session.Get(vdc.Href)
    defer r.Close()

    _ = xml.NewDecoder(r).Decode(a)
    for k1, v1 := range a.Records.Entities {
        u, _ := url.Parse(v1.Href)
        a.Records.Entities[k1].Href = u.Path 
    }
}