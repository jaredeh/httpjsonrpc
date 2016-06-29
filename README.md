httpjsonrpc

Right now this is just a simple client for json-rpc over http.  I really really tried hard to find an "off-the-shelf" implementation that just worked.  What I found didn't actually connect to my server (targetd) maybe it's a flawed server?  I don't know.  It's not that complex to write it to satisfy my needs.  And it gets me practice with   So here we are.

example usage:

func main() {
	params := make(map[string]string)
	params["name"] = "vgdumb"
	params["path"] = "/dev/sda"

	jq := new(JsonrpcHttpClient)

	jq.Http.User = "admin"
	jq.Http.Password = "targetd"
	jq.Http.Host = "192.168.1.211"
	jq.Http.Port = "18700"
	result, err := jq.Execute("vg_create", params)

	if err != nil {
		log.Fatal("JsonrpcHttpClient error:", err)
	}
	spew.Dump(result)
}
