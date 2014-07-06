
package vcloudoslicenses

import (
    "log"
    "fmt"

    "encoding/xml"
    "net/url"
    "errors"
)

type OrganisationReference struct {
    XMLName string `xml:"OrgReference"`
    Name    string `xml:"name,attr"`
    Id      string `xml:"id,attr"`
    Href    string `xml:"href,attr"`
}

type Organisations struct {
    XMLName     string  `xml:"OrgReferences"`
    Records     []*OrganisationReference `xml:"OrgReference"`
}

func FindOrganisations (session *VCloudSession, max_page_size, max_pages int) (Orgs *Organisations, err error) {

	Orgs = &Organisations{}

    if max_page_size <= 0 {
        max_page_size = 1
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    for i := 1; i <= max_pages; i++ {
    	o := &Organisations{}

        uri := fmt.Sprintf("/api/query?type=organization&format=references&pageSize=%v&page=%v", max_page_size, i)
        r := session.Get(uri)
        defer r.Body.Close()

        if r.StatusCode == 400 {
            break 
        }

        _ = xml.NewDecoder(r.Body).Decode(o)
     
        for k, v := range o.Records {
            u, _ := url.Parse(v.Href)
            o.Records[k].Href = u.Path 

            Orgs.Records = append(Orgs.Records, v)
        }

        log.Printf("i = %v | uri = %s | status code = %v | me = %+v", i, uri, r.StatusCode, o.Records)
    }

    if len(Organisations.Records) <= 0 {
    	return &Organisations{}, errors.New("No organisations returned.")
    } else {
    	return Orgs, nil
    }
}