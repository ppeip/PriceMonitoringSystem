import random
import requests
lastpricelist = []
openpricelist = []
maxpricelist = []
minpricelist = []
limitlist= []
totallist = []
num = random.randint(8000,9500)
openpricelist.append(num)
for i in range(14):
    num = random.randint(8000,9500)
    total = random.randint(120000,130000)
    temp1 =random.randint(15,150)
    temp2 =random.randint(15,150)
    openpricelist.append(num)
    maxpricelist.append(num+temp1)
    minpricelist.append(num-temp2)
    lastprice = random.randint(num-temp2,num+temp1)
    lastpricelist.append(lastprice)
    totallist.append(total)
    #if i==0:
    #    limitlist.append(0)
    #else :
    limitlist.append((num-openpricelist[i-1])*100/num)
print(lastpricelist)
print(openpricelist)
print(limitlist)
print(maxpricelist)
print(minpricelist)
print(totallist)

for i in range(14):
    requests.post(url="http://localhost:8080/api/records",json=
        {
        "variety": "Ru",
        "cnname": "钌",
        "latestpri": lastpricelist[i],
        "openpri": openpricelist[1+i],
        "maxpri": maxpricelist[i],
        "minpri": minpricelist[i],
        "limit": f"{limitlist[i]:.2f}%",
        "yespri": openpricelist[i],
        "totalvol": totallist[i],
        "time": f"2022-11-{1+i}"
    }
    )