multipass launch --cpus 2
multipass launch --cpus 2

for vm in `multipass list | awk '{print $1}'`
do
    multipass delete -p $vm
done