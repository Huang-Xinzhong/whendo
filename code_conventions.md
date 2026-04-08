# 代码规范 (Code Conventions)

## 1. Go 后端

### 1.1 目录与包
- `app/`：仅存放 Wails 生命周期、绑定方法、菜单/托盘。业务逻辑禁止写在这里。
- `internal/services/`：纯业务逻辑，不依赖 Wails runtime，便于单元测试。
- `internal/database/`：仅负责连接、迁移、原始查询封装。
- `internal/models/`：仅放 struct 定义，零业务逻辑。

### 1.2 命名
- 包名：全小写，无下划线（`scheduler` 而非 `task_scheduler`）。
- 接口名：以 `er` 结尾或行为命名，如 `TaskStore`、`ReminderNotifier`。
- 导出方法：PascalCase；未导出：camelCase。
- 错误变量前缀 `Err`，如 `ErrWorkspaceNotFound`。

### 1.3 错误处理
- Service 层返回带上下文的错误：`fmt.Errorf("create task: %w", err)`。
- 不裸抛 SQLite 错误给前端，App 层做统一映射（后续可扩展错误码）。

### 1.4 并发
- Scheduler 运行时必须是 `App` 初始化后的独立 goroutine，使用 `context cancel` 优雅退出。
- 共享状态（如当前暂停日期缓存）使用 `sync.RWMutex` 保护。

### 1.5 单元测试
- Service 层必须可 mock 数据库；优先用 interface + fake 实现，而非直接连 SQLite。

---

## 2. React 前端

### 2.1 目录结构（按职责分层）
```
src/
├── components/          # 可复用 UI 组件层
│   ├── common/          # 纯原子/通用组件，无业务语义（仅供 UI 拼装使用，后续手动调整时统一放此目录）
│   │   ├── button.jsx   # 用于：所有需要按钮的地方（表单、弹窗、工具栏）
│   │   ├── input.jsx    # 用于：表单输入、搜索框
│   │   ├── modal.jsx    # 用于：弹窗容器（提醒弹窗、确认删除弹窗）
│   │   └── card.jsx     # 用于：内容卡片外壳（任务卡片、设置项卡片）
│   └── business/        # 业务组件，绑定领域语义（仅在页面或特定功能模块中使用）
│       ├── task-card.jsx         # 用于：task-list 页面 / 任务相关页面
│       ├── task-form.jsx         # 用于：新建/编辑任务页面
│       ├── reminder-modal.jsx    # 用于：全局提醒触发时由 app.jsx 统一挂载
│       ├── workspace-item.jsx    # 用于：sidebar 工作区列表
│       └── sidebar.jsx           # 用于：main-layout 布局组件
├── layouts/             # 页面布局骨架（仅在路由顶层使用）
│   └── main-layout.jsx  # 用于：app.jsx 路由配置中包裹需要侧边栏+顶部栏的页面
├── pages/               # 路由级页面组件（每个页面对应一个路由，负责组装业务组件与数据流）
│   ├── workspace-page.jsx        # 用于：/workspace/:id 路由
│   ├── task-form-page.jsx        # 用于：/task/new、/task/:id/edit 路由
│   └── settings-page.jsx         # 用于：/settings 路由
├── hooks/               # 自定义 Hooks（仅在组件或页面中调用）
│   ├── use-tasks.js     # 用于：workspace-page、task-form-page
│   ├── use-workspaces.js # 用于：sidebar、workspace-page、全局默认工作区初始化
│   └── use-reminders.js # 用于：app.jsx（全局监听提醒事件）
├── services/            # 对 wailsjs/go 的再封装（仅供 hooks/pages 调用，组件层禁止直接 import）
│   └── api.js
├── stores/              # 全局状态（Context / Zustand），谨慎使用（仅在跨页面时由 hooks 消费）
│   └── workspaceStore.js
└── styles/              # 全局 CSS 变量、主题
    ├── variables.css    # 用于：全局注入（颜色、间距、字体）
    └── global.css       # 用于：main.jsx 全局引入
```

### 2.2 组件封装规范
- **全部使用函数组件 + Hooks**，每个组件只做一件事。
- **common 组件**：接收纯 Props，不调用 hooks/services，保证任意上下文中可复用。
- **business 组件**：可调用 hooks 获取局部数据，但禁止直接调用 `wailsjs/go/...`，数据由父组件或 hooks 注入。
- **pages 组件**：负责数据获取、错误处理和路由跳转，不处理精细 UI 逻辑。
- 组件文件名大驼峰，与默认导出同名，如 `TaskCard.jsx` export `TaskCard`。
- Props 必须简短、语义化；复杂对象用解构 + 默认值。

### 2.3 状态管理
- **优先本地状态**（`useState`），仅在真正需要跨页面共享时才提升到 `stores/`。
- 跨页面共享数据（如当前工作区列表）放 `stores/` 或 Context，由 `hooks/` 封装后供组件消费。
- **禁止在 components/ 目录下的组件里直接调用 `wailsjs/go/...`**，统一走 `services/api.js` 封装层，再由 hooks 调用。

### 2.4 样式规范
- **common 组件**：使用 `*.module.css`，命名 `Button.module.css`。
- **business 组件**：使用 `*.module.css`，命名 `TaskCard.module.css`。
- **pages/layouts**：允许直接使用全局类名或少量行内样式，复杂样式也使用 module。
- 颜色、间距、字体**统一使用 CSS 变量**（`var(--color-primary)`），为后续深色模式做准备。

---

## 3. 提交规范

- 格式：`type(scope): subject`
- 常用 type：`feat`、`fix`、`docs`、`refactor`、`test`、`chore`
- 示例：
  - `feat(scheduler): add next_trigger_at calculator`
  - `fix(ui): reminder modal z-index issue`
