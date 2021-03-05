go-http-api-sample
==================

这是一份 golang http api 应用样例。  

WARNING: 项目当前状态**未经测试、细节粗糙，不能用于生产环境**，仅供设计思路参考！  

已完成：  

- 架构设计
- 错误码 与 response 设计
- logger、recovery
- API：
    - /sms/login -- 登录短信验证码（90%）

todo:  

- 单元测试代码
- token 身份认证模块
- rbac 鉴权模块
- 标准 crud 例子
- swagger 接入
- Dockfile 编写
- ... 等

项目布局
--------

本项目目录结构参考了 [golang 标准项目布局](https://github.com/golang-standards/project-layout)，但对部分文件夹定义了自己的语义。    

由于项目是一个直接面向用户的应用程序，而非为其他包调用而设计。故舍弃/修改了部分文件夹语义：  

- `internal`：舍弃。整个项目均私有，不供外部调用。这样减少了项目目录的层级。
- `pkg`：修改。当前意为：“公共包目录”，不供外部调用。

项目布局规划如下（假定项目包含两个可执行文件：`app`、`console`）：  

``` 
.
├── app                     # 应用程序 app 目录
│   ├── config              # 配置文件
│   ├── dao                 # dao 层
│   ├── handler             # handler 层，包含逻辑上的控制器和中间件
│   │   ├── pkg             # handler 内部公共包。ps：其他包内也可定义 pkg，语义相同。
│   ├── model               # gorm model 定义
│   ├── pkg                 # app 内部公共包
│   ├── service             # service 层
│   ├── test                # 测试相关
│   │   ├── helper          # 测试公共助手函数
│   │   ├── mock            # gomock 等工具生成的 `.go` 格式 mock 文件
│   │   ├── testdata        # 公共 testdata
│   ├── Makefile            # app 程序的 Makefile 文件
├── build                   # 打包和持续集成，包含 Dockfile 等
│   ├── app                 # 对应 app 程序
│   ├── console             # 对应 console 程序
├── cmd                     # main 目录
│   ├── app                 # 对应 app 程序
│   ├── console             # 对应 console 程序
├── console                 # 应用程序 console 目录
├── deployments             # 容器编排
├── docs                    # 文档
├── scripts                 # 脚本，以减轻 Makefile 文件复杂度
├── Makefile                # 包含所有应用程序的 Makefile 文件
```

应用程序 app 的一些具体设计：  

```
/app
├── config              
│   ├── secret.yaml     # 机密配置文件 
│   ├── template.yaml   # 通用配置文件
├── handler                     # 所有 gin.HandlerFunc，包含逻辑上的控制器与中间件
│   ├── pkg
│   │   ├── e                   # 业务错误码（参照谷歌API设计指南而设计）
│   │   │   │── code.go         # 业务错误码及映射定义
│   │   │   │── error_detail.go # 定义具体错误
│   │   ├── ginvalidator        # 用于初始化 gin 内部 validator，包含自定义验证器、翻译等
│   │   │   │── ... ... ...
│   ├── ctrl_sms_login.go       # 登录验证码控制器
│   ├── ... ... ...             # 其他控制器（文件命令统一使用 ctrl 前缀）
│   ├── handler.go              # handler 通用函数
│   ├── mw_logger.go            # http 日志中间件
│   ├── mw_recovery.go          # recovery 中间件
│   ├── mw_authentication.go    # 鉴权中间件
│   ├── mw_authorization.go     # 身份认证中间件
│   ├── ... ... ...             # 其他中间件（文件命令统一使用 mw 前缀）
├── pkg                         # /app 下的通用包
│   ├── cache                   # 各种 cache 实例
│   │   ├── go_cache.go         # 外部库 go_cache 实例
│   │   ├── go_redis.go         # 外部库 go_redis 实例
│   ├── config                  # 配置模块
│   │   ├── testdata            # 配置模块下的 testdata（非共用）
│   │   │   │── config.yaml     
│   │   │   │── sercet.yaml
│   │   ├── config.go           # 使用 viper 实现配置文件读取
│   │   ├── config_test.go
│   ├── sms                     # 短信验证码模块的接口定义及实现
│   │   ├── aliyun.go           # 阿里云实现
│   │   ├── tencent.go          # 腾讯云实现
│   │   ├── interface.go        # 接口定义
│   ├── util                    # util 包（个人认为将包命名为 util 是可取的）
│   │   ├── rand.go             # 生成随机值系列函数
│   │   ├── rand_test.go        
├── service                     # serice 层
│   ├── interface.go            # handler 中使用的接口定义在此
│   ├── service_login_sms.go    # 登录短信验证码 service
│   ├── ... ... ...
├── ... ... ...
├── app.go                      # 实例化 app
```

业务错误设计（面向客户端而非日志）
--------------

业务错误设计参照了 “谷歌API设计指南” 中的错误设计。[See more](https://www.bookstack.cn/read/API-design-guide/API-design-guide-07-%E9%94%99%E8%AF%AF.md)    

### 错误码

- 使用了 谷歌标准错误码 和 标准错误详情类型
- 为每个错误码定义了一组固定的映射信息
- 通过 handler 层的公共方法 `success()` `fail()` 对错误返回进行了封装

具体参见：  

- app/handler/e 模块
- app/handler/handler.go 文件

一个典型的错误响应示例：  

```json
// 该错误映射 http status 400
{
    "code": 3,  
    "status": "INVALID_ARGUMENT",
    "message": "客户端指定了无效的参数",
    "error": [
        {
            "field_violations": [
                {
                    "field": "cn_cell_phone_number",
                    "description": "cn_cell_phone_number为必填字段"
                }
            ]
        }
    ]
}
```

程序中应当尽量使用`谷歌标准错误码`和`标准错误详情类型`，而非自行定义，避免破窗效应。  

对于不在`标准错误码`之内的业务错误，例如：“超过30分钟未支付订单，订单已自动取消，无法继续支付”，
我们可以自定义错误码及映射信息如下：  

```
const (
    CodeOK iota
    ...
    CodePayTimeExpired
)

var codeDetails map[Code]CodeDetail = map[Code]CodeDetail{
    ...
    ...
    CodePayTimeExpired: {
		Code:       CodePayTimeExpired,
		Status:     "PAY_TIME_EXPIRED",
		Message:    "由于长时间未支付订单，订单已自动取消，无法继续支付",
		HttpStatus: http.StatusOK, // 412 Precondition Failed 客户端请求信息的先决条件错误
		LogLevel:   zapcore.WarnLevel,
    }
}
```

层级设计
------------

### handler

handler 层对应 gin.HandlerFunc 类型的实现。  

主要包含以下设计规划：  

- app/handler/pkg 表示公共包，但仅应由 handler 域内使用，不应由 handler 之外的其他包调用。
- controller 文件统一设置前缀为 `ctrl`
- middleware 文件统一设置前缀为 `mw`
- 对于 handler 中注入的 service 接口，接口定义位置视情况而定，定义在 handler 包或 service 包均可。

### service、dao

略  

response 封装 && log
---------------------

- handler/handler.go 中的 success() failed() mustBind() 方法实现 response 封装
- success()、failed()、mustBind() 方法 + logger 中间件，共同实现了 http log，记录 http 请求信息、响应信息、错误。

#### response 封装：

略，详见代码注释

#### log

由于在 gin 中间件中获取到的 response body 是一个 []byte 类型，无法直接使用 zap 记录日志。  

所以本项目在封装的 response 函数中使用 c.Set() 向 gin.Context 附加了日志信息，诸如：
原始错误、日志级别（ps：日志级别已在错误码定义中做了映射）。  

同时在 logger 中间件中读取日志信息，完成日志记录。  

ps：当前日志中间件使用的 zap 实例仅为示意，可以自定义将日志记录到本地文件或远端。  

wire 依赖注入
------------

项目使用 wire 自动生成注入函数。以下是一些 wire 使用约定：  

- 避免在多个文件中定义 provider set。一个应用程序应只在 main 包中定义一次 provider set。

测试
-----

- 单元测试推荐使用库 `testify/assert`+`httpexpect`+`gomock` 组合实现。  
- 对于非公共 testdata，一律定义在测试文件所在包
- 对于 mock 文件，放在 test/mock 文件夹中，mock 文件与目标文件路径相互映射。 

test 包的结构布局，例子如下：  

``` 
├── test                                # 所有测试相关公共文件
│   ├── helper                          # 公共测试助手函数
│   ├── mock        
│   │   ├── handler          
│   │   │   │── ctrl_user_mock.go       # 对应 handler/ctrl_user.go
│   │   ├── service   
│   │   │   │── service_order_mock.go   # 对应 service/service_order.go
│   ├── testdata                        # 公共 testdata
```
















