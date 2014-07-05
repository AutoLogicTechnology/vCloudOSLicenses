
package main 

import (
    "log"

    "github.com/davecheney/profile"
)

func main() {
    profile_cfg := profile.Config{CPUProfile: true, MemProfile: true, BlockProfile: true, NoShutdownHook: false,}
    profile     := profile.Start(&profile_cfg)

    defer profile.Stop()

    session := &vCloudSession{}
    session.Login("https://vcloud", "...", "...", "...")

    orgs := &Organisations{}
    vdcs := &VDCs{}
    vapps := &vApps{}

    orgs.GetAll(session, "references", 500)
    for _, o := range orgs.Records {
        log.Printf("Organisation: %v\n", o.Name)

        vdcs.GetAll(session, o)
        
        if len(vdcs.Records) <= 0 {
            log.Printf("Organisation %s has no VDCs.", o.Name)
            continue 
        }

        for _, v := range vdcs.Records {
            if v.Type == "application/vnd.vmware.vcloud.vdc+xml" {
                log.Printf("VDC: %v\n", v.Name)

                vapps.GetAll(session, v)
            
                for _, a := range vapps.Records {
                    log.Printf("vApp: %v\n", a.Name)
                }
            }
        }
    }
}