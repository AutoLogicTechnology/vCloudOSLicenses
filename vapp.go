
package vcloudoslicenses 

import (
    "log"
    // "io/ioutil"

    "encoding/xml"
    // "net/url"
)

type VmOS struct {
    XMLName     string `xml:"OperatingSystemSection"`

    Id          string `xml:"deployed,attr"`
    Type        string `xml:"deployed,attr"`
    Href        string `xml:"deployed,attr"`
    OSType      string `xml:"osType,attr"`
}

type VAppVm struct {
    XMLName     string `xml:"Vm"`

    Deployed    string `xml:"deployed,attr"`
    Status      string `xml:"deployed,attr"`
    Name        string `xml:"deployed,attr"`
    Id          string `xml:"deployed,attr"`
    Type        string `xml:"deployed,attr"`
    Href        string `xml:"deployed,attr"`

    OperatingSystemSection *VmOS `xml:"OperatingSystemSection"`
}

type VAppChildren struct {
    XMLName     string `xml:"Children"`
    VM          []*VAppVm `xml:"Vm"`
}

type VDCVApp struct {
    XMLName     string `xml:"VApp"`

    Deployed    string `xml:"deployed,attr"`
    Status      string `xml:"deployed,attr"`
    Name        string `xml:"deployed,attr"`
    Id          string `xml:"deployed,attr"`
    Type        string `xml:"deployed,attr"`
    Href        string `xml:"deployed,attr"`

    VMs         VAppChildren `xml:"Children"`
}

func (a *VDCVApp) Get (session *VCloudSession, vdc *VdcResourceEntity) {
    r := session.Get(vdc.Href)
    defer r.Body.Close()

    // raw, _ := ioutil.ReadAll(r.Body)
    // log.Printf("RAW VAPP BODY: %v", string(raw))

    _ = xml.NewDecoder(r.Body).Decode(a)

    log.Printf("vApp Children Selfie: %+v", a)

    // for k1, v1 := range a.Children.Child {
    //     u, _ := url.Parse(v1.Href)
    //     a.Children.Child[k1].Href = u.Path 

    //     u, _ = url.Parse(v1.OperatingSystemSection.Href)
    //     a.Children.Child[k1].OperatingSystemSection.Href = u.Path
    // }

    // log.Printf("vApp Selfie: %+v", a)
}