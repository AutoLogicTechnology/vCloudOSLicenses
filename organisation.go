
package vcloudoslicenses 

import (
    "log"
    "fmt"

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

func (o *Organisations) GetAll (session *VCloudSession, format string, max_page_size, max_pages int) {
    
    if max_page_size <= 0 {
        max_page_size = 1
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    for i := 1; i <= max_pages; i++ {
        uri := fmt.Sprintf("/api/query?type=organization&format=%v&pageSize=%v&page=%v", format, max_page_size, i)

        r := session.Get(uri)
        defer r.Body.Close()

        if r.StatusCode == 400 {
            break 
        }

        log.Printf("i = %v | uri = %s | status code = %v", i, uri, r.StatusCode)

        _ = xml.NewDecoder(r.Body).Decode(o)
     
        for k, v := range o.Records {
            u, _ := url.Parse(v.Href)
            o.Records[k].Href = u.Path 
        }
    }
}
