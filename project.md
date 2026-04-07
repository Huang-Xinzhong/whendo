# WhenDo - 项目文档

## 1. 项目概述

**WhenDo** 是一款结合 Todo List 与定时提醒的跨平台桌面应用。它帮助用户将待办事项与特定时间点或周期性提醒相结合，并通过"工作区"（如家庭、工作、个人）对任务进行分类管理。

### 1.1 目标用户
- 需要按时间或周期触发提醒的个人用户
- 希望将生活与工作任务分开管理的效率工具使用者

### 1.2 核心价值
- **任务 + 提醒一体化**：不仅是待办清单，更能在指定时间主动提醒。
- **工作区隔离**：不同场景的任务互不干扰。
- **离线优先**：无需网络即可使用，未来可扩展云端同步。

---

## 2. 功能清单

### 2.1 MVP (V0.1) 功能

#### 工作区模块
- [ ] 预设工作区：家庭、工作、个人
- [ ] 创建新工作区
- [ ] 重命名工作区
- [ ] 删除工作区（迁移或清空其下任务）
- [ ] 在工作区之间切换

#### 任务管理模块
- [ ] 添加任务（标题、描述、任务类型）
  - **任务类型**：分为「代办」与「定时提醒」
    - **代办**：需填写截止时间，可额外设置单次提醒时间
    - **定时提醒**：需设置提醒时间范围（开始时间、结束时间）、提醒周期（每 15/30/60/120 分钟）以及重复模式：
      - **每天**
      - **每个工作日**
      - **每周（自定义）**：选择具体的星期几（如周一、周三）
      - **每月（自定义）**：选择具体的日期（如每月 15 号）
- [ ] 查看任务列表（按创建时间/截止时间排序）
- [ ] 编辑任务
- [ ] 删除任务
- [ ] 标记任务为已完成 / 未完成
  - **快捷操作按钮**：任务卡片右侧提供一键操作。普通待办显示「完成」按钮；定时提醒任务显示「暂停」按钮，可临时暂停当日提醒，再次点击恢复，次日自动恢复正常触发。
- [ ] 按状态筛选（全部 / 未完成 / 已完成）

#### 定时提醒模块
- [ ] **定时提醒任务**：设置提醒时间范围（如 09:00 ~ 18:00）、提醒周期（每 15/30/60/120 分钟）以及重复模式（每天 / 每个工作日 / 每周特定几天 / 每月特定日期）
- [ ] 应用内弹窗提醒（倒计时到期后在前端弹出提示）
- [ ] 提醒关联具体任务，点击弹窗可跳转任务详情
- [ ] 关闭或延迟（Snooze）提醒
- [ ] 暂停今日提醒：任务列表右侧暂停按钮可将当日提醒静音，次日自动恢复正常触发

#### 设置模块
- [ ] 默认工作区设置
- [ ] 任务列表默认排序方式
- [ ] 数据导出 / 导入（SQLite 备份）
- [ ] 清除已完成任务

---

## 3. 技术架构

### 3.1 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端 | Go | 业务逻辑、定时调度、数据库操作 |
| 桌面框架 | Wails v2 | 跨平台桌面应用框架，以轻量方式将 Go 后端与 Web 前端结合 |
| 前端 UI | React 18 | 用户交互界面，使用函数组件 + Hooks |
| 样式 | CSS Modules / Tailwind CSS（可选） | 保持组件样式隔离与灵活 |
| 数据存储 | SQLite | 本地嵌入式数据库，无需额外服务 |

### 3.2 架构图

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

### 3.3 离线优先与同步扩展

- **离线优先**：所有任务和提醒规则均存储在本地 SQLite 中，应用在无网络环境下完全可用。
- **未来同步**：待后端服务搭建后，Go 后端将新增同步模块，通过 REST API 将本地 SQLite 数据与云端进行增量同步，实现多平台数据一致。

### 3.4 提醒机制

- Go 后端维护一个调度器（Scheduler），轮询数据库中的提醒规则。
- 当某个提醒到期时，调度器通过 Wails 事件（Events）通知前端 React。
- React 接收到事件后，弹出全局提醒弹窗。弹窗显示任务标题与操作按钮（"知道了" / "稍后提醒"）。

---

## 4. 数据库设计

### 4.1 表结构

