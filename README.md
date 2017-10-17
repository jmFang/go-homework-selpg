# go-homework-selpg
课程-服务计算的第二次作业-selpg linux CMI的golang实现

## 1. selpg 简介

### 1.1 selpg 程序逻辑

如前面所说的那样，selpg 是从文本输入选择页范围的实用程序。该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。

selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。
参数处理

**“-sNumber”和“-eNumber”强制选项：**

selpg 要求用户用两个命令行参数“-sNumber”（例如，“-s10”表示从第 10 页开始）和“-eNumber”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。selpg 对所给的页号进行合理性检查；换句话说，它会检查两个数字是否为有效的正整数以及结束页是否不小于起始页。这两个选项，“-sNumber”和“-eNumber”是强制性的，而且必须是命令行上在命令名 selpg 之后的头两个参数：

> $ selpg -s10 -e20 ...

（... 是命令的余下部分，下面对它们做了描述）。

** “-lNumber”和“-f”可选选项： **
selpg 可以处理两种输入文本：

**类型 1** ：该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-lNumber”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。

选择 72 作为缺省值是因为在行打印机上这是很常见的页长度。这样做的意图是将最常见的命令用法作为缺省值，这样用户就不必输入多余的选项。该缺省值可以用“-lNumber”选项覆盖，如下所示：

> $ selpg -s10 -e20 -l66 ...

这表明页有固定长度，每页为 66 行。

**类型 2** ：该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界。该格式与“每页行数固定”格式相比的好处在于，当每页的行数有很大不同而且文件有很多页时，该格式可以节省磁盘空间。在含有文本的行后面，类型 2 的页只需要一个字符 ― 换页 ― 就可以表示该页的结束。打印机会识别换页符并自动根据在新的页开始新行所需的行数移动打印头。

将这一点与类型 1 比较：在类型 1 中，文件必须包含 PAGELEN - CURRENTPAGELEN 个新的行以将文本移至下一页，在这里 PAGELEN 是固定的页大小而 CURRENTPAGELEN 是当前页上实际文本行的数目。在此情况下，为了使打印头移至下一页的页首，打印机实际上必须打印许多新行。这在磁盘空间利用和打印机速度方面效率都很低（尽管实际的区别可能不太大）。

类型 2 格式由“-f”选项表示，如下所示：

> $ selpg -s10 -e20 -f ...

该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。

**注：“-lNumber”和“-f”选项是互斥的。 **

**“-dDestination”可选选项： **

selpg 还允许用户使用“-dDestination”选项将选定的页直接发送至打印机。这里，“Destination”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。该目的地应该存在 ― selpg 不检查这一点。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“Destination”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出。这一特性是用 popen() 系统调用实现的，该系统调用允许一个进程打开到另一个进程的管道，将管道用于输出或输入。在下面的示例中，我们打开到命令

> $ lp -dDestination

的管道以便输出，并写至该管道而不是标准输出：

> selpg -s10 -e20 -dlp1

该命令将选定的页作为打印作业发送至 lp1 打印目的地。您应该可以看到类似“request id is lp1-6”的消息。该消息来自 lp 命令；它显示打印作业标识。如果在运行 selpg 命令之后立即运行命令 lpstat -t | grep lp1 ，您应该看见 lp1 队列中的作业。如果在运行 lpstat 命令前耽搁了一些时间，那么您可能看不到该作业，因为它一旦被打印就从队列中消失了。
输入处理

一旦处理了所有的命令行参数，就使用这些指定的选项以及输入、输出源和目标来开始输入的实际处理。

selpg 通过以下方法记住当前页号：如果输入是每页行数固定的，则 selpg 统计新行数，直到达到页长度后增加页计数器。如果输入是换页定界的，则 selpg 改为统计换页符。这两种情况下，只要页计数器的值在起始页和结束页之间这一条件保持为真，selpg 就会输出文本（逐行或逐字）。当那个条件为假（也就是说，页计数器的值小于起始页或大于结束页）时，则 selpg 不再写任何输出。瞧！您得到了想输出的那些页。

<hr>

## 2. Go lang 设计与实现

方法1：golang 自带处理命令行参数的flag包，调用该package，使用flag来解析参数。

方法二：使用os包下的os.Args来解析命令，同时使用一些数据结构来使得输入操作符合要求。

因为flag包处理参数的方法中要求子命令和命令的值分开才能实现给子命令赋值，如 ./selpg -s 1 -e 2， 为了符合要求，使得输入类似 ：。/ selpg -s -e2，使用一些中间变量保存解析的结果。具体说明如下：

### 2.1 数据结构

#### 2.1.1 方法一中使用flag包解析命令需要用到的自命令集合 flag subcommand-set

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/subcommandset.png)

#### 2.1.2 方法二中处理一页或输入输出需要用到的参数变量

struct

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/struct.png)

### 2.2 参数处理

#### 格式认证：

方法二中利用string的向量和切片性质对输入的命令做字符串比较和切割，从而完成格式认证

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/format-f2.png)

方法一中要求输入的自命令和子命令的值用空格分开，因此可直接判断，无须多余的解析

![]()

将解析完成的结果存放进结构体中，传递给IO处理模块，做下一步处理，以上工作由processArgs模块完成

### 2.3 IO处理

#### 2.3.1 根据参数解析的结果，设置输入输出# 2.4 

如果输入文件名为空，那么从标准输入（键盘输入），否则从文件输入；如果-d参数为空值，则只是标准的输出，否则根据-d启动另一个线程，利用shell执行该命令。

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/sh.png)

#### 2.3.2 根据参数解析的结果，设置换页类型

如果是默认的l-type，则根据行数来换页；

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/pagetype-l.png)

如果是强加的-f类型，则根据换页符来换页。

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/pagetype-f.png)

如果-d为空，对于默认的-l类型，因为是标准输出所以写出的时候可以按行写出如果-d不为空，那么先把读取的数据存放在一个buffer，最后再一次写出。

![](https://github.com/jmFang/go-homework-selpg/blob/master/image/d-f.png)

### 2.4 总结

重点是解析命令和管道重定向，命令解析可以使用flag包，也可以使用os.Args，对于不同的输入格式，两者各有千秋。golang中的重定向可以通过exec包的cmd开启新线程，执行重定向。输入输出使用io包，等等。程序实现逻辑不难，难在刚接触golang，对于golang的语法和包使用不熟悉，以及相关学习资源较少。整个程序基本是按照原版的selpg.c用golang翻译过来的，算是一次对golang语法使用的练手。

**如果有bug，请不吝赐教！**

## 3. 参考资料

【1】. Building a Simple CLI Tool with Golang ：https://blog.komand.com/build-a-simple-cli-tool-with-golang

【2】. CLI: Command Line Programming with Go - The New Stack：https://thenewstack.io/cli-command-line-programming-with-go/

【3】. 开发 Linux 命令行实用程序：https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html

【4】. Package flag: https://go-zh.org/pkg/flag/

【5】. Package os: https://go-zh.org/pkg/os/





















