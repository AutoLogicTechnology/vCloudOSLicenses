
package vcloudoslicenses 

type vCloud struct {
	Hostname 	string 
	Port 		string 
	Username 	string 
	Password 	string 
	Context 	string
}

type ElasticSearch struct {
	Hostname 	string 
	Port 		string 
	Index		string 
	DocType	 	string
}

type Configuration struct {
	vCloud 			vCloud
	ElasticSearch 	ElasticSearch
	Session 		*vCloudToken 

	Debugging 		bool
}
