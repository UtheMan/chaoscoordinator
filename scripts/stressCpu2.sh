cpus=$(cat /proc/cpuinfo | awk "/^processor/{print $3}" | wc -l)
pids=""
echo "Stressing $cpus CPUs for $duration seconds."
#trap 'for p in $pids; do kill $p; done' 0
trap 'pkill -P $$' 0
for i in `seq 1 $cpus`
do
    sha1sum /dev/zero | sha1sum /dev/zero | sha1sum /dev/zero | sha1sum /dev/zero &
#    pids="$pids $!";
done
sleep $duration