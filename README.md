# NUVI


1) Install GO (golang) from https://golang.org/doc/install

2) Install and start redis from https://redis.io/download

3) Run console and make command
```
git clone https://github.com/iostrovok/NUVI.git
cd ./NUVI/
make test
make install
```

4) Start download:
```
# See start option
./bin/index -help

# Download 3 files
./bin/index -files=3 -goes=3 -d=1 -old=false -port=6379 -host=localhost -url="http://bitly.com/nuvi-plz"

# Download all files 
./bin/index -goes=10 -d=1 -old=false -port=6379 -host=localhost -url="http://bitly.com/nuvi-plz"
```
