f = open("default-ol/config.json","r")
strings = f.read()

for worker in range(1,64):
    port = 5003+worker
    strings_cp = strings.replace("5003",str(port))
    strings_cp = strings_cp.replace("default-ol","worker"+str(worker))
    f2 = open("{}/config.json".format("worker"+str(worker)),"w")
    f2.write(strings_cp)
    f2.close()
f.close()
