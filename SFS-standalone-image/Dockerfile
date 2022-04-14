FROM ubuntu:20.04
RUN apt update -y && apt install schedtool software-properties-common -y && add-apt-repository ppa:longsleep/golang-backports && apt update -y && apt install golang-go python python3 -y
#RUN apt update -y && apt install schedtool golang-go python python3 -y
RUN mkdir /SFS-standalone
COPY ./SFS-standalone/* /SFS-standalone/
WORKDIR /SFS-standalone
RUN go build && mkdir /result
CMD ["./test.sh"]
