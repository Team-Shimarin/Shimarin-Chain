from time import sleep
import urllib.request
import urllib.parse
import random

def get_sample():
    account_list = ["shima1", "shima2", "shima3", "shima4"]
    result = random.sample(account_list, 2)
    return result

def make_request():
    ids = get_sample()
    url = "http://localhost:8081/api/v1/balance/remit"
    param = [
        ( "to", ids[0]),
        ( "from", ids[1]),
        ( "value", random.randrange(100, 1000))
    ]
    url += "?{0}".format( urllib.parse.urlencode( param))
    request = urllib.request.Request(url, method="POST")
    f = urllib.request.urlopen(request)
    print(f.read().decode('utf-8'))

if __name__ == "__main__":
    while(True):
        try:
            print("Connect")
            make_request()
            sleep(random.randrange(10, 15))
        except as e:
            print("Failed To Connect", e)
