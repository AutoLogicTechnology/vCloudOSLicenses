
package main 

import (
    "log"
)

func main() {
    session := &vCloudSession{}
    session.Login(...)
    

    orgs := &Organisations{}
    vdcs := &VDCs{}
    vapps := &vApps{}

    log.Print("Getting all organisations")
    orgs.GetAll(session, "references", 500)

    for _, o := range orgs.Records {
        log.Printf("Organisation: %v\n", o.Name)

        log.Print("Getting all VDCs, if applicable")
        vdcs.GetAll(session, o)
        
        if len(vdcs.Records) <= 0 {
            log.Printf("Organisation %s has no VDCs.", o.Name)
            continue 
        }

        for _, v := range vdcs.Records {
            if v.Type == "application/vnd.vmware.vcloud.vdc+xml" {
                log.Printf("VDC: %v\n", v.Name)

                log.Printf("Getting all vApps")
                vapps.GetAll(session, v)
            
                for _, a := range vapps.Records {
                    log.Printf("vApp: %v\n", a.Name)
                }
            }
        }
    }
}