#### `workspaces` - 工作区表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | 工作区唯一 ID |
| `name` | TEXT NOT NULL | 工作区名称 |
| `is_preset` | BOOLEAN DEFAULT 0 | 是否为预设工作区（不可删除） |
| `sort_order` | INTEGER DEFAULT 0 | 排序权重 |
| `created_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 更新时间 |

#### `tasks` - 任务表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | 任务唯一 ID |
| `workspace_id` | INTEGER NOT NULL | 所属工作区 ID（外键） |
| `title` | TEXT NOT NULL | 任务标题 |
| `description` | TEXT | 任务描述 |
| `type` | TEXT DEFAULT 'todo' | 任务类型：`todo` 代办 / `reminder` 定时提醒 |
| `due_at` | DATETIME | 代办任务的截止时间（可选） |
| `remind_at` | DATETIME | 代办任务的提醒时间（可选） |
| `is_completed` | BOOLEAN DEFAULT 0 | 是否已完成 |
| `created_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 更新时间 |

#### `reminders` - 提醒规则表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | 提醒唯一 ID |
| `task_id` | INTEGER NOT NULL | 关联任务 ID（外键，级联删除） |
| `start_time` | TEXT | 每日提醒开始时间（如 `09:00`） |
| `end_time` | TEXT | 每日提醒结束时间（如 `18:00`） |
| `interval_value` | INTEGER | 提醒周期数值（如 30） |
| `interval_unit` | TEXT | 提醒周期单位：`minute`, `hour` |
| `repeat_mode` | TEXT | 重复模式：`daily` 每天 / `workday` 每个工作日 / `weekly` 每周 / `monthly` 每月 |
| `weekdays` | TEXT | 每周模式时存储选中的星期几，格式如 `1,3,5`（周一、三、五） |
| `month_day` | INTEGER | 每月模式时存储具体的日期（1~31） |
| `is_enabled` | BOOLEAN DEFAULT 1 | 是否启用 |
| `next_trigger_at` | DATETIME | 下次触发时间（调度器用） |
| `paused_date` | DATE | 当日暂停日期（YYYY-MM-DD），调度器遇到该日期跳过触发，次日自动恢复 |
| `created_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 更新时间 |

### 4.2 索引

```sql
CREATE INDEX idx_tasks_workspace_id ON tasks(workspace_id);
CREATE INDEX idx_tasks_is_completed ON tasks(is_completed);
CREATE INDEX idx_reminders_task_id ON reminders(task_id);
CREATE INDEX idx_reminders_next_trigger_at ON reminders(next_trigger_at);
```

---

## 5. UI 页面结构

### 5.1 页面 / 路由规划

| 页面 | 路径 | 说明 |
|------|------|------|
| 主布局 | `/` | 包含侧边栏、顶部栏、内容区 |
| 任务列表 | `/workspace/:id` | 显示某个工作区下的任务列表 |
| 新建/编辑任务 | `/task/new` 或 `/task/:id/edit` | 表单页，支持设置提醒规则 |
| 设置 | `/settings` | 应用设置、数据导入导出 |

### 5.2 关键组件

- **Sidebar 侧边栏**：工作区列表 + 新建工作区按钮
- **TaskList 任务列表**：任务卡片 + 筛选/排序控件 + 新建任务入口
- **TaskForm 任务表单**：标题、描述、截止时间、提醒规则设置
- **ReminderModal 提醒弹窗**：全局弹窗，到期时自动弹出
- **SettingsPage 设置页**：配置项与数据管理

---

## 6. 里程碑与版本规划

### V0.1 - MVP（最小可用版本）
- 工作区预设 + 自定义管理
- 任务 CRUD
- 单次提醒 + 循环提醒（每 N 分钟/小时/天/周/月）
- 应用内弹窗提醒
- 设置页面（基础配置 + 数据导入导出）
- 平台：Windows / Linux 桌面端

### V0.2 - 增强体验
- 增加声音提醒
- 支持"稍后提醒"（Snooze）自定义时长
- 循环规则进阶：支持工作日/周末跳过
- 任务标签/优先级
- 深色模式

### V0.3 - 多端同步
- 后端服务搭建（Go REST API）
- 用户账号体系
- 本地 SQLite 与云端增量同步
- 移动端适配 / macOS 支持

---

## 7. 项目目录结构建议

```
whendo/
├── app.go                 # Wails 应用入口（Go）
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

---

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

---

## 9. 术语表

| 术语 | 说明 |
|------|------|
| Wails | Go 编写的轻量级跨平台桌面应用框架 |
| 工作区 (Workspace) | 任务分类的顶层容器 |
| 单次提醒 | 在唯一指定时间点触发一次 |
| 循环提醒 | 按照固定周期重复触发 |
| Snooze | 用户选择延后提醒 |
