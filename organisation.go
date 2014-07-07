
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"
)

type OrgLink struct {
    XMLName     string `xml:"Link"`

    Type        string `xml:"type,attr"` 
    Name        string `xml:"name,attr"`
    Href        string `xml:"href,attr"`
}

type Organisation struct {
    XMLName     string `xml:"Org"`

    Name        string `xml:"name,attr"`
    Id          string `xml:"id,attr"`
    Type        string `xml:"type,attr"`
    Href        string `xml:"href,attr"`

    Links       []*OrgLink `xml:"Link"`
}

func (o *Organisation) Get (session *VCloudSession, org_url string) {
    uri, _ := url.Parse(org_url)
    r, _ := session.Get(uri.Path)
    defer r.Body.Close()

    _ = xml.NewDecoder(r.Body).Decode(o)

    for k, v := range o.Links {
        u, _ := url.Parse(v.Href)
        o.Links[k].Href = u.Path
    }
}