# Write a docker  
### Namespace  
Linux has 6 different namespace:  
**UTS Namespace** :UTS Namespace is mainly used to isolate "nodename" and "hostname". You can change your hostname after executed CLONE_NEWUTS.   
**IPC Namespace** :IPC Namespace is used to isolate "System V IPC" and "POSIX message queue", like shared memory, message queue or semaphore. To achieve this, you can call CLONE_NEWIPC.   
**PID Namespace** :PID Namespace is used to isolate PIDs. If we use "ps -ef" command in a container, we shall see the PID in the front-page is 1. But outside the container, It's totally different.  We can just use CLONE_NEWPID to achieve this.   
**Mount Namespace** :Mount Namespace is used to isolate "mount" and "umount" command. It's the first namespace that linux have ever achieved, so we need to call this function by "CLONE_NEWNS"(clone new namespace)  

