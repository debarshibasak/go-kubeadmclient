for vm in `multipass list | awk '{print $1}'`
do
    multipass delete -p $vm
done

multipass launch --cpus 2
multipass launch --cpus 2
multipass launch --cpus 2

for vm in `multipass list | awk '{print $1}'`
do
    multipass exec $vm -- sh -c 'echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCjM5hwX1cM4BP7ozU9vsekBe6DYR8Lc35YgLlU8ZKY0MAM8ABx3j41eB2Pr2zpoXJ2Q7N/QkMqfIIjTh90RPIkU3Nlpo4WnyQdPrP/YDU7//H7xpZZg814NKJQG+Qatd5j4HJklEdKtTy1PPnCwgSBhDRjH2WUxvQYvt+FRtPKtrKmjCtjJ/X1T4QMNzpO33K93WJxyjMtFK+oem83xUrcdyaIXxTWN/1iU6yoR6TAABS9kpKx3bdikQnu1kb1o+J8HyAPeU424XOILe5fdyYdhxPvNLvoP0EJ6xcg3f58wmZbdlSe8l9B1zRW0CwHZ5VN+Zy4GkNmFxeTMB9tJv9u9vGnwtL1mKPR6YysYT3rv3V38KzZVHU5u8MusteZ+POUL/QCIDfGrVXg+Sr7iDj5QqumxrxId7OSLg5bUhLEQ6GeK7AbC4SxbhP+93E2PPjG9iOgGcGEs0M30w4Rx0/PXGWd7wwrbC/dgeTMatiXBcvhS/DSMiiSv7aJMVeXgr5VA6pK5sM4FHGmYH0r9OPBxnkz0Ge3RMudt/t0aqcYGpG+l8+/btgqTZrr0RtUKgIzk+xmPLUH1usJVeTuIHcgrnb382OfqRusJbOJTn4nWicStAq5rnbfOgsc7CKVT9+zY1H7eNejCD4czhloUuwZ4IzYrtGHHOT5y3/D4Y37w== debarshri@gmail.com" >> /home/ubuntu/.ssh/authorized_keys'
done


multipass list
