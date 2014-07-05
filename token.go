
package vcloudoslicenses 

import (
    "net/http"
    "crypto/tls"
    "log"
    "fmt"
    "encoding/base64"
)

type VCloudSession struct {
    Host            string 
    Username        string 
    Password        string 
    Context         string

    Transport       *http.Transport
    Token           string
    Accessible      bool
}

func (v *VCloudSession) Login () {
    credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s@%s:%s", v.Username, v.Context, v.Password)))

    v.Transport     = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
    
    request, _      := http.NewRequest("GET", fmt.Sprintf("%s/api/sessions", v.Host), nil)
    request.Header.Add("Authorization", fmt.Sprintf("Basic %s", credentials))
    request.Header.Add("Accept", "application/*+xml;version=5.1")

    client          := &http.Client{Transport: v.Transport}
    response, err   := client.Do(request)

    defer response.Body.Close()

    if err != nil {
        log.Fatal(err)
    }

    v.Token         = response.Header.Get("x-vcloud-authorization")
    v.Accessible    = true
}

func (v *VCloudSession) Get (uri string) (body *http.Response) {
    var err error 
    var response *http.Response 

    if v.Accessible {
        request, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", v.Host, uri), nil)
        request.Header.Add("x-vcloud-authorization", v.Token)
        request.Header.Add("Accept", "application/*+xml;version=5.1")

        client := &http.Client{Transport: v.Transport}
        response, err = client.Do(request)

        if err != nil {
            log.Fatal(fmt.Sprintf("Call to %s failed: %v", uri, err))
        }

    } else {
        log.Fatal("NewRequest() called, but no accessible session available.")
    }

    return response
}
