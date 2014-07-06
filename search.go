
package vcloudoslicenses

import (
    "fmt"

    "encoding/xml"
    "net/url"
    "errors"
)

type OrgReference struct {
    XMLName string `xml:"OrgReference"`
    Name    string `xml:"name,attr"`
    Id      string `xml:"id,attr"`
    Href    string `xml:"href,attr"`
}

type OrgReferences struct {
    XMLName     string  `xml:"OrgReferences"`
    Records     []*OrgReference `xml:"OrgReference"`
}

func (v *VCloudSession) FindVApps (max_page_size, max_pages int) (VApps []*AdminVAppRecord, err error) {

    if max_page_size <= 0 {
        max_page_size = 1
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    for i := 1; i <= max_pages; i++ {
        vapp := &VAppQueryResultsRecords{}
        uri := fmt.Sprintf("/api/query?type=adminVApp&pageSize=%v&page=%v", max_page_size, i)

        r := v.Get(uri)
        defer r.Body.Close()

        if r.StatusCode == 400 {
            break
        }

        _ = xml.NewDecoder(r.Body).Decode(vapp)
        for k,v := range vapp.Records {
            u, _ := url.Parse(v.Href)
            vapp.Records[k].Href = u.Path

            VApps = append(VApps, v)
            v.Counters.VApps++
        }
    }

    return VApps, nil 
}

func (v *VCloudSession) FindOrganisations (max_page_size, max_pages int) (Orgs []*Organisation, err error) {

    if max_page_size <= 0 {
        max_page_size = 1
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    for i := 1; i <= max_pages; i++ {
    	page := &OrgReferences{}
        uri := fmt.Sprintf("/api/query?type=organization&format=references&pageSize=%v&page=%v", max_page_size, i)
        r := session.Get(uri)
        defer r.Body.Close()

        if r.StatusCode == 400 {
            break 
        }

        _ = xml.NewDecoder(r.Body).Decode(page)
     
        for _, v := range page.Records {
            u, _ := url.Parse(v.Href)
            r := session.Get(u.Path)
            defer r.Body.Close()

            if r.StatusCode != 200 {
            	continue 
            }

            new_org := &Organisation{}
            _ = xml.NewDecoder(r.Body).Decode(new_org)

            for k, val := range new_org.Links {
                u, _ := url.Parse(val.Href)
                new_org.Links[k].Href = u.Path
            }

            Orgs = append(Orgs, new_org)
        }
    }

    if len(Orgs) <= 0 {
    	return Orgs, errors.New("No organisations returned.")
    } else {
    	return Orgs, nil
    }
}