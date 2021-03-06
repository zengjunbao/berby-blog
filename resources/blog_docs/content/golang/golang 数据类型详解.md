

[TOC]



## string

字符串本质是字符数组。

标准库`builtin` 对string的描述：

> string is the set of all strings of 8-bit bytes, conventionally but not necessarily representing UTF-8-encoded text. A string may be empty, but not nil. Values of string type are immutable.
>
> 字符串是字节的一个序列，约定但不必须是 UTF-8 编码的文本。字符串可以为空但不能是nil，其值不可变。



### 源码

Go 中字符串的源码定义在 `src/runtime/string.go`：

```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```

str 指针虽然是 unsafe.Pointer 类型，但它最后其实指向了一个 byte 类型的数组，

```go
// go:nosplit
func gostringnocopy(str *byte) string {
	ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)}
	s := *(*string)(unsafe.Pointer(&ss))
	return s
}
```

字符串的赋值其实是指针的复制，同时字符串长度其实调用了 findnull 函数

```go
func findnullw(s *uint16) int {
	if s == nil {
		return 0
	}
	p := (*[maxAlloc/2/2 - 1]uint16)(unsafe.Pointer(s))
	l := 0
	for p[l] != 0 {
		l++
	}
	return l
}

// 在 findnull 的实现中，maxAlloc 是允许用户分配的最大虚拟内存空间。
// 在 64 位，理论上可分配最大 1 << heapAddrBits 字节。
// 在 32 位，最大可分配小于 1 << 32 字节。
// 所以，求长度的逻辑是：
// 如果指针悬空，那么字符串长度为0，否则将指针转换为一个字符数组的指针，然后判断这个字符数组的每个值是否存在，第一个为0的值对应的索引就是字符串的长度。
```



字符串的值不可改变这个特性是通过禁止访问 str 指针指向的内存的值实现的，但 str 指针本身的值是可以改变的，也就是说它指向的内存区域可以改变，所以字符串可以重复赋值。

字符串同时也支持切片操作，我们可以理解为 str 的重新赋值和 len 的重新计算，比如下面的语句中，hello 和 world 其实都指向 s 所指向的内存区域，只是指针的位置不一样。

虽然字符串底层指向一个 byte 数组，单独访问其元素得到的类型也是 byte，但使用 for range 语法遍历时，单个值的类型却是 rune。主要是因为 Go 专门做了一个解码操作，如下：（注意这里的代码不是真的底层实现，只是用来说明逻辑的）

```go
func forOnString(s string, forBody func(i int, r rune)) {
    for i := 0; len(s) > 0; {
        r,size := utf8.DecodeRuneInString(s)
        forBody(i,r)
        s = s[size:]
        i += size
    }
}
```



### 字符串转化

由以上可知字符串单个字符可能是 byte 或 rune，这也是我们使用字符串时经常做的强制类型转换。它们隐含者内存的重新分配，代价可能是不一样的。



1. string 和 []byte，string 和 []rune 的转换都会进行内存的重新分配，有一定代价；
2. 直接访问 string 中的成员，类型为 byte，使用 for range 结构，类型为 rune；
3. 需要修改 string 中的成员时，需要转换 []byte。



#### 1. string->[]byte

```go
func stringtoslicebyte(buf *tmpBuf, s string) []byte {
	var b []byte
	if buf != nil && len(s) <= len(buf) {
		*buf = tmpBuf{}
		b = buf[:len(s)]
	} else {
		b = rawbyteslice(len(s))
	}
	copy(b, s)
	return b
}

// rawbyteslice allocates a new byte slice. The byte slice is not zeroed.
func rawbyteslice(size int) (b []byte) {
	cap := roundupsize(uintptr(size))
	p := mallocgc(cap, nil, false)
	if cap != uintptr(size) {
		memclrNoHeapPointers(add(p, uintptr(size)), cap-uintptr(size))
	}

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, int(cap)}
	return
}
```

可以看到其实做了一次内存的重新分配，得到了新的字符数组 b，然后将 s 复制给 b。至于 copy 函数可以直接把 string 复制给 []byte，是因为 go 源码单独实现了一个`slicestringcopy`函数来实现，具体可以看`src/runtime/slice.go`。



#### 2. []byte->string

