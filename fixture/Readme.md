#### How to test

In order to test, you have to install `multipass`.

Replace the public key with your own `id_rsa.pub` in this [multipass.sh](https://github.com/debarshibasak/go-kubeadmclient/blob/master/fixture/multipass.sh)

```
chmod 755 multipass.sh
./multipass.sh 
```

This will create all the 3 multipass instances. 
One can be used for master, two of them can be used as worker nodes.

```.env
multipass list
```

This will show the details of the instances. Use the IP and the user name `ubuntu` in the testing code.

