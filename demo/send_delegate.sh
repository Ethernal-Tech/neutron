#!/bin/bash

cd $(dirname "$0")

export OWNER=$(neutrond keys show rly1 -a --keyring-backend test --home ../data/test-1) && echo $OWNER;

../build/neutrond tx bank send cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs cosmos1t5rejj0jgw63fc2mjqxdqjh6cdhhwak6rrsx5sjuc6npkpqd553s50gqpl 10000stake --chain-id test-2 --home ../data/test-2 --node tcp://localhost:26657 --keyring-backend test -y

sleep 3;
OWNER="cosmos14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s4hmalr"

../build/neutrond tx interchaintxs submit-tx connection-0 $OWNER \
    ./test_tx_delegate.json --from demowallet1 --gas 10000000 --gas-adjustment 1.4 --gas-prices 0.5stake \
    --broadcast-mode block --chain-id test-1 --keyring-backend test --home ../data/test-1 --node tcp://127.0.0.1:16657