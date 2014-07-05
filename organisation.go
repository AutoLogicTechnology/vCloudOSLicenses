
package main 

import (
    "fmt"
    // "log"

    "encoding/xml"
    "net/url"
)

type OrganisationReference struct {
    Name    string `xml:"name,attr"`
    Id      string `xml:"id,attr"`
    Href    string `xml:"href,attr"`
}

type Organisations struct {
    Records []*OrganisationReference `xml:"OrgReference"`
}

func (o *Organisations) GetAll (session *vCloudSession, format string, max int) {
    uri := fmt.Sprintf("/api/query?type=organization&format=%v&pageSize=%v", format, max)

    r := session.Get(uri)
    defer r.Close()

    _ = xml.NewDecoder(r).Decode(o)
 
    for k, v := range o.Records {
        u, _ := url.Parse(v.Href)
        o.Records[k].Href = u.Path 
    }
}
