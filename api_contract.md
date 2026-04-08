# Wails API 契约

> 本文档定义 Go 后端通过 Wails 暴露给 React 前端的全部方法（Binding）与事件（Events）。
> 任何接口变更须同步更新此文档。

---

## 约定

- **请求结构体**：以 `Req` 结尾，必填字段用 `json:"name,required"`  tagging（Go 侧校验在 Service 层）。
- **响应结构体**：统一使用 `Result<T>` 包装（需前端配合），或显式返回 `(*T, error)`，Wails 会将 `error` 透传给前端。
- **事件命名**：采用 `domain:action` 小写驼峰/冒号风格，如 `reminder:triggered`。

---

## 1. 工作区 (Workspace)

### 1.1 获取所有工作区
```go
func (a *App) WorkspaceList() ([]models.Workspace, error)
```
**返回**：工作区数组，按 `created_at` 升序。

### 1.2 创建工作区
```go
func (a *App) WorkspaceCreate(req requests.WorkspaceCreateReq) (*models.Workspace, error)
```
| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 工作区名称，1~50 字符 |

### 1.3 更新工作区
```go
func (a *App) WorkspaceUpdate(req requests.WorkspaceUpdateReq) (*models.Workspace, error)
```
| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 工作区 ID |
| `name` | string | 新名称 |

### 1.4 删除工作区
```go
func (a *App) WorkspaceDelete(id int64) error
```
**行为**：删除工作区及其下所有任务（后续可扩展为迁移）。

---

## 2. 任务 (Task)

### 2.1 获取任务列表
```go
func (a *App) TaskList(workspaceID int64, filter string) ([]models.Task, error)
```
| 参数 | 说明 |
|------|------|
| `workspaceID` | 工作区 ID |
| `filter` | 筛选条件：`all` / `pending` / `completed` |

### 2.2 获取单个任务
```go
func (a *App) TaskGet(id int64) (*models.Task, error)
```

### 2.3 创建任务
```go
func (a *App) TaskCreate(req requests.TaskCreateReq) (*models.Task, error)
```
**注意**：根据 `type` 字段校验不同字段：
- `type=todo`：`title`、`workspace_id` 必填，`due_at` 可选。
- `type=reminder`：`title`、`workspace_id`、`start_time`、`end_time`、`interval_value`、`interval_unit`、`repeat_mode` 必填。

### 2.4 更新任务
```go
func (a *App) TaskUpdate(req requests.TaskUpdateReq) (*models.Task, error)
```

### 2.5 删除任务
```go
func (a *App) TaskDelete(id int64) error
```

### 2.6 切换完成状态
```go
func (a *App) TaskToggleCompleted(id int64) (*models.Task, error)
```
**行为**：对普通待办翻转 `is_completed`；对定时提醒任务无意义（前端可隐藏按钮）。

### 2.7 暂停 / 恢复今日提醒
```go
func (a *App) TaskTogglePause(id int64) (*models.Task, error)
```
**行为**：
- 若今日未暂停，则写入 `paused_date = today`，调度器今日跳过。
- 若今日已暂停，则清空 `paused_date`，恢复触发。

---

## 3. 设置 (Settings)

### 3.1 获取全部设置
```go
func (a *App) SettingsGet() (map[string]string, error)
```
返回一个 KV map，如 `{"default_workspace_id":"1","default_sort":"created_at_desc"}`。

### 3.2 更新设置
```go
func (a *App) SettingsUpdate(req requests.SettingsUpdateReq) error
```
| 字段 | 类型 | 说明 |
|------|------|------|
| `defaultWorkspaceID` | int64 | 默认打开的工作区 |
| `defaultSort` | string | 默认排序方式 |

### 3.3 导出数据
```go
func (a *App) DataExport() (string, error)
```
**返回**：SQLite 数据库文件的 Base64 字符串（或临时文件路径）。

### 3.4 导入数据
```go
func (a *App) DataImport(fileData string) error
```
**参数**：Base64 编码的 SQLite 数据库文件内容。

### 3.5 清除已完成任务
```go
func (a *App) DataClearCompleted() error
```

---

## 4. 事件 (Events)

调度器通过 Wails `runtime.EventsEmit` 向前端推送事件。

### 4.1 提醒触发
```
事件名：reminder:triggered
Payload：models.Task（JSON）
```
**时机**：某个定时提醒任务到达触发时间，且今日未被暂停。

### 4.2 提醒关闭（用户操作反馈）
前端主动调用 Go 方法即可，无需事件；若后续需要后端知道用户已查看，可追加 `AcknowledgeReminder(taskID)` 方法。

---

## 附录：前端调用示例

```js
import { EventsOn } from "../../wailsjs/runtime";
import {
  WorkspaceList,
  TaskList,
  TaskTogglePause,
} from "../../wailsjs/go/app/App";

// Call Go method
const workspaces = await WorkspaceList();

// Listen event
EventsOn("reminder:triggered", (task) => {
  showReminderModal(task);
});
```
