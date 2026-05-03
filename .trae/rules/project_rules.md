# Jukebox Lite 项目规则

## 项目概述
点歌台应用，微信小程序前端 + Go 后端，用户可搜索/点歌，歌手可接单演唱。

## 技术栈
- **前端**: 微信小程序原生开发（JS + WXML + WXSS），无框架，ES5 语法
- **后端**: Go 1.21 + Gin v1.9，内存存储（无数据库）

## 目录结构
```
miniprogram/          # 微信小程序前端
  app.js/json/wxss    # 入口，globalData.apiBase/apiBase.userInfo
  pages/
    index/            # 首页：歌曲搜索/分类/歌手端订单管理
    login/            # 登录：wx.login → 后端换 token
    confirm/          # 确认下单：选歌手+留言
    success/          # 下单成功
  utils/
    api.js            # HTTP 请求封装，Bearer token 鉴权，401 自动跳登录
backend/              # Go 后端
  main.go             # 入口，手动依赖注入组装各层
  config/config.go    # 全局配置（端口/音乐API/微信凭证）
  model/model.go      # 数据模型：Song/Category/Order/Singer/User + 请求参数结构体
  repository/         # 存储层接口 + 内存实现
    interfaces.go     # Order/Singer/User 三个 Repository 接口
    memory.go         # 内存实现，sync.RWMutex 并发安全，Singer 有种子数据
  service/            # 业务逻辑层
    auth.go           # 微信 code2Session 登录，HMAC-SHA256 token 生成/验证
    song.go           # 歌曲搜索/分类/推荐，委托 MusicAPIService
    order.go          # 订单 CRUD
    singer.go         # 歌手 CRUD
    music_api.go      # 外部音乐 API 代理（sayqz.com），失败时返回 mock 数据
  controller/         # HTTP 处理层，绑定参数→调 service→返回统一 APIResponse
  middleware/
    auth.go           # Bearer token 鉴权中间件，设置 open_id/user/user_nick/user_role
    cors.go           # CORS 跨域中间件
  router/router.go    # 路由注册，/api 前缀，公开/鉴权分组
```

## API 路由
**公开**: POST /api/auth/login | GET /api/songs/search,categories,category,info | GET /api/singers,/:id
**鉴权**: GET /api/auth/profile | PUT /api/auth/role | POST /api/orders | GET /api/orders/:id,/singer,/user | PUT /api/orders/:id/status | POST /api/singers

## 统一响应格式
```json
{"code": 0, "message": "success", "data": ...}
```
code=0 成功，非0 失败；401 时前端自动清 token 跳登录页。

## 前端约定
- 页面四件套：.js/.json/.wxml/.wxss
- API 调用统一通过 `utils/api.js`，token 存 wx.StorageSync
- 角色双端：user（搜索点歌）/ singer（接单管理），首页按角色切换视图
- 音乐源：netease/kuwo/qq，picker 切换

## 后端约定
- 分层：router → controller → service → repository，严格单向依赖
- 依赖注入：main.go 手动构造，New* 构造函数
- 分页默认 page=1, limit=20
- 订单状态流转：pending → singing → done
- 用户角色：user / singer
- 开发模式：微信 AppID/AppSecret 为空时，code 直接映射为 dev_{code}

## 编码规范
- 前端：ES5（var/函数表达式），无箭头函数/模板字符串
- 后端：标准 Go 风格，错误返回而非 panic
- 不添加注释
