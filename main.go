
package main 

import (
    "log"
)

func main() {
    log.SetFlags(log.Lmicroseconds)

    session := &vCloudSession{}
    session.Login(...)

    log.Printf("Total Report Count: %d", len(Report(session)))
}