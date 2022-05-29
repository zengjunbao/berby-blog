https://wu-sheng.gitbooks.io/opentracing-io/content/

# Open Tracing

当一个生产系统面对真正的高并发，或者解耦成大量微服务时，以前很容易实现的重点任务变得困难。

OpenTracing通过提供平台无关、厂商无关的API，使得开发人员能够方便的实现追踪系统。



- trace

> 在广义上，一个trace代表了一个事务或者流程在（分布式）系统中的执行过程。
>
> 在OpenTracing标准中，trace是多个span组成的一个有向无环图（DAG），每一个span代表trace中被命名并计时的连续性的执行片段。在一个常规的RPC调用过程中，OpenTracing推荐在RPC的客户端和服务端，至少各有一个span，用于记录RPC调用的客户端和服务端信息。
>
> (传统系统调用图它不能很好显示组件的调用时间，是串行调用还是并行调用，如果展现更复杂的调用关系，会更加复杂)
>
> 显示了执行时间的上下文，相关服务间的层次关系，进程或者任务的串行或并行调用关系。这样的视图有助于发现系统调用的关键路径。通过关注关键路径的执行过程，项目团队可能专注于优化路径中的关键位置，最大幅度的提升系统性能。

## ChildOf 和 FollowsFrom

这两种引用类型代表了子节点和父节点间的直接因果关系。

> **`ChildOf`引用:** 一个span可能是一个父级span的孩子，即"ChildOf"关系。在"ChildOf"引用关系下，父级span某种程度上取决于子span。
>
> **`FollowsFrom` 引用:** 一些父级节点不以任何方式依然他们子节点的执行结果。

## Key

> ## Traces
>
> 一个trace代表一个潜在的，分布式的，存在并行数据或并行执行轨迹（潜在的分布式、并行）的系统。一个trace可以认为是多个span的有向无环图（DAG）。
>
> ## Spans
>
> 一个span代表系统中具有开始时间和执行时长的逻辑运行单元。span之间通过嵌套或者顺序排列建立逻辑因果关系。
>
> ### Operation Names
>
> 每一个span都有一个操作名称，这个名称简单，并具有可读性高。（例如：一个RPC方法的名称，一个函数名，或者一个大型计算过程中的子任务或阶段）。span的操作名应该是一个抽象、通用的标识，能够明确的、具有统计意义的名称；更具体的子类型的描述，请使用Tags
>
> ### Logs
>
> 每个span可以进行多次**Logs**操作，每一次**Logs**操作，都需要一个带时间戳的时间名称，以及可选的任意大小的存储结构。
>
> ### Tags
>
> 每个span可以有多个键值对（key:value）形式的**Tags**，**Tags**是没有时间戳的，支持简单的对span进行注解和补充。
>
> ## SpanContext
>
> 每个span必须提供方法访问**SpanContext**。SpanContext代表跨越进程边界，传递到下级span的状态。(例如，包含`<trace_id, span_id, sampled>`元组)，并用于封装**Baggage** (关于Baggage的解释，请参考下文)。
>
> ### Baggage
>
> **Baggage**是存储在SpanContext中的一个键值对(SpanContext)集合。它会在一条追踪链路上的所有span内*全局*传输，包含这些span对应的SpanContexts。
>
> Baggage拥有强大功能，也会有很大的*消耗*。由于Baggage的全局传输，如果包含的数量量太大，或者元素太多，它将降低系统的吞吐量或增加RPC的延迟。
>
> ## Baggage vs. Span Tags
>
> - Baggage在全局范围内，（伴随业务系统的调用）跨进程传输数据。Span的tag不会进行传输，因为他们不会被子级的span继承。
> - span的tag可以用来记录业务相关的数据，并存储于追踪系统中。实现OpenTracing时，可以选择是否存储Baggage中的非业务数据，OpenTracing标准不强制要求实现此特性。
>
> 





## 跨进程追踪

`Inject` 和 `Extract` 允许开发者进行跨进程追踪时，不用和特定的OpenTracing实现进行紧耦合

- Inject, Extract, 和 Carriers

> #### Carrier格式
>
> 所有的Carrier都有自己的格式。在一些语言的OpenTracing实现中，格式必须必须作为一个常量或者字符串来指定； 另一些，则通过Carrier的静态类型来指定。



## Inject/Extract Carrier 所必须的格式

> OpenTracing标准所有平台的实现者支持两种Carrier格式：基于"text map"（基于字符串的map）的格式和基于"binary"（二进制）的格式。
>
> - *text map* 格式的 Carrier是一个平台惯用的map格式，基于unicode编码的`字符串`对`字符串`键值对
> - *binary* 格式的 Carrier 是一个不透明的二进制数组（可能更紧凑和有效）
>
> 



