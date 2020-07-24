package gin

//是path.Clean的url版本，他返回不包含.和..元素的url路径
//将重复执行以下规则，直到无法进行下一步为止
//1. 用单个斜杆替换多个斜杆
//2. 在当前目录消除.元素
//3. 消除父目录中的..元素和前面的非..元素
//4. 消除以..开头的根目录，也就是说在根目录用/替换/..
func cleanPath(p string) string {
	//const stackBufSize = 128
	//
	//if p == "" {
	//	return "/"
	//}
	//
	////在堆栈上合理大小的缓冲区可避免在通常情况下进行分配
	////如果需要一个更大的内存，则他会自动分配
	//buf := make([]byte, 0, stackBufSize)
	//
	//n := len(p)
	//
	////r:下一个需要读取的字节位置
	////w:下一个需要写入的字节位置
	//r := 1
	//w := 1
	//
	//if p[0] != '/' {
	//	r = 0
	//	if n+1 > stackBufSize {
	//		buf = make([]byte, n+1)
	//	} else {
	//		buf = buf[:n+1]
	//	}
	//	buf[0] = '/'
	//}
	//
	//trailing := n > 1 && p[n-1] == '/'
	return ""
}
