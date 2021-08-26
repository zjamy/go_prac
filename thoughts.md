对当下自己项目中的业务，进行一个微服务改造，需要考虑如下技术点：

1）微服务架构（BFF、Service、Admin、Job、Task 分模块）

2）API 设计（包括 API 定义、错误码规范、Error 的使用）

3）gRPC 的使用

4）Go 项目工程化（项目结构、DI、代码分层、ORM 框架）

5）并发的使用（errgroup 的并行链路请求

6）微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）

7）缓存的使用优化（一致性处理、Pipeline 优化）

My Answer based on the learnings of the class:
**Q: 1）微服务架构（BFF、Service、Admin、Job、Task 分模块）**

A: 根据业务的具体需求，适当的进行逐步的拆分。尽量让微服务之间减少相互依赖。

一般的服务也是模块化逻辑，但是最终它还是会被打包并部署为单体式应用。其中最主要的问题是这个应用太复杂，以至于难以维护，难以扩展（？），可靠性难保障（？），最终，敏捷性开发和部署难以完成。所以需要化繁为简，分而治之。

微服务也满足SOA（[Service-oriented architecture](https://en.wikipedia.org/wiki/Service-oriented_architecture) ）的各种需求

* 小：小在服务代码少，bug少，易测试，易维护，易迭代
* 单一职责，一个服务只做好一件事情
* 尽可能早的创建原型：尽可能早的提供服务API，建立服务契约，达成服务间沟通的一致性约定
* 可移植性高：服务间的轻量级交互协议在效率和可移植性二者间，首先依然考虑兼容性和移植性

![image-20210825075300497](C:\Users\jzhao26\OneDrive - Intel Corporation\SourceCode\go_prac\data\image-20210825075300497.png)



Q: 2）API 设计（包括 API 定义、错误码规范、Error 的使用）

A: 符合API设计的理念。

| API 本身的含义指应用程序接口，包括所依赖的库、平台、操作系统提供的能力都可以叫做 API。我们在讨论微服务场景下的 API 设计都是指 WEB API，一般的实现有 RESTful、RPC等。API 代表了一个微服务实例对外提供的能力，因此 API 的传输格式（XML、JSON）对我们在设计 API 时的影响并不大。 |
| ------------------------------------------------------------ |

Q: 3）gRPC 的使用

A: 虽然client和Server是用restful进行沟通的，但是我们Server端内部最好还是使用GRPC来进行沟通。所以，我们需要用GRPC来模拟对外的restful的API的接口。

## *gRPC是什么？*

*gRPC是什么可以用官网的一句话来概括*

> *A high-performance, open-source universal RPC framework*

***所谓RPC(remote procedure call 远程过程调用)框架实际是提供了一套机制，使得应用程序之间可以进行通信，而且也遵从server/client模型。使用的时候客户端调用server端提供的接口就像是调用本地的函数一样。**如下图所示就是一个典型的RPC结构图。*

*![img](https:////upload-images.jianshu.io/upload_images/3959253-76284b64125a8673.png?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp)*

*RPC通信*

## *gRPC有什么好处以及在什么场景下需要用gRPC*

*既然是server/client模型，那么我们直接用restful api不是也可以满足吗，为什么还需要RPC呢？下面我们就来看看RPC到底有哪些优势*

### *gRPC vs. Restful API*

*gRPC和restful API都提供了一套通信机制，用于server/client模型通信，而且它们都使用http作为底层的传输协议(严格地说, gRPC使用的http2.0，而restful api则不一定)。不过gRPC还是有些特有的优势，如下：*

- *gRPC可以通过protobuf来定义接口，从而可以有更加严格的接口约束条件。关于protobuf可以参见笔者之前的小文[Google Protobuf简明教程](https://www.jianshu.com/p/b723053a86a6)*
- *另外，通过protobuf可以将数据序列化为二进制编码，这会大幅减少需要传输的数据量，从而大幅提高性能。*
- *gRPC可以方便地支持流式通信(理论上通过http2.0就可以使用streaming模式, 但是通常web服务的restful api似乎很少这么用，通常的流式数据应用如视频流，一般都会使用专门的协议如HLS，RTMP等，这些就不是我们通常web服务了，而是有专门的服务器应用。）*

### *使用场景*

- *需要对接口进行严格约束的情况，比如我们提供了一个公共的服务，很多人，甚至公司外部的人也可以访问这个服务，这时对于接口我们希望有更加严格的约束，我们不希望客户端给我们传递任意的数据，尤其是考虑到安全性的因素，我们通常需要对接口进行更加严格的约束。这时gRPC就可以通过protobuf来提供严格的接口约束。*
- *对于性能有更高的要求时。有时我们的服务需要传递大量的数据，而又希望不影响我们的性能，这个时候也可以考虑gRPC服务，因为通过protobuf我们可以将数据压缩编码转化为二进制格式，通常传递的数据量要小得多，而且通过http2我们可以实现异步的请求，从而大大提高了通信效率。*

*但是，通常我们不会去单独使用gRPC，而是将gRPC作为一个部件进行使用，这是因为在生产环境，我们面对大并发的情况下，需要使用分布式系统来去处理，而gRPC并没有提供分布式系统相关的一些必要组件。而且，真正的线上服务还需要提供包括负载均衡，限流熔断，监控报警，服务注册和发现等等必要的组件。不过，这就不属于本篇文章讨论的主题了，我们还是先继续看下如何使用gRPC。*

## *gRPC HelloWorld实例详解*

*gRPC的使用通常包括如下几个步骤：*

1. *通过protobuf来定义接口和数据类型*

2. *编写gRPC server端代码*

3. *编写gRPC client端代码*
    *下面来通过一个实例来详细讲解上述的三步。*
    *下边的hello world实例完成之后，其目录结果如下：*

   *![img](https:////upload-images.jianshu.io/upload_images/3959253-df25b1a5150fe55d.png?imageMogr2/auto-orient/strip|imageView2/2/w/712/format/webp)*

   *project helloworld*

### *定义接口和数据类型*

- *通过protobuf定义接口和数据类型*



```cpp
syntax = "proto3";

package rpc_package;

// define a service
service HelloWorldService {
    // define the interface and data type
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// define the data type of request
message HelloRequest {
    string name = 1;
}

// define the data type of response
message HelloReply {
    string message = 1;
}
```

- *使用gRPC protobuf生成工具生成对应语言的库函数*



```undefined
python -m grpc_tools.protoc -I=./protos --python_out=./rpc_package --grpc_python_out=./rpc_package ./protos/user_info.proto
```

*这个指令会自动生成rpc_package文件夹中的`helloworld_pb2.py`和`helloworld_pb2_grpc.py`，但是不会自动生成`__init__.py`文件，需要我们手动添加*

*关于protobuf的详细解释请参考[Google Protobuf简明教程](https://www.jianshu.com/p/b723053a86a6)*

### *gRPC server端代码*



```python
#!/usr/bin/env python
# -*-coding: utf-8 -*-

from concurrent import futures
import grpc
import logging
import time

from rpc_package.helloworld_pb2_grpc import add_HelloWorldServiceServicer_to_server, \ 
    HelloWorldServiceServicer
from rpc_package.helloworld_pb2 import HelloRequest, HelloReply


class Hello(HelloWorldServiceServicer):

    # 这里实现我们定义的接口
    def SayHello(self, request, context):
        return HelloReply(message='Hello, %s!' % request.name)


def serve():
    # 这里通过thread pool来并发处理server的任务
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # 将对应的任务处理函数添加到rpc server中
    add_HelloWorldServiceServicer_to_server(Hello(), server)

    # 这里使用的非安全接口，世界gRPC支持TLS/SSL安全连接，以及各种鉴权机制
    server.add_insecure_port('[::]:50000')
    server.start()
    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    logging.basicConfig()
    serve()
```

### *gRPC client端代码*



```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function
import logging

import grpc
from rpc_package.helloworld_pb2 import HelloRequest, HelloReply
from rpc_package.helloworld_pb2_grpc import HelloWorldServiceStub

def run():
    # 使用with语法保证channel自动close
    with grpc.insecure_channel('localhost:50000') as channel:
        # 客户端通过stub来实现rpc通信
        stub = HelloWorldServiceStub(channel)

        # 客户端必须使用定义好的类型，这里是HelloRequest类型
        response = stub.SayHello(HelloRequest(name='eric'))
    print ("hello client received: " + response.message)

if __name__ == "__main__":
    logging.basicConfig()
    run()
```

### *演示*

*先执行server端代码*



```css
python hello_server.py
```

*接着执行client端代码如下：*



```css
➜  grpc_test python hello_client.py
hello client received: Hello, eric!
```

## *References*

- *[gRPC官网](https://grpc.io/)*
- *[Google Protobuf简明教程*](https://www.jianshu.com/p/b723053a86a6)



*作者*：geekpy
链接：https://www.jianshu.com/p/9c947d98e192
来源：简书
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

Q: 4）Go 项目工程化（项目结构、DI、代码分层、ORM 框架）

A: 符合课程中工程化的建议。

*作者：Go语言进阶*
*链接：https://zhuanlan.zhihu.com/p/399101012*
*来源：知乎*
*著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。*



*我们在微服务框架**[kratos v2](https://link.zhihu.com/?target=https%3A//github.com/go-kratos/kratos)**的默认项目模板中**[kratos-layout](https://link.zhihu.com/?target=https%3A//github.com/go-kratos/kratos-layout)**使用了**[google/wire](https://link.zhihu.com/?target=https%3A//github.com/google/wire)**进行依赖注入，也建议开发者在维护项目时使用该工具。*

*wire 乍看起来比较违反直觉，导致很多同学不理解为什么要用或不清楚如何用（也包括曾经的我），本文来帮助大家理解 wire 的使用。*

## ***What*** 

***[wire](https://link.zhihu.com/?target=https%3A//github.com/google/wire)**是由 google 开源的一个供 Go 语言使用的依赖注入代码生成工具。它能够根据你的代码，生成相应的依赖注入 go 代码。*

*而与其它依靠反射实现的依赖注入工具不同的是，wire 能在编译期（准确地说是代码生成时）如果依赖注入有问题，在代码生成时即可报出来，不会拖到运行时才报，更便于 debug。*

## ***Why*** 

### ***理解依赖注入***

*什么是依赖注入？为什么要依赖注入？ 依赖注入就是 Java 遗毒（不是）*

***[依赖注入](https://link.zhihu.com/?target=https%3A//zh.wikipedia.org/wiki/%E4%BE%9D%E8%B5%96%E6%B3%A8%E5%85%A5)** (Dependency Injection，缩写为 DI)，可以理解为一种代码的构造模式（就是写法），按照这样的方式来写，能够让你的代码更加容易维护。*

*对于很多软件设计模式和架构的理念，我们都无法理解他们要绕好大一圈做复杂的体操、用奇怪的方式进行实现的意义。他们通常都只是丢出来一段样例，说这样写就很好很优雅，由于省略掉了这种模式是如何发展出来的推导过程，我们只看到了结果，导致理解起来很困难。那么接下来我们来尝试推导还原一下整个过程，看看代码是如何和为什么演进到依赖注入模式的，以便能够更好理解使用依赖注入的意义。*

### ***依赖是什么？***

*这里的依赖是个名词，不是指软件包的依赖（比如那坨塞在 node_modules 里面的东西），而是指软件中某一个模块（对象/实例）所依赖的其它外部模块（对象/实例）。*

### ***注入到哪里？***

*被依赖的模块，在创建模块时，被注入到（即当作参数传入）模块的里面。*

### ***不 DI 是啥样？DI 了又样子？***

*下面用 go 伪代码来做例子，领会精神即可。*

*假设个场景，你在打工搞一个 web 应用，它有一个简单接口。最开始的项目代码可能长这个样子：*

```text
# 下面为伪代码，忽略了很多与主题无关的细节

type App struct {
}

# 假设这个方法将会匹配并处理 GET /biu/<id> 这样的请求
func (a *App) GetData(id string) string {
    # todo: write your data query
    return "some data"
}

func NewApp() *App {
    return &App{}
}

app := App()
app.Run()
```

*你要做的是接一个 mysql，从里面把数据按照 id 查出来，返回。 要连 mysql 的话，假设我们已经有了个`NewMySQLClient`的方法返回 client 给你，初始化时传个地址进去就能拿到数据库连接，并假设它有个`Exec`的方法给你执行参数。*

### ***不用 DI，通过全局变量传递依赖实例***

*一种写法是，在外面全局初始化好 client，然后 App 直接拿来调用。*

```text
var mysqlUrl = "mysql://blabla"
var db = NewMySQLClient(mysqlUrl)


type App struct {

}

func (a *App) GetData(id string) string {
    data := db.Exec("select data from biu where id = ? limit 1", id)
    return data
}


func NewApp() *App {
    return &App{}
}
func main() {
    app := App()
    app.Run()
}
```

*这就是没用依赖注入，app 依赖了全局变量 db，这是比较糟糕的一种做法。db 这个对象游离在全局作用域，暴露给包下的其他模块，比较危险。（设想如果这个包里其他代码在运行时悄悄把你的这个 db 变量替换掉会发生啥）*

### ***不用 DI，在 App 的初始化方法里创建依赖实例***

*另一种方式是这样的：*

```text
type App struct {
    db *MySQLClient
}

func (a *App) GetData(id string) string {
    data := a.db.Exec("select data from biu where id = ? limit 1", id)
    return data
}


func NewApp() *App {
    return &App{db: NewMySQLClient(mysqlUrl)}
}
func main() {
    app := NewApp("mysql://blabla")
    app.Run()
}
```

*这种方法稍微好一些，db 被塞到 app 里面了，不会有 app 之外的无关代码碰它，比较安全，但这依然不是依赖注入，而是在内部创建了依赖，接下来你会看到它带来的问题。*

### ***老板：我们的数据要换个地方存 （需要变更实现）***

*你的老板不知道从哪听说——Redis 贼特么快，要不我们的数据改从 Redis 里读吧。这个时候你的内心有点崩溃，但毕竟要恰饭的，就硬着头皮改上面的代码。*

```text
type App struct {
    ds *RedisClient
}

func (a *App) GetData(id string) string {
    data := a.ds.Do("GET", "biu_"+id)
    return data
}


func NewApp() *App {
    return &App{ds: NewRedisClient(redisAddr)}
}

func main() {
    app := NewApp("redis://ooo")
    app.Run()
}
```

*上面基本进行了 3 处修改：*

1. *App 初始化方法里改成了初始化 RedisClient*
2. *get_data 里取数据时改用 run 方法，并且查询语句也换了*
3. *App 实例化时传入的参数改成了 redis 地址*

### ***老板：要不，我们再换个地方存？/我们要加测试，需要 Mock***

*老板的思路总是很广的，又过了两天他又想换成 Postgres 存了；或者让你们给 App 写点测试代码，只测接口里面的逻辑，通常我们不太愿意在旁边再起一个数据库，那么就需要 mock 掉数据源这块东西，让它直接返回数据给请求的 handler 用，来进行针对性的测试。*

*这种情况怎么办？再改里面的代码？这不科学。*

### ***面向接口编程***

*一个很重要的思路就是要**面向接口(interface)编程**，而不是面向具体实现编程。*

*什么叫面向具体实现编程呢？比如上述的例子里改动的部分：调 mysqlclient 的 exec_sql 执行一条 sql，被改成了：调 redisclient 的 do 执行一句 get 指令。由于每种 client 的接口设计不同，每换一个实现，就得改一遍。*

*而面向接口编程的思路，则完全不同。我们不要听老板想用啥就马上写代码。首先就得预料到，这个数据源的实现很有可能被更换，因此在一开始就应该做好准备（设计）。*

### ***设计接口***

*Python 里面有个概念叫鸭子类型(duck-typing)，就是如果你叫起来像鸭子，走路像鸭子，游泳像鸭子，那么你就是一只鸭子。这里的叫、走路、游泳就是我们约定的鸭子接口，而你如果完整实现了这些接口，我们可以像对待一个鸭子一样对待你。*

*在我们上面的例子中，不论是 Mysql 实现还是 Redis 实现，他们都有个共同的功能：用一个 id，查一个数据出来，那么这就是共同的接口。*

*我们可以约定一个叫 DataSource 的接口，它必须有一个方法叫 GetById，功能是要接收一个 id，返回一个字符串*

```text
type DataSource interface {
    GetById(id string) string
}
```

*然后我们就可以把各个数据源分别进行封装，按照这个 interface 定义实现接口，这样我们的 App 里处理请求的部分就可以稳定地调用 GetById 这个方法，而底层数据实现只要实现了 DataSource 这个 interface 就能花式替换，不用改 App 内部的代码了。*

```text
// 封装个redis
type redis struct {
    r *RedisClient
}

func NewRedis(addr string) *redis {
    return &redis{db: NewRedisClient(addr)}
}

func (r *redis) GetById(id string) string {
    return r.r.Do("GET", "biu_"+id)
}


// 再封装个mysql
type mysql struct {
    m *MySQLClient
}

func NewMySQL(addr string) *redis {
    return &mysql{db: NewMySQLClient(addr)}
}

func (m *mysql) GetById(id string) string {
    return r.m.Exec("select data from biu where id = ? limit 1", id)
}


type App struct {
    ds DataSource
}

func NewApp(addr string) *App {
    //需要用Mysql的时候
    return &App{ds: NewMySQLClient(addr)}

    //需要用Redis的时候
    return &App{ds: NewRedisClient(addr)}
}
```

*由于两种数据源都实现了 DataSource 接口，因此可以直接创建一个塞到 App 里面了，想用哪个用哪个，看着还不错？*

### ***等一等，好像少了些什么***

*addr 作为参数，是不是有点简单？通常初始化一个数据库连接，可能有一堆参数，配在一个 yaml 文件里，需要解析到一个 struct 里面，然后再传给对应的 New 方法。*

*配置文件可能是这样的：*

```text
redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
```

*解析结构体是这样的：*

```text
type RedisConfig struct {
 Network      string             `json:"network,omitempty"`
 Addr         string             `json:"addr,omitempty"`
 ReadTimeout  *duration.Duration `json:"read_timeout,omitempty"`
 WriteTimeout *duration.Duration `json:"write_timeout,omitempty"`
}
```

*结果你的`NewApp`方法可能就变成了这个德性：*

```text
func NewApp() *App {
    var conf *RedisConfig
    yamlFile, err := ioutil.ReadFile("redis_conf.yaml")
    if err != nil {
        panic(err)
    }
    err = yaml.Unmarshal(yamlFile, &conf)
    if err != nil {
        panic(err)
    }
    return &App{ds: NewRedisClient(conf)}
}
```

*NewApp 说，停停，你们年轻人不讲武德，我的责任就是创建一个 App 实例，我只需要一个 DataSource 注册进去，至于这个 DataSource 是怎么来的我不想管，这么一坨处理 conf 的代码凭什么要放在我这里，我也不想关心你这配置文件是通过网络请求拿来的还是从本地磁盘读的，我只想把 App 组装好扔出去直接下班。*

### ***依赖注入终于可以登场了***

*还记得前面是怎么说依赖注入的吗？被依赖的模块，在创建模块时，被注入到（即当作参数传入）初始化函数里面。通过这种模式，正好可以让 NewApp 早点下班。我们在外面初始化好 NewRedis 或者 NewMysql，得到的 DataSource 直接扔给 NewApp。*

*也就是这样*

```text
func NewApp(ds DataSource) *App {
    return &App{ds: ds}
}
```

*那坨读配置文件初始化 redis 的代码扔到初始化 DataSource 的方法里去*

```text
func NewRedis() DataSource {
    var conf *RedisConfig
    yamlFile, err := ioutil.ReadFile("redis_conf.yaml")
    if err != nil {
        panic(err)
    }
    err = yaml.Unmarshal(yamlFile, &conf)
    if err != nil {
        panic(err)
    }
    return &redis{r: NewRedisClient(conf)}
}
```

*更进一步，NewRedis 这个方法甚至也不需要关心文件是怎么读的，它的责任只是通过 conf 初始化一个 DataSource 出来，因此你可以继续把读 config 的代码往外抽，把 NewRedis 做成接收一个 conf，输出一个 DataSource*

```text
func GetRedisConf() *RedisConfig
func NewRedis(conf *RedisConfig) DataSource
```

*因为之前整个组装过程是散放在 main 函数下面的，我们把它抽出来搞成一个独立的 initApp 方法。最后你的 App 初始化逻辑就变成了这样*

```text
func initApp() *App {
    c := GetRedisConf()
    r := NewRedis(c)
    app := NewApp(r)
    return app
}

func main() {
    app := initApp()
    app.Run()
}
```

*然后你可以通过实现 DataSource 的接口，更换前面的读取配置文件的方法，和更换创建 DataSource 的方法，来任意修改你的底层实现（读配置文件的实现，和用哪种 DataSource 来查数据），而不用每次都改一大堆代码。这使得你的代码层次划分得更加清楚，更容易维护了。*

*这就是依赖注入。*

### ***手工依赖注入的问题***

*上文这一坨代码，把各个实例初始化好，再按照各个初始化方法的需求塞进去，最终构造出 app 的这坨代码，就是注入依赖的过程。*

```text
c := GetRedisConf()
r := NewRedis(c)
app := NewApp(r)
```

*目前只有一个 DataSource，这样手写注入过程还可以，一旦你要维护的东西多了，比如你的 NewApp 是这样的`NewApp(r *Redis, es *ES, us *UserSerivce, db *MySQL) *App`然后其中 UserService 是这样的`UserService(pg *Postgres, mm *Memcached)`，这样形成了多层次的一堆依赖需要注入，徒手去写非常麻烦。*

*而这部分，就是 wire 这样的依赖注入工具能够起作用的地方了——他的功能只是通过生成代码**帮你注入依赖**，而实际的依赖实例需要你自己创建（初始化）。*

## ***How*** 

*wire 的主要问题是，看文档学不会。反正我最初看完文档之后是一头雾水——这是啥，这要干啥？但通过我们刚才的推导过程，应该大概理解了为什么要用依赖注入，以及 wire 在这其中起到什么作用——通过生成代码**帮你注入依赖**，而实际的依赖实例需要你自己创建（初始化）。*

*接下来就比较清楚了。*

*首先要实现一个`wire.go`的文件，里面定义好 Injector。*

```text
// +build wireinject

func initApp() (*App) {
 panic(wire.Build(GetRedisConf, NewRedis, SomeProviderSet, NewApp))
}
```

*然后分别实现好 Provider。*

*执行`wire`命令后 他会扫描整个项目，并帮你生成一个`wire_gen.go`文件，如果你有什么没有实现好，它会报错出来。*

*你学会了吗？*

### ***重新理解***

*等一等，先别放弃治疗，让我们用神奇的中文编程来解释一下要怎么做。*

### ***谁参与编译？***

*上面那个`initApp`方法，官方文档叫它 Injector，由于文件里首行`// +build wireinject`这句注释，这个 wire.go 文件只会由 wire 读取，在 go 编译器在编译代码时不会去管它，实际会读的是生成的 wire_gen.go 文件。*

*而 Provider 就是你代码的一部分，肯定会参与到编译过程。*

### ***Injector 是什么鬼东西？***

*Injector 就是你最终想要的结果——最终的 App 对象的初始化函数，也就是前面那个例子里的`initApp`方法。*

*把它理解为你去吃金拱门，进门看到点餐机，噼里啪啦点了一堆，最后打出一张单子。*

```text
// +build wireinject

func 来一袋垃圾食品() 一袋垃圾食品 {
    panic(wire.Build(来一份巨无霸套餐, 来一份双层鳕鱼堡套餐, 来一盒麦乐鸡, 垃圾食品打包))
}
```

*这就是你点的单子，它不参与编译，实际参与编译的代码是由 wire 帮你生成的。*

### ***Provider 是什么鬼东西？***

*Provider 就是创建各个依赖的方法，比如前面例子里的 NewRedis 和 NewApp 等。*

*你可以理解为，这些是金拱门的服务员和后厨要干的事情： 金拱门后厨需要提供这些食品的制作服务——实现这些实例初始化方法。*

```text
func 来一盒麦乐鸡() 一盒麦乐鸡 {}
func 垃圾食品打包(一份巨无霸套餐, 一份双层鳕鱼堡套餐, 一盒麦乐鸡) 一袋垃圾食品 {}
```

*wire 里面还有个 ProviderSet 的概念，就是把一组 Provider 打包，因为通常你点单的时候很懒，不想这样点你的巨无霸套餐：我要一杯可乐，一包薯条，一个巨无霸汉堡；你想直接戳一下就好了，来一份巨无霸套餐。这个套餐就是 ProviderSet，一组约定好的配方，不然你的点单列表（injector 里的 Build）就会变得超级长，这样你很麻烦，服务员看着也很累。*

*用其中一个套餐举例*

```text
// 先定义套餐内容
var 巨无霸套餐 = wire.NewSet(来一杯可乐，来一包薯条，来一个巨无霸汉堡)

// 然后实现各个食品的做法
func 来一杯可乐() 一杯可乐 {}
func 来一包薯条() 一包薯条 {}
func 来一个巨无霸汉堡() 一个巨无霸汉堡 {}
```

### ***wire 工具做了啥？***

*重要的事情说三遍，通过生成代码**帮你注入依赖**。*

*在金拱门的例子里就是，wire 就是个服务员，它按照你的订单，去叫做相应的同事把各个食物/套餐做好，然后最终按需求打包给你。这个中间协调构建的过程，就是注入依赖。*

*这样的好处就是， 对于金拱门，假设他们突然换可乐供应商了，直接把`来一杯可乐`替换掉就行，返回一种新的可乐，而对于顾客不需要有啥改动。 对于顾客来说，点单内容可以变换，比如我今天不想要麦乐鸡了，或者想加点别的，只要改动我的点单(只要金拱门能做得出来)，然后通过 wire 重新去生成即可，不需要关注这个服务员是如何去做这个订单的。*

*现在你应该大概理解 wire 的用处和好处了。*

### ***总结***

*让我们从金拱门回来，重新总结一下用 wire 做依赖注入的过程。*

### ***1. 定义 Injector***

*创建`wire.go`文件，定义下你最终想用的实例初始化函数例如`initApp`（即 Injector），定好它返回的东西`*App`，在方法里用`panic(wire.Build(NewRedis, SomeProviderSet, NewApp))`罗列出它依赖哪些实例的初始化方法（即 Provider）/或者哪些组初始化方法（ProviderSet）*

### ***2. 定义 ProviderSet（如果有的话）***

*ProviderSet 就是一组初始化函数，是为了少写一些代码，能够更清晰的组织各个模块的依赖才出现的。也可以不用，但 Injector 里面的东西就需要写一堆。 像这样 `var SomeProviderSet = wire.NewSet(NewES,NewDB)`定义 ProviderSet 里面包含哪些 Provider*

### ***3. 实现各个 Provider***

*Provider 就是初始化方法，你需要自己实现，比如 NewApp，NewRedis，NewMySQL，GetConfig 等，注意他们们各自的输入输出*

### ***4. 生成代码***

*执行 wire 命令生成代码，工具会扫描你的代码，依照你的 Injector 定义来组织各个 Provider 的执行顺序，并自动按照 Provider 们的类型需求来按照顺序执行和安排参数传递，如果有哪些 Provider 的要求没有满足，会在终端报出来，持续修复执行 wire，直到成功生成`wire_gen.go`文件。接下来就可以正常使用`initApp`来写你后续的代码了。*

*如果需要替换实现，对 Injector 进行相应的修改，实现必须的 Provider，重新生成即可。*

*它生成的代码其实就是类似我们之前需要手写的这个*

```text
func initApp() *App {  // injector
    c := GetRedisConf() // provider
    r := NewRedis(c)  // provider
    app := NewApp(r) // provider
    return app
}
```

*由于我们的例子比较简单，通过 wire 生成体现不出优势，但如果我们的软件复杂，有很多层级的依赖，使用 wire 自动生成注入逻辑，无疑更加方便和准确。*

### ***5. 高级用法***

*wire 还有更多功能，比如 cleanup, bind 等等，请参考官方文档来使用。*

*最后，其实多折腾几次，就会使用了，希望本文能对您起到一定程度上的帮助。*

## ***相关文献*** 

- *[https://github.com/google/wire](https://link.zhihu.com/?target=https%3A//github.com/google/wire)*
- *[https://go-kratos.dev/docs/getting-started/wire](https://link.zhihu.com/?target=https%3A//go-kratos.dev/docs/getting-started/wire)*
- *[https://github.com/go-kratos/kratos-layout](https://link.zhihu.com/?target=https%3A//github.com/go-kratos/kratos-layout)*

*\- END -*



Q: 5）并发的使用（errgroup 的并行链路请求）

A: 我们给错误做了分级定义，不是所有的错误都是fatal错误。

以下内容来自：https://zhuanlan.zhihu.com/p/397673996

> *哈喽，大家好，我是`asong`，今天给大家介绍一个并发编程包`errgroup`，其实这个包就是对`sync.waitGroup`的封装。我们在之前的文章—— **[源码剖析sync.WaitGroup(文末思考题你能解释一下吗?)](https://link.zhihu.com/?target=https%3A//mp.weixin.qq.com/s/hofXXzFhu-rk3_6i2X4m6A)**，从源码层面分析了`sync.WaitGroup`的实现，使用`waitGroup`可以实现一个`goroutine`等待一组`goroutine`干活结束，更好的实现了任务同步，但是`waitGroup`却无法返回错误，当一组`Goroutine`中的某个`goroutine`出错时，我们是无法感知到的，所以`errGroup`对`waitGroup`进行了一层封装，封装代码仅仅不到`50`行，下面我们就来看一看他是如何封装的？*

## ***`errGroup`如何使用***

*老规矩，我们先看一下`errGroup`是如何使用的，前面吹了这么久，先来验验货；*

*以下来自官方文档的例子：*

```text
var (
 Web   = fakeSearch("web")
 Image = fakeSearch("image")
 Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
 return func(_ context.Context, query string) (Result, error) {
  return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
 }
}

func main() {
 Google := func(ctx context.Context, query string) ([]Result, error) {
  g, ctx := errgroup.WithContext(ctx)

  searches := []Search{Web, Image, Video}
  results := make([]Result, len(searches))
  for i, search := range searches {
   i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
   g.Go(func() error {
    result, err := search(ctx, query)
    if err == nil {
     results[i] = result
    }
    return err
   })
  }
  if err := g.Wait(); err != nil {
   return nil, err
  }
  return results, nil
 }

 results, err := Google(context.Background(), "golang")
 if err != nil {
  fmt.Fprintln(os.Stderr, err)
  return
 }
 for _, result := range results {
  fmt.Println(result)
 }

}
```

*上面这个例子来自官方文档，代码量有点多，但是核心主要是在`Google`这个闭包中，首先我们使用`errgroup.WithContext`创建一个`errGroup`对象和`ctx`对象，然后我们直接调用`errGroup`对象的`Go`方法就可以启动一个协程了，`Go`方法中已经封装了`waitGroup`的控制操作，不需要我们手动添加了，最后我们调用`Wait`方法，其实就是调用了`waitGroup`方法。这个包不仅减少了我们的代码量，而且还增加了错误处理，对于一些业务可以更好的进行并发处理。*

## ***赏析`errGroup`***

### ***数据结构***

*我们先看一下`Group`的数据结构：*

```text
type Group struct {
 cancel func() // 这个存的是context的cancel方法

 wg sync.WaitGroup // 封装sync.WaitGroup

 errOnce sync.Once // 保证只接受一次错误
 err     error // 保存第一个返回的错误
}
```

### ***方法解析***

```text
func WithContext(ctx context.Context) (*Group, context.Context)
func (g *Group) Go(f func() error)
func (g *Group) Wait() error
```

*`errGroup`总共只有三个方法：*

- *`WithContext`方法*

```text
func WithContext(ctx context.Context) (*Group, context.Context) {
 ctx, cancel := context.WithCancel(ctx)
 return &Group{cancel: cancel}, ctx
}
```

*这个方法只有两步：*

- *使用`context`的`WithCancel()`方法创建一个可取消的`Context`*
- *创建`cancel()`方法赋值给`Group`对象*
- *`Go`方法*

```text
func (g *Group) Go(f func() error) {
 g.wg.Add(1)

 go func() {
  defer g.wg.Done()

  if err := f(); err != nil {
   g.errOnce.Do(func() {
    g.err = err
    if g.cancel != nil {
     g.cancel()
    }
   })
  }
 }()
}
```

*`Go`方法中运行步骤如下：*

- *执行`Add()`方法增加一个计数器*
- *开启一个协程，运行我们传入的函数`f`，使用`waitGroup`的`Done()`方法控制是否结束*
- *如果有一个函数`f`运行出错了，我们把它保存起来，如果有`cancel()`方法，则执行`cancel()`取消其他`goroutine`*

*这里大家应该会好奇为什么使用`errOnce`，也就是`sync.Once`，这里的目的就是保证获取到第一个出错的信息，避免被后面的`Goroutine`的错误覆盖。*

- *`wait`方法*

```text
func (g *Group) Wait() error {
 g.wg.Wait()
 if g.cancel != nil {
  g.cancel()
 }
 return g.err
}
```

*总结一下`wait`方法的执行逻辑：*

- *调用`waitGroup`的`Wait()`等待一组`Goroutine`的运行结束*
- *这里为了保证代码的健壮性，如果前面赋值了`cancel`，要执行`cancel()`方法*
- *返回错误信息，如果有`goroutine`出现了错误才会有值*

### ***小结***

*到这里我们就分析完了`errGroup`包，总共就`1`个结构体和`3`个方法，理解起来还是比较简单的，针对上面的知识点我们做一个小结：*

- *我们可以使用`withContext`方法创建一个可取消的`Group`，也可以直接使用一个零值的`Group`或`new`一个`Group`，不过直接使用零值的`Group`和`new`出来的`Group`出现错误之后就不能取消其他`Goroutine`了。*
- *如果多个`Goroutine`出现错误，我们只会获取到第一个出错的`Goroutine`的错误信息，晚于第一个出错的`Goroutine`的错误信息将不会被感知到。*
- *`errGroup`中没有做`panic`处理，我们在`Go`方法中传入`func() error`方法时要保证程序的健壮性*

## ***踩坑日记***

*使用`errGroup`也并不是一番风顺的，我之前在项目中使用`errGroup`就出现了一个`BUG`，把它分享出来，避免踩坑。*

*这个需求是这样的(并不是真实业务场景，由`asong`虚构的)：开启多个`Goroutine`去缓存中设置数据，同时开启一个`Goroutine`去异步写日志，很快我的代码就写出来了：*

```text
func main()  {
 g, ctx := errgroup.WithContext(context.Background())

 // 单独开一个协程去做其他的事情，不参与waitGroup
 go WriteChangeLog(ctx)

 for i:=0 ; i< 3; i++{
  g.Go(func() error {
   return errors.New("访问redis失败\n")
  })
 }
 if err := g.Wait();err != nil{
  fmt.Printf("appear error and err is %s",err.Error())
 }
 time.Sleep(1 * time.Second)
}

func WriteChangeLog(ctx context.Context) error {
 select {
 case <- ctx.Done():
  return nil
 case <- time.After(time.Millisecond * 50):
  fmt.Println("write changelog")
 }
 return nil
}
// 运行结果
appear error and err is 访问redis失败
```

*代码没啥问题吧，但是日志一直没有写入，排查了好久，终于找到问题原因。原因就是这个`ctx`。*

*因为这个`ctx`是`WithContext`方法返回的一个带取消的`ctx`，我们把这个`ctx`当作父`context`传入`WriteChangeLog`方法中了，如果`errGroup`取消了，也会导致上下文的`context`都取消了，所以`WriteChangelog`方法就一直执行不到。*

*这个点是我们在日常开发中想不到的，所以需要注意一下～。*

## ***总结***

*因为最近看很多朋友都不知道这个库，所以今天就把他分享出来了，封装代码仅仅不到`50`行，真的是很厉害，如果让你来封装，你能封装的更好吗？*

*欢迎关注公众号：【Golang梦工厂】*







Q: 6）微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）

A: 没有使用kafka，没有使用公开的中间件。

Q: 7）缓存的使用优化（一致性处理、Pipeline 优化）

A: 减少依赖，减少缓存穿透，根据需要合适布置缓存。

