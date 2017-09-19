### file_watcher

Recursively watches directories and calls `vagrant rsync` when a file change is detected. 

For larger projects, you might have to increase the open file limit:

```
$ echo kern.maxfiles=65536 | sudo tee -a /etc/sysctl.conf
$ echo kern.maxfilesperproc=65536 | sudo tee -a /etc/sysctl.conf
$ sudo sysctl -w kern.maxfiles=65536
$ sudo sysctl -w kern.maxfilesperproc=65536
$ ulimit -n 65536 65536    
$ echo 'ulimit -n 65536 65536' | tee -a ~/.bashrc

```
