FROM golang:1.18-bullseye
EXPOSE 16657
EXPOSE 16656
EXPOSE 6060
EXPOSE 1316 
EXPOSE 8090 
EXPOSE 8091 
EXPOSE 8080 
EXPOSE 26656 
EXPOSE 26657 
EXPOSE 1317 
EXPOSE 9090 
EXPOSE 9091
EXPOSE 8081
ADD . /opt/neutron
RUN cd /opt/neutron && make install && PLATFORM=`uname -m` && \
    curl -L "https://github.com/informalsystems/ibc-rs/releases/download/v1.0.0/hermes-v1.0.0-${PLATFORM}-unknown-linux-gnu.tar.gz" > hermes.tar.gz && \
    mkdir -p $HOME/.hermes/bin && \
    tar -C $HOME/.hermes/bin/ -vxzf hermes.tar.gz && \
    rm -f hermes.tar.gz 
ENV PATH="/root/.hermes/bin:${PATH}"
WORKDIR /opt/neutron
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD \
    curl -f http://127.0.0.1:1317/blocks/1 >/dev/null 2>&1 || exit 1
CMD ./network/init.sh && \
	./network/start.sh && \
	./network/hermes/restore-keys.sh && \
	./network/hermes/create-conn.sh && \
    ./network/hermes/start.sh
