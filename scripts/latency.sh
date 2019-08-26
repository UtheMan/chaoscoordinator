echo "Increasing latency by $latencyIncrease ms"
sudo tc qdisc add dev eth0 root netem delay ${latencyIncrease}ms
sleep $duration
sudo tc qdisc del dev eth0 root