```go
func slicebytetostring(buf *tmpBuf, ptr *byte, n int) (str string) {
	if n == 0 {
		// Turns out to be a relatively common case.
		// Consider that you want to parse out data between parens in "foo()bar",
		// you find the indices and convert the subslice to string.
		return ""
	}
	if raceenabled {
		racereadrangepc(unsafe.Pointer(ptr),
			uintptr(n),
			getcallerpc(),
			funcPC(slicebytetostring))
	}
	if msanenabled {
		msanread(unsafe.Pointer(ptr), uintptr(n))
	}
	if n == 1 {
		p := unsafe.Pointer(&staticuint64s[*ptr])
		if sys.BigEndian {
			p = add(p, 7)
		}
		stringStructOf(&str).str = p
		stringStructOf(&str).len = 1
		return
	}

	var p unsafe.Pointer
	if buf != nil && n <= len(buf) {
		p = unsafe.Pointer(buf)
	} else {
		p = mallocgc(uintptr(n), nil, false)
	}
	stringStructOf(&str).str = p
	stringStructOf(&str).len = n
	memmove(p, unsafe.Pointer(ptr), uintptr(n))
	return
}

func stringStructOf(sp *string) *stringStruct {
	return (*stringStruct)(unsafe.Pointer(sp))
}
```

转换的思路是新分配 s，然后将 b 复制给它，所以依然有内存的重新分配。



#### 3. string->[]rune

```go
func stringtoslicerune(buf *[tmpStringBufSize]rune, s string) []rune {
	// two passes.
	// unlike slicerunetostring, no race because strings are immutable.
	n := 0
	for range s {
		n++
	}

	var a []rune
	if buf != nil && n <= len(buf) {
		*buf = [tmpStringBufSize]rune{}
		a = buf[:n]
	} else {
		a = rawruneslice(n)
	}

	n = 0
	for _, r := range s {
		a[n] = r
		n++
	}
	return a
}

// rawruneslice allocates a new rune slice. The rune slice is not zeroed.
func rawruneslice(size int) (b []rune) {
	if uintptr(size) > maxAlloc/4 {
		throw("out of memory")
	}
	mem := roundupsize(uintptr(size) * 4)
	p := mallocgc(mem, nil, false)
	if mem != uintptr(size)*4 {
		memclrNoHeapPointers(add(p, uintptr(size)*4), mem-uintptr(size)*4)
	}

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, int(mem / 4)}
	return
}
```

由于 byte 和 rune 类型的差异，进行内存的重新分配。



#### 4. []rune->string

```go
func slicerunetostring(buf *tmpBuf, a []rune) string {
	if raceenabled && len(a) > 0 {
		racereadrangepc(unsafe.Pointer(&a[0]),
			uintptr(len(a))*unsafe.Sizeof(a[0]),
			getcallerpc(),
			funcPC(slicerunetostring))
	}
	if msanenabled && len(a) > 0 {
		msanread(unsafe.Pointer(&a[0]), uintptr(len(a))*unsafe.Sizeof(a[0]))
	}
	var dum [4]byte
	size1 := 0
	for _, r := range a {
		size1 += encoderune(dum[:], r)
	}
	s, b := rawstringtmp(buf, size1+3)
	size2 := 0
	for _, r := range a {
		// check for race
		if size2 >= size1 {
			break
		}
		size2 += encoderune(b[size2:], r)
	}
	return s[:size2]
}

func rawstringtmp(buf *tmpBuf, l int) (s string, b []byte) {
	if buf != nil && l <= len(buf) {
		b = buf[:l]
		s = slicebytetostringtmp(&b[0], len(b))
	} else {
		s, b = rawstring(l)
	}
	return
}

func rawstring(size int) (s string, b []byte) {
	p := mallocgc(uintptr(size), nil, false)

	stringStructOf(&s).str = p
	stringStructOf(&s).len = size

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}

	return
}
```



## slice

### 源码

```go
type slice struct {
    array unsafe.Pointer   // 用来存储实际数据的数组指针，指向一块连续的内存
    len   int              // 切片中元素的数量
    cap   int              // array数组的长度
}
```

切片可以等于 nil 的，只要其底层指针等于 nil，一般情况是切片声明而未初始化的时候出现该情况, 因为切片等于 nil 一般意味着没有初始化，也就没有使用的价值，切片一旦初始化，底层指针就指向了一个确定的内存区域，但指向的内存区域大小可以为0，所以很少将切片直接和 nil 作比较，使用更多的还是判断切片的长度是否为0。



### 扩容

如果追加的元素不大于剩余容量，则只操作指针，若超过容量则造成内存重新分配。



扩容规则：

1. 如果需要的容量超过原切片容量的两倍，直接使用需要的容量作为新容量；
2. 如果原切片的长度小于 1024，新切片的容量翻倍；
3. 如果原切片的长度大于1024，则每次增加25%，直到新容量超过所需要的容量；

