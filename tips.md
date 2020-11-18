# Tips 助记   
exec.Command("sh") :新进程的初始命令，之后的命令都在sh环境下执行    
mount -t type FROM TO 指定文件类型，从FROM挂到TO  
SysProcAttr:设置attr属性 封装系统调用    
SysProcAttr.Credential:设置进程的执行用户    
stress命令用于模拟系统负载较高的场景   
/proc/N pid为N的进程信息  
/proc/N/stat 诸如此类 查看详细信息  
/proc/self 链接到当前进程   
os.exec会开始一个新进程 exec.Command("/proc/self/exe")会执行本身，覆盖那个新进程。   
cmd.Start()不阻塞 cmd.Run()阻塞 cmd.Wait()阻塞   
os.Stat 在指定路径上执行stat  
os.IsNotExist() 如果不存在该文件  
os.Getwd() get work directory  

