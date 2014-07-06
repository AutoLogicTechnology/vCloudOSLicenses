
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"
)

type ResourceEntity struct {
    XMLName string `xml:"ResourceEntity"`

    Name    string `xml:"name,attr"`
    Type    string `xml:"type,attr"`
    Href    string `xml:"href,attr"` 
}

type VDC struct {
    XMLName string `xml:"Vdc"`

    Status  string `xml:"status,attr"`
    Name    string `xml:"name,attr"`
    Id      string `xml:"id,attr"`
    Type    string `xml:"type,attr"`
    Href    string `xml:"href,attr"`

    ResourceEntities []*ResourceEntity `xml:"ResourceEntities"`
}

func (v *VDC) Get (session *VCloudSession, org *Organisation) {
    for link := range org.Links {
        if link.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            r := session.Get(link.Href)
            defer r.Body.Close()

            _ = xml.NewDecoder(r.Body).Decode(v)

            for k, val := range v.ResourceEntities {
                u, _ := url.Parse(val.Href)
                v.ResourceEntities[k].Href = u.Path
            }

            break
        }
    }    
}