```go
// src/runtime/slice.go growslice 函数
newcap := old.cap
doublecap := newcap + newcap
if cap > doublecap {
    newcap = cap
} else {
    if old.len < 1024 {
        newcap = doublecap
    } else {
        // Check 0 < newcap to detect overflow
        // and prevent an infinite loop.
        for 0 < newcap && newcap < cap {
            newcap += newcap / 4
        }
        // Set newcap to the requested cap when
        // the newcap calculation overflowed.
        if newcap <= 0 {
            newcap = cap
        }
    }
}
```

注意：一般扩容后的容量不会是整数倍，因为会进行内存对齐。详细解释可参考链接

https://shuzang.github.io/2020/golang-deep-learning-3-slice/





## map

map 是通过哈希表+链地址法解决冲突；

map每次遍历得到的元素的顺序不一定相同，和元素的插入顺序无关。基本逻辑是先调用 mapiterinit 初始化 hiter 结构体，然后利用该结构体进行遍历。

### 源码

Go 中映射（map）的底层实现是哈希表，位于 `src/runtime/map.go` 中，数据被放到一个 buckets 数组里，每个 bucket 包含最多 8 个键值对。key 的哈希值低 8 位用于选择 bucket，高 8 位用于区分 bucket 中存放的多个键值。如果超过 8 个键被放到同一个 bucket，使用一个额外的 bucket 来存储。

核心的结构体主要是 hmap 和 bmap，前者就是这个 bucket 数组，后者就是单个 bucket 的结构。

```go
// map的基础数据结构
type hmap struct {
	count     int	 // map存储的元素对计数，len()函数返回此值，所以map的len()时间复杂度是O(1)
	flags     uint8  
	B         uint8  // buckets数组的长度，也就是桶的数量为2^B个
	noverflow uint16 // 溢出的桶的数量的近似值
	hash0     uint32 // hash种子

	buckets    unsafe.Pointer // 指向2^B个桶组成的数组的指针，数据存在这里
	oldbuckets unsafe.Pointer // 指向扩容前的旧buckets数组，只在map增长时有效
	nevacuate  uintptr        // 计数器，标示扩容后搬迁的进度

	extra *mapextra // 保存溢出桶的指针数组和未使用的溢出桶数组的首地址
}

type mapextra struct {
	overflow    *[]*bmap // overflow contains overflow buckets for hmap.buckets.
	oldoverflow *[]*bmap // oldoverflow contains overflow buckets for hmap.oldbuckets.

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}

// 桶的实现结构, hmap的buckets指针指向该结构
type bmap struct {
	// tophash存储桶内每个key的hash值的高字节
	// tophash[0] < minTopHash表示桶的疏散状态
	// 当前版本bucketCnt的值是8，一个桶最多存储8个key-value对
	tophash [bucketCnt]uint8
    // 下面紧跟存放的键值对，存放的格式是所有的 key，然后是所有的 value，
	// 之所以不是一个 key 跟随一个 value，是为了消除填充所需要的间隙，因为
    // key 与 value 的类型不一致，占用的内存大小不一致
    
	// 最后是一个溢出指针
}
```

![https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_1480383-20191104215659319-1712154558.jpg](https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_1480383-20191104215659319-1712154558.jpg)



### 访问

![https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_7515493-599f9d40d5c56e61.webp](https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_7515493-599f9d40d5c56e61.webp)



### 添加

![https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_7515493-54c06b9844da39bd.webp](https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_7515493-54c06b9844da39bd.webp)





### 删除

