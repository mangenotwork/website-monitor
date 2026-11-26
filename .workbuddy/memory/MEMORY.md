# website-monitor 项目长期记忆

## 项目形态
- 网站监测平台，单仓库两个**独立 Go module**：`master`（Web 管理端 + UDP 服务端，18888/18889）、`monitor`（探针客户端）。无 go.work。
- 技术栈：Gin + Go html/template（定界符 `{[ ]}`）+ Vue3/ ECharts / Bootstrap（CDN 引入）；存储为 **BoltDB 单文件** `master/data.db`；日志按 `log/<hostID>_<YYYYMMDD>.log` 落地，行内 `|` 分隔 19 字段。
- 分布式通信：`github.com/mangenotwork/udp_comm v0.0.1`（自研 UDP 应用层协议，包头 1+7+7 字节，DES-ECB + gzip，data 建议 <533 字节）。
- master 与 monitor 之间是**双通道**：UDP 做通知/拉取/上报，HTTP `/data/*` 做配置下发；`/data` 组无鉴权。

## 代码分层
- master: `routers`(路由+JWT/CSRF 中间件) → `handler`(参数解析与输出) → `dao`(BoltDB 封装 + 采集工具) → `entity`；`business`(定时扫描 Collect、告警规则)；`udpServer`(接收 monitor 上报)。
- monitor: `business`(内存 sync.Map 存网站/URL/监测点 + 每秒 tick 计频监测) → `handler`(UDP 通知/Get 回调)。

## 已知关键问题（迭代优先候选）
1. `dao/db.go` 每次 Get/Set/Delete 都 `bolt.Open` + `Close`，且 `GetDBConn()` 复用并关闭共享 `DB.Conn`，并发下有文件锁阻塞与竞态风险。
2. `WebsiteList` 调 `dao.GetClientList()`，对每个 monitor 同步 UDP Get（ipAddr/osInfo，各 1s 超时）→ 面板加载慢的主因。
3. `master/business/scan.go`：`go func(){ dao.Scan(v.Host,...) }()` 闭包捕获循环变量（module go 1.19，无 loopvar 语义）；`% scan.ScanRate` 为 0 时除零 panic。
4. `dao/whois.go`：`net.Dial` 错误被忽略，conn 为 nil 直接 Write → panic。
5. `dao/monitor_log.go ReadLog`：`os.Open` 失败后仍 `defer f.Close()`，f 为 nil → panic。
6. monitor `PingActive` 硬编码 1000ms，未使用 `PingTime`；`request()` 无超时参数（gathertool 默认 30s），故网站超时配置不生效。
7. monitor 复用同一个 `mLog` 对象，Random/Point 阶段未重置 `UriCode/UriMs`；错误分支误写 `ContrastUriCode=0`。
8. monitor `Business()` 在单 goroutine 顺序执行，慢站点会阻塞所有站点。
9. `dao/alert.go` 的 `AlertTimeInfoMap` 无淘汰机制；`GetAtWebsite` 为 O(n*m) 双层遍历。
10. 安全问题：`app.yaml` 明文账号密码；`/data/*` 无鉴权；DES-ECB 强度不足。

## 迭代约定
- 改动前先确认范围；DAO 沿用 interface + `NewXxx()` 返回风格；保持日志行 19 字段格式兼容（master 解析依赖）。
- 两个 module 需分别 `go build ./...` 验证（已验证均通过，Go 1.25 工具链）。

## 前端架构（2026-11-26 UI 改版后）
- 样式唯一入口 `static/css/theme.css`（Bootstrap 先加载，再覆盖）；`page.css` 已不再引用。设计 token 用 CSS 变量，`html[data-accent]` 支持靛蓝/青绿/石墨，`html[data-density]` 支持舒适/紧凑，偏好存 localStorage。
- 布局：`__header__.html` 同时输出 `.sidebar` + `.topbar` + 确认框/遮罩；页面用 `<main class="viewport">`。图标是 `__header__` 顶部的 SVG sprite（`#i-*`），不再依赖 `SVG` 模板函数（旧弹窗仍可用）。
- 公共 JS：`public.js` 提供 `UI`（toast/抽屉/确认框/强调色/徽标/时间格式化）与 `newSiteDetail()`、`newAlertDetail()`；页面 JS 只保留业务方法。
- 详情统一右侧抽屉（`section/site_drawer.html`、`section/alert_drawer.html`），Bootstrap 仅用于弹窗与下拉；弹窗 z-index 1080 高于抽屉 1070。
- 页面：`/home` 面板、`/sites` 站点（新增）、`/monitor` 探针、`/alert` 报警、`/tool` 左导航工具、`/requester` 三区请求器。
- 密度：`theme.css` 顶部 `--d-*` 令牌成对定义（`:root` 舒适 / `html[data-density="compact"]` 紧凑），覆盖页面留白、行高、表头、卡片、列表项、抽屉、表单、图表高度；改密度只需改这两组值。
- 版本展示：`app.yaml` 的 `app.version` / `app.runType` 由 `handler.ginH` 注入 `.Version` / `.RunType`，`MasterVersion` 仅作兜底；不要再硬编码 "dev"。
- 弹窗：Bootstrap 尺寸类（`modal-lg`/`modal-xl`）必须写在 `.modal-dialog`，写到外层 `.modal` 会通过 `max-width` 压缩 fixed 容器导致弹窗靠左；统一定义 `.modal-dialog` + `.modal-wide` + `modal-dialog-centered`。
- KPI 卡片：四个 `.stat` 是 `<a>`，分别指向 `/sites`、`/sites?seg=monitoring`、`/sites?seg=alert`、`/alert?seg=unread`；两个列表页从 URL 读取 `seg` 并做白名单校验。
- 注意：Bootstrap 的 `.alert` 会污染自定义 `.status.alert`（padding/margin 撑高行），状态修饰类统一用 `danger`，不要用 `alert`。
- 注意：Vue 3 in-DOM 模板会压缩静态换行，静态代码块用 `<br>`；抽屉变量名必须与页面 data 对齐（统一 `site` / `alertDetail`）。

