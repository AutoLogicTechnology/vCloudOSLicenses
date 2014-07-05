
package vcloudoslicenses 

type vCloud struct {
    Hostname    string 
    Port        string 
    Username    string 
    Password    string 
    Context     string
}

type Configuration struct {
    vCloud          vCloud
    Session         *vCloudToken 

    Debugging       bool
}