![https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_7515493-a3221dbfcd6249ab.webp](https://picped-1301226557.cos.ap-beijing.myqcloud.com/Go_20200725_7515493-a3221dbfcd6249ab.webp)



### 扩容

map 扩容每次增加一倍的空间，分配一个新的 Bucket 数组，然后将就数组复制过去。



// 判断是否需要扩容

```go
func (h *hmap) growing() bool {
    return h.oldbuckets != nil
}
```

在分配assign逻辑中，当没有位置给key使用，而且满足测试条件(装载因子>6.5或有太多溢出通)时，会触发（hashGrow）扩容。







## channel

golang 的 channel 就是一个环形队列（ringbuffer）的实现。

> 发生panic的情景：
>
> 1. 关闭已关闭的channel
> 2. 关闭 nil channel
> 3. 向已关闭的channel发送数据





![chan_flow_chart](https://raw.githubusercontent.com/ywanbing/golearning/master/chan_flow_chart/chan_flow_chart.png)



### 源码

```go
type hchan struct {
 qcount   uint           // queue 里面有效用户元素，这个字段是在元素出对，入队改变的；
 dataqsiz uint           // 初始化的时候赋值，之后不再改变，指明数组 buffer 的大小；
 buf      unsafe.Pointer // 指明 buffer 数组的地址，初始化赋值，之后不会再改变；
 elemsize uint16  // 指明元素的大小，和 dataqsiz 配合使用就能知道 buffer 内存块的大小了；
 closed   uint32
 elemtype *_type // 元素类型，初始化赋值；
 sendx    uint   // send index
 recvx    uint   // receive index
 recvq    waitq  // 等待 recv 响应的对象列表，抽象成 waiters
 sendq    waitq  // 等待 sedn 响应的对象列表，抽象成 waiters

 // 互斥资源的保护锁，官方特意说明，在持有本互斥锁的时候，绝对不要修改 Goroutine 的状态，不能很有可能在栈扩缩容的时候，出现死锁
 lock mutex
}
```



```go
// waitq 类型其实就是一个双向列表的实现，和 linux 里面的 LIST 实现非常相像
type waitq struct {
 first *sudog
 last  *sudog
}
```



### send

当我们在 golang 里面执行 `c <- x` 这么一行代码意图投递一个元素到 channel 的时候，其实就是调用到 chansend 函数。这个函数分几个场景来处理，总结来说：

1. 场景一：如果有人（ goroutine ）等着取 channel 的元素，这种场景最快乐，直接把元素给他就完了，然后把它唤醒，hchan 本身递增下 ringbuffer 索引；
2. 场景二：如果 ringbuffer 还有空间，那么就把元素存着，这种也是场景的流程，存和取走的是异步流程，可以把 channel 理解成消息队列，生产者和消费者解耦；
3. 场景三：ringbuffer 没空间，这个时候就要是否需要 block 了，一般来讲，`c <- x` 编译出的代码都是 `block = true` ，那么什么时候 chansend 的 block 参数会是 false 呢？答案是：select 的时候；



```go
// return true 则入队成功

func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    // channel 的所有操作，都在互斥锁下；
    lock(&c.lock)
    // 如果投递的目标是已经关闭的 channel，那么直接 panic；
    if c.closed != 0 {
        unlock(&c.lock)
        panic(plainError("send on closed channel"))
    }
    // 场景一：性能最好的场景，我投递的元素刚好有人在等着（那我直接给他就完了）;
    // 调用的是 send 函数，这个函数后面详细阐述，其实非常简单，递增 sendx, recvx 的索引，然后直接把元素给到等他的人，并且唤醒他；
    if sg := c.recvq.dequeue(); sg != nil {
        send(c, sg, ep, func() { unlock(&c.lock) }, 3)
        return true
    }
    // 场景二：ringbuffer 还有空间，那么把元素放好，递增索引，就可以返回了；
    if c.qcount < c.dataqsiz {
        // 复制，赋值好元素；
        qp := chanbuf(c, c.sendx)
        typedmemmove(c.elemtype, qp, ep)
        // 递增索引
        c.sendx++
        // 回环空间
        if c.sendx == c.dataqsiz {
            c.sendx = 0
        }
        // 递增元素个数
        c.qcount++
        unlock(&c.lock)
        return true
    }
    // 判断是否需要阻塞？如果是非阻塞的，那么就直接解锁返回了，如果是阻塞的场景，那么就会走到下面的逻辑哦；
    // chan <- 和 <-chan 的场景，都是 true，但是会有其他场景这里是 false，可以提前想下？
    if !block {
        unlock(&c.lock)
        return false
    }
    // 代码走到这里，说明都是因为条件不满足，要阻塞当前 goroutine，所以做的事情本质上就是保留好通知路径，等待条件满足，会在这个地方唤醒；
    gp := getg()
    mysg := acquireSudog()
    mysg.releasetime = 0
    mysg.elem = ep
    mysg.waitlink = nil
    mysg.g = gp
    mysg.isSelect = false
    mysg.c = c
    gp.waiting = mysg
    gp.param = nil
    // 把 goroutine 相关的线索结构入队，等待条件满足的唤醒；
    c.sendq.enqueue(mysg)
    // goroutine 切走，让出 cpu 执行权限；
    gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)

    // 到这就是某些人唤醒该 goroutine 了。
    // 下面就是唤醒之后的逻辑了；
    if mysg != gp.waiting {
        throw("G waiting list is corrupted")
    }
    // 做一些资源的释放和环境的清理。
    gp.waiting = nil
    gp.activeStackChans = false
    if gp.param == nil {
        // 做一些校验
        if c.closed == 0 {
            throw("chansend: spurious wakeup")
        }
        panic(plainError("send on closed channel"))
    }
    gp.param = nil
    mysg.c = nil
    releaseSudog(mysg)
    return true
}
```





### recv

chanrecv 函数的返回值有两个值，selected，received，其中 selected 一般作为 select 结合的函数返回值，指明是否要进入 select-case 的代码分支，received 表明是否从队列中成功获取到元素，有几种情况：

1. 如果是非阻塞模式（ block=false ），并且没有任何可用元素，返回 （selected=false，received=false），这样就不会进到 select 的 case 分支；

2. 如果是阻塞模式（ block=true ），如果 chan 已经 closed 了，那么返回的是 （selected=true，received=false），说明需要进到 select 的分支，但是是没有取到元素的；

3. 如果是阻塞模式，chan 还是正常状态，那么返回（selected=true，recived=true），说明正常取到了元素；

   

```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
 // 特殊场景：非阻塞模式，并且没有元素的场景直接就可以返回了，这个分支是快速分支，下面的代码都是在锁内的；
 if !block && (c.dataqsiz == 0 && c.sendq.first == nil ||
  c.dataqsiz > 0 && atomic.Loaduint(&c.qcount) == 0) &&
  atomic.Load(&c.closed) == 0 {
  return
 }

 // 以下所有的逻辑都在锁内；
 lock(&c.lock)

 if c.closed != 0 && c.qcount == 0 {
  if raceenabled {
   raceacquire(c.raceaddr())
  }
  unlock(&c.lock)
  if ep != nil {
   typedmemclr(c.elemtype, ep)
  }
  return true, false
 }

 // 场景：如果发现有个人（sender）正在等着别人接收，那么刚刚好，直接把它的元素给到我们这里就好了；
 if sg := c.sendq.dequeue(); sg != nil {
  recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
  return true, true
 }

 // 场景：ringbuffer 还有空间存元素，那么下面就可以把元素放到 ringbuffer 放好，递增索引，就可以返回了；
 if c.qcount > 0 {
  // 存元素
  qp := chanbuf(c, c.recvx)
  if ep != nil {
   typedmemmove(c.elemtype, ep, qp)
  }
  typedmemclr(c.elemtype, qp)
  // 递增索引
  c.recvx++
  if c.recvx == c.dataqsiz {
   c.recvx = 0
  }
  c.qcount--
  unlock(&c.lock)
  return true, true
 }

 // 代码到这说明 ringbuffer 空间是不够的，后面学会要做两个事情，是否需要阻塞？
 // 如果 block 为 false ，那么直接就退出了，返回对应的返回值；
 if !block {
  unlock(&c.lock)
  return false, false
 }

 // 到这就说明要阻塞等待了，下面唯一要做的就是给阻塞做准备（准备好唤醒的条件）
 gp := getg()
 mysg := acquireSudog()
 mysg.releasetime = 0
 mysg.elem = ep
 mysg.waitlink = nil
 gp.waiting = mysg
 mysg.g = gp
 mysg.isSelect = false
 mysg.c = c
 gp.param = nil
 // goroutine 作为一个 waiter 入队列，等待条件满足之后，从这个队列里取出来唤醒；
 c.recvq.enqueue(mysg)
 // goroutine 切走，交出 cpu 执行权限
 goparkunlock(&c.lock, waitReasonChanReceive, traceEvGoBlockRecv, 3)

 // 这里是被唤醒的开始的地方；
 if mysg != gp.waiting {
  throw("G waiting list is corrupted")
 }
 // 下面做一些资源的清理
 gp.waiting = nil
 closed := gp.param == nil
 gp.param = nil
 mysg.c = nil
 releaseSudog(mysg)
 return true, !closed
}
```



### 遍历

使用`for-range`会编译成`chanrecv2( c, ep )`方法，所以只有这个 chan 被 close 了，否则一直会处于这个死循环内部。

```text
for (   ; ok = chanrecv2( c, ep )  ;   ) {
    // do something
}

func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
    // 注意了，这个 block=true，说明 chanrecv 内部是阻塞的；
    _, received = chanrecv(c, elem, true)        
    return
}
```



参考链接：https://zhuanlan.zhihu.com/p/297053654















