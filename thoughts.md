对当下自己项目中的业务，进行一个微服务改造，需要考虑如下技术点：
1）微服务架构（BFF、Service、Admin、Job、Task 分模块）
2）API 设计（包括 API 定义、错误码规范、Error 的使用）
3）gRPC 的使用
4）Go 项目工程化（项目结构、DI、代码分层、ORM 框架）
5）并发的使用（errgroup 的并行链路请求
6）微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）
7）缓存的使用优化（一致性处理、Pipeline 优化）

My Answer:
Q: 1）微服务架构（BFF、Service、Admin、Job、Task 分模块）
A: 根据业务的具体需求，适当的进行逐步的拆分。尽量让微服务之间减少相互依赖。
Q: 2）API 设计（包括 API 定义、错误码规范、Error 的使用）
A: 符合API设计的理念。
Q: 3）gRPC 的使用
A: 虽然client和Server是用restful进行沟通的，但是我们Server端内部最好还是使用GRPC来进行沟通。所以，我们需要用GRPC来模拟对外的restful的API的接口。
Q: 4）Go 项目工程化（项目结构、DI、代码分层、ORM 框架）
A: 符合课程中工程化的建议。
Q: 5）并发的使用（errgroup 的并行链路请求）
A: 
Q: 6）微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）
A: 我们没有使用kafka.
Q: 7）缓存的使用优化（一致性处理、Pipeline 优化）
A: 
