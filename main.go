
package main 

import "log"

func main() {
    session := &vCloudSession{}
    // session.Login(...)
    

    log.Printf("Total Report Count: %d", len(Report(session)))
}