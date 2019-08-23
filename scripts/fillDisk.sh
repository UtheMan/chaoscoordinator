echo "Filling Disk with $amount MB of random data for $duration seconds."

nohup dd if=/dev/urandom of=/root/burn bs=1G count="$amount" iflag=fullblock
sleep "$duration"
rm /root/burn