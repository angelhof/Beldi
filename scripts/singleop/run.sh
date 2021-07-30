#!/bin/bash
read -p "Choose mode (fast or full, default: full): " mode
mode=${mode:-"full"}
if [ "$mode" == "fast" ]; then
  duration=60
else
  duration=600
fi
echo "Cleaning logs at AWS"
python ./scripts/singleop/singleop.py --command clean
echo "Compiling"
make clean >/dev/null
make singleop >/dev/null
echo "Initializing Database"
go run ./internal/singleop/init/init.go >/dev/null
echo "Deploying"
sls deploy -c singleop.yml

## Get the endpoints from a file or ask them from the user and save them to the file if it doesnt exist
endpoints_file="singleop.endpoints"

if [ -f "$endpoints_file" ]; then
  echo "Endpoints acquired from: $endpoints_file"
  source "$endpoints_file"
else
  read -p "Please Input HTTP gateway url for beldi-dev-bsingleop: " bp
  read -p "Please Input HTTP gateway url for beldi-dev-singleop: " p
  read -p "Please Input HTTP gateway url for beldi-dev-tsingleop: " tp
  echo "bp=$bp"
  echo "p=$p"
  echo "tp=$tp"
fi

## TODO: Install wrk in submodules
wrk_bin=../wrk2/wrk
# wrk_bin=./tools/wrk
echo "Running baseline"
ENDPOINT="$bp" "$wrk_bin" -t1 -c1 -d"$duration"s -R1 -s ./benchmark/singleop/workload.lua --timeout 10s "$bp" >/dev/null
echo "Running beldi"
ENDPOINT="$p" "$wrk_bin" -t1 -c1 -d"$duration"s -R1 -s ./benchmark/singleop/workload.lua --timeout 10s "$p" >/dev/null
echo "Running beldi-txn"
ENDPOINT="$tp" "$wrk_bin" -t1 -c1 -d"$duration"s -R1 -s ./benchmark/singleop/workload.lua --timeout 10s "$tp" >/dev/null
echo "Collecting metrics"
python ./scripts/singleop/singleop.py --command run
echo "Cleanup"
go run ./internal/singleop/init/init.go clean >/dev/null
