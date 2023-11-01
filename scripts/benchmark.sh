curl -X POST http://localhost:8082/reset
sleep 1
wrk -t4 -c100 -d30s --latency -s benchmark.lua http://localhost:8081
