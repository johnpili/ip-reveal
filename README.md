# Develop a whatismyipaddress.com clone website using Golang

Do you want to build a clone website similar to [whatismyipaddress.com](https://whatismyipaddress.com)? It is actually easy to develop. I was working on an automated DNS client that will check my public IP address and I decided to build this tool. Perhaps, somebody might need this as well in the future. I already made the completed tool available online at [ip.johnpili.com](https://ip.johnpili.com)
It works by reading the HTTP header request which contains information such as IP Address, User-Agent, Scheme, etc. If you are using a reverse proxy like Cloudflare, you can extract IP information from header keys ***Cf-Connecting-Ip*** or ***X-Real-Ip***
###Sample HTTP header
<pre>
map[
    Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8]
    Accept-Encoding:[gzip] 
    Accept-Language:[en-US,en;q=0.5] 
    Cache-Control:[no-cache]
    Cdn-Loop:[cloudflare] 
    Cf-Connecting-Ip:[193.169.145.66] 
    Cf-Ipcountry:[T1] 
    Cf-Visitor:[{"scheme":"https"}] 
    Connection:[upgrade] 
    Dnt:[1] Pragma:[no-cache] 
    Upgrade-Insecure-Requests:[1] 
    User-Agent:[Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0] 
    X-Forwarded-For:[193.169.145.66, 193.169.145.66] 
    X-Forwarded-Proto:[https] 
    X-Real-Ip:[193.169.145.66]
]
</pre>

### Snippet extracting IP address from header
<pre>
func (z *PageController) getIPDetails(r *http.Request) models.IPInfo {
	ip := ""
	if len(z.Configuration.Extraction.HeaderKey) > 0 {
		ip = r.Header.Get(z.Configuration.Extraction.HeaderKey) // Extract IP from header because we are using reverse proxy example X-Real-Ip
	}

	if len(ip) == 0 { // Fallback
		ip = extractIPAddress(r.RemoteAddr)
	}

	ipInfo := models.IPInfo{
		IP:        ip,
		UserAgent: r.Header.Get("User-Agent"),
	}
	return ipInfo
}
</pre>
You can checkout my blog post about this at [https://johnpili.com](https://johnpili.com/develop-a-whatismyipaddress-com-clone-website-using-golang/)

### Screenshot
![IP Echo - ip.johnpili.com](https://johnpili.com/wp-content/uploads/2020/01/Screen-Shot-2020-01-10-at-11.27.25-PM.png)