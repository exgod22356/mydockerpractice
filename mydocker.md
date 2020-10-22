# Write a docker  
### Namespace  
Linux has 6 different namespace:  
**UTS Namespace** :UTS Namespace is mainly used to isolate "nodename" and "hostname". You can change your hostname after executed CLONE_NEWUTS.   
**IPC Namespace** :IPC Namespace is used to isolate "System V IPC" and "POSIX message queue", like shared memory, message queue or semaphore. To achieve this, you can call CLONE_NEWIPC.   
**PID Namespace** :PID Namespace is used to isolate PIDs. If we use "ps -ef" command in a container, we shall see the PID in the front-page is 1. But outside the container, It's totally different.  We can just use CLONE_NEWPID to achieve this.   
**Mount Namespace** :Mount Namespace is used to isolate "mount" and "umount" command. It's the first namespace that linux have ever achieved, so we need to call this function by "CLONE_NEWNS"(clone new namespace)  
**User Namespace** :User Namespace is used to isolate user id and user group id. You can give users root permission in the namespace to limit the permission. Call this fuction by "CLONE_NEWUSER".  
**Network Namespace** :Network Namespace is used to isolate network devices and IP ports. By using namespace, we can avoid conflicts of ports. So that we can build bridge on the host between different containers to easily achieve the communication.Call this function by "CLONE_NEWNET".    
  
### Cgroups   
Cgroups provide a method to limit sub-processes' resources. Cgroups have 3 components: cgroup, subsystem, hierarchy.   
**cgroup** is a system to manage process groups. A cgroup has a set of progresses, and you can associate them with different subsystem configurations.   
**subsystem** is a module that could control resources. You can "apt-get install cgroup-bin" and "lssubsys" to check the subsystem your kernal could support.    
**hierarchy** turns a string of cgroup into a tree structure. By this tree structure, Cgroups can be inherited.   




