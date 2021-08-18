#!/bin/bash
read -p "Choose mode (fast or full, default: full): " mode
mode=${mode:-"full"}
if [ "$mode" == "fast" ]; then
  duration=60
else
  duration=600
fi
echo "Cleaning logs at AWS"
python ./scripts/append/append.py --command clean
echo "Compiling"
make clean >/dev/null
make append >/dev/null

## TODO: Refactor so that the same source nop is used for all
## TODO: Add the main.go files for the append app
echo "Initializing Database"
go run ./internal/append/init/init.go >/dev/null
echo "Deploying"
sls deploy -c append.yml

## Get the endpoints from a file or ask them from the user and save them to the file if it doesnt exist
endpoints_file="append.endpoints"

if [ -f "$endpoints_file" ]; then
  echo "Endpoints acquired from: $endpoints_file"
  source "$endpoints_file"
else
  read -p "Please Input HTTP gateway url for beldi-dev-bappend: " bp
  read -p "Please Input HTTP gateway url for beldi-dev-append: " p
  read -p "Please Input HTTP gateway url for beldi-dev-tappend: " tp
  echo "bp=$bp" | tee -a "$endpoints_file"
  echo "p=$p" | tee -a "$endpoints_file"
  echo "tp=$tp" | tee -a "$endpoints_file"
fi

## TODO: Install wrk in submodules
wrk_bin=../wrk2/wrk
# wrk_bin=./tools/wrk
echo "Running baseline"
ENDPOINT="$bp" "$wrk_bin" -t1 -c1 -d"$duration"s -R1 -s ./benchmark/append/workload.lua --timeout 10s "$bp" >/dev/null
echo "Running beldi"
ENDPOINT="$p" "$wrk_bin" -t1 -c1 -d"$duration"s -R1 -s ./benchmark/append/workload.lua --timeout 10s "$p" >/dev/null
echo "Running beldi-txn"
ENDPOINT="$tp" "$wrk_bin" -t1 -c1 -d"$duration"s -R1 -s ./benchmark/append/workload.lua --timeout 10s "$tp" >/dev/null
echo "Collecting metrics"
python ./scripts/append/append.py --command run
echo "Cleanup"
go run ./internal/append/init/init.go clean >/dev/null
