# 技术架构

## 1. 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端 | Go | 业务逻辑、定时调度、数据库操作 |
| 桌面框架 | Wails v2 | 跨平台桌面应用框架，以轻量方式将 Go 后端与 Web 前端结合 |
| 前端 UI | React 18 | 用户交互界面，使用函数组件 + Hooks |
| 样式 | CSS Modules / Tailwind CSS（可选） | 保持组件样式隔离与灵活 |
| 数据存储 | SQLite | 本地嵌入式数据库，无需额外服务 |

## 2. 架构图

```
┌─────────────────────────────────────┐
│              React 前端              │
│  ┌─────────┐ ┌─────────┐ ┌────────┐ │
│  │ 任务列表 │ │ 提醒弹窗 │ │ 设置页  │ │
│  └─────────┘ └─────────┘ └────────┘ │
└────────────┬────────────────────────┘
             │ Wails 绑定 (JS -> Go)
┌────────────▼────────────────────────┐
│              Go 后端                 │
│  ┌─────────┐ ┌─────────┐ ┌────────┐ │
│  │ 任务服务 │ │ 调度器   │ │ 设置服务│ │
│  └────┬────┘ └────┬────┘ └────┬───┘ │
│       └───────────┴───────────┘     │
│              SQLite (本地)           │
└─────────────────────────────────────┘
```

## 3. 离线优先与同步扩展

- **离线优先**：所有任务和提醒规则均存储在本地 SQLite 中，应用在无网络环境下完全可用。
- **未来同步**：待后端服务搭建后，Go 后端将新增同步模块，通过 REST API 将本地 SQLite 数据与云端进行增量同步，实现多平台数据一致。

## 4. 提醒机制

- Go 后端维护一个调度器（Scheduler），轮询数据库中的提醒规则。
- 当某个提醒到期时，调度器通过 Wails 事件（Events）通知前端 React。
- React 接收到事件后，弹出全局提醒弹窗。弹窗显示任务标题与操作按钮（"知道了" / "稍后提醒"）。

## 5. 数据库设计

### 5.1 表结构

#### `workspaces` - 工作区表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | 工作区唯一 ID |
| `name` | TEXT NOT NULL | 工作区名称 |
| `created_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 更新时间 |

#### `tasks` - 任务表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | 任务唯一 ID |
| `workspace_id` | INTEGER NOT NULL | 所属工作区 ID（外键） |
| `title` | TEXT NOT NULL | 任务标题 |
| `description` | TEXT | 任务描述 |
| `type` | TEXT DEFAULT 'todo' | 任务类型：`todo` 待办 / `reminder` 定时提醒 |
| `due_at` | DATETIME | 待办任务的截止时间（可选） |
| `remind_at` | DATETIME | 待办任务的提醒时间（可选） |
| `is_completed` | BOOLEAN DEFAULT 0 | 是否已完成 |
| `start_time` | TEXT | 定时提醒开始时间（如 `09:00`） |
| `end_time` | TEXT | 定时提醒结束时间（如 `18:00`） |
| `interval_value` | INTEGER | 提醒周期数值（如 30） |
| `interval_unit` | TEXT | 提醒周期单位：`minute`, `hour` |
| `repeat_mode` | TEXT | 重复模式：`daily` 每天 / `workday` 每个工作日 / `weekly` 每周 / `monthly` 每月 |
| `weekdays` | TEXT | 每周模式时存储选中的星期几，格式如 `1,3,5`（周一、三、五） |
| `month_day` | INTEGER | 每月模式时存储具体的日期（1~31） |
| `next_trigger_at` | DATETIME | 下次触发时间（调度器用） |
| `paused_date` | DATE | 当日暂停日期（YYYY-MM-DD），调度器遇到该日期跳过触发，次日自动恢复 |
| `created_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 更新时间 |

### 5.2 索引

```sql
CREATE INDEX idx_tasks_workspace_id ON tasks(workspace_id);
CREATE INDEX idx_tasks_is_completed ON tasks(is_completed);
CREATE INDEX idx_tasks_next_trigger_at ON tasks(next_trigger_at);
```

## 6. UI 页面结构

### 6.1 页面 / 路由规划

| 页面 | 路径 | 说明 |
|------|------|------|
| 主布局 | `/` | 包含侧边栏、顶部栏、内容区 |
| 任务列表 | `/workspace/:id` | 显示某个工作区下的任务列表 |
| 新建/编辑任务 | `/task/new` 或 `/task/:id/edit` | 表单页，支持设置提醒规则 |
| 设置 | `/settings` | 应用设置、数据导入导出 |

### 6.2 关键组件

- **Sidebar 侧边栏**：工作区列表 + 新建工作区按钮
- **TaskList 任务列表**：任务卡片 + 筛选/排序控件 + 新建任务入口
- **TaskForm 任务表单**：标题、描述、截止时间、提醒规则设置
- **ReminderModal 提醒弹窗**：全局弹窗，到期时自动弹出
- **SettingsPage 设置页**：配置项与数据管理

## 7. 项目目录结构建议

```
whendo/
├── app/                   # Wails 应用入口包
│   ├── app.go             # 应用结构体与 Wails 绑定方法
│   ├── menu.go            # 应用菜单配置（扩展）
│   └── tray.go            # 系统托盘相关（扩展）
├── main.go                # 程序主入口
├── go.mod / go.sum        # Go 依赖
├── wails.json             # Wails 配置文件
├── frontend/              # React 前端目录
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   ├── src/
│   │   ├── main.jsx       # React 挂载入口
│   │   ├── App.jsx        # 根组件
│   │   ├── components/    # 通用 UI 组件
│   │   ├── pages/         # 页面级组件
│   │   ├── hooks/         # 自定义 Hooks
│   │   ├── services/      # 与 Go 后端交互的 API 封装
│   │   └── styles/        # 全局样式
│   └── public/
├── internal/              # Go 内部模块
│   ├── database/          # SQLite 连接与迁移
│   ├── models/            # 数据模型定义
│   ├── services/          # 业务逻辑（任务、提醒、工作区）
│   ├── scheduler/         # 定时调度器
│   └── utils/             # 工具函数
├── build/                 # 构建产物与打包配置
└── docs/                  # 文档目录
    └── 项目文档.md
```

## 8. 开发环境准备

### 8.1 前置依赖
- Go 1.21+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### 8.2 初始化命令
```bash
# 进入项目目录后初始化 Wails + React 项目
wails init -n whendo -t react
```

## 9. 术语表

| 术语 | 说明 |
|------|------|
| Wails | Go 编写的轻量级跨平台桌面应用框架 |
| 工作区 (Workspace) | 任务分类的顶层容器 |
| 单次提醒 | 在唯一指定时间点触发一次 |
| 循环提醒 | 按照固定周期重复触发 |
| Snooze | 用户选择延后提醒 |
