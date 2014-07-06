
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"

    "log"
)

type VdcResourceEntity struct {
    XMLName string `xml:"ResourceEntity"`

    Name    string `xml:"name,attr"`
    Type    string `xml:"type,attr"`
    Href    string `xml:"href,attr"` 
}

type VdcResourceEntities struct {
    XMLName string `xml:"ResourceEntities"`
    ResourceEntity []*VdcResourceEntity `xml:"ResourceEntity"`
}

type VDC struct {
    XMLName string `xml:"Vdc"`

    Status  string `xml:"status,attr"`
    Name    string `xml:"name,attr"`
    Id      string `xml:"id,attr"`
    Type    string `xml:"type,attr"`
    Href    string `xml:"href,attr"`

    ResourceEntities VdcResourceEntities `xml:"ResourceEntities"`
}

func (v *VDC) Get (session *VCloudSession, org *OrgLink) {
    r := session.Get(org.Href)
    defer r.Body.Close()

    _ = xml.NewDecoder(r.Body).Decode(v)

    log.Printf("VDC Selfie: %+v", v)

    for k, val := range v.ResourceEntities {
        u, _ := url.Parse(val.Href)
        v.ResourceEntities[k].Href = u.Path

        log.Printf("VDC: %s | Href: %s", v.Name, v.ResourceEntities[k].Href)
    }
}