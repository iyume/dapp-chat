# mypy: ignore-errors
import os
import shutil
import subprocess
import tempfile

cfg1 = tempfile.NamedTemporaryFile("w")  # bootnode cfg
cfg2 = tempfile.NamedTemporaryFile("w")
cfg3 = tempfile.NamedTemporaryFile("w")

cfg_template = """
data_dir = {data_dir}

[http]
address = 127.0.0.1:{port}
token = 123

[backend]
private_key = {hexkey}
address = 127.0.0.1:{peer_port}
bootnodes = {bootnodes}
"""

data_dirs = [f"chatdata-test-{i}" for i in range(3)]

hexkeys = [
    "3006931cdfa95453be77a0f9c0b46e6444b98ef9088608f3c4ed91e89baa00f9",
    "94fa73a584836623b488b0db234f5dee3147d85f04e383fdb8a474109c0f68f2",
    "d2f400a1d8dd2c8ff84f21826efab69ecca05e797b6605429769d87909d384fd",
]

ports = [44445, 44446, 44447]

peer_ports = [44400, 0, 0]

bootnode = "enode://075c4cdca8604e3e5a5a67122615714a426435a11a3874b16d66fac604b08cc0ebc55c9019efe8ca464bf7c9e5bfc25ee8a4736e24a26a4e8adaadce0d45a7a5@127.0.0.1:44400"

for i, c in enumerate((cfg1, cfg2, cfg3)):
    c.write(
        cfg_template.format(
            data_dir=data_dirs[i],
            port=ports[i],
            hexkey=hexkeys[i],
            peer_port=peer_ports[i],
            bootnodes=bootnode if i != 0 else "",
        )
    )
    c.flush()


# Run one bootnode and two regular nodes

p2pchat_executable = "p2pchat-bin"
os.system(f"go build -o {p2pchat_executable} ./p2pchat")


def start_server(id: int, config: str) -> subprocess.Popen:
    logfile = open(f"test-backend-{id}.log", "w")
    return subprocess.Popen(
        [f"./{p2pchat_executable}", "--config", config],
        stdout=logfile,
        stderr=logfile,
    )


p1 = start_server(0, cfg1.name)
p2 = start_server(1, cfg2.name)
p3 = start_server(2, cfg3.name)
ps = (p1, p2, p3)

print("#1", "102e0de7d9586b40990d986e3c5baee68678a16b2d90af3a086fb8f048594541")
print("#2", "a7c32ffba6e9449229886522f8e37d5684fc95e07ad92d05dd3577922dcd0321")
print("#3", "436e3172a545d485587c0e75b18b1e0759d9154396a15511a34d23200c4d2c89")


import os
import time

# kill p3 after connected for inactive node test
time.sleep(10)
p3.terminate()
retcode = p3.wait()
print("retcode", retcode)


try:
    for p in ps:
        p.wait()
except KeyboardInterrupt:
    for p in ps:
        p.terminate()
    for p in ps:
        if (retcode := p.wait()) and (retcode not in (0, 1)):
            print(f"{p.args} was not properly closed, retcode {retcode}")
finally:
    for d in data_dirs:
        shutil.rmtree(d, ignore_errors=True)
    os.remove(p2pchat_executable)
