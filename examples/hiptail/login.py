#!/usr/local/bin/python
import mechanize
import cookielib
import sys
import requests
from lxml import html
import json

uri = str(sys.argv[1])
entity_id = uri.split("/")[4]
print entity_id
print uri

#flags
mClass=False
mDept=False
mSubClass=False

# Browser
br = mechanize.Browser()

# Cookie Jar
cj = cookielib.LWPCookieJar()
br.set_cookiejar(cj)

# Browser options
br.set_handle_equiv(True)
br.set_handle_gzip(True)
br.set_handle_redirect(True)
br.set_handle_referer(True)
br.set_handle_robots(False)

# Follows refresh 0 but not hangs on refresh > 0
br.set_handle_refresh(mechanize._http.HTTPRefreshProcessor(), max_time=1)

# Want debugging messages?
#br.set_debug_http(True)
#br.set_debug_redirects(True)
#br.set_debug_responses(True)

# User-Agent (this is cheating, ok?)
br.addheaders = [('User-agent', 'Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.1) Gecko/2008071615 Fedora/3.0.1-1.fc9 Firefox/3.0.1')]


# Open some site, let's pick a random one, the first that pops in mind:
br.open(uri)
br.select_form(nr=0)
br["username"]="pedrosimoes@parceiro.sonae.pt"
br["password"]="aTOEdm7W5AUW"
br.submit()
response = br.response()
page = response.read()
tree = html.fromstring(page)
stringasjson = tree.xpath('//*[@id="order_json"]')[0].text_content();
actualjson = json.loads(stringasjson)


for i in actualjson['stash']['checks']['order_items']['msg']:
    
    subject = ''.join([f for f in i if not f.isdigit()])
    if subject == 'Item  missing Dept':
            mDept=True
    if subject == 'Item  missing Class':
            mClass=True
    if subject == 'Item  missing SubClass':
            mSubClass=True

#curl -X PUT http://orlando-ws-prd.sonaesr.net/api/1/orlando/orders/30677222/status -d '{ "status" : "FIXING" }'

data = {
    "status":"Fixing"
}

if mDept==True and mClass==True and mSubClass==True:
    r = requests.put("http://orlando-ws-prd.sonaesr.net/api/1/orlando/orders/"+str(entity_id)+"/status", json={"status": "FIXING"})
    print(r.status_code, r.reason)
    #print actualjson['canonical'][0]['stash']
