-- 001_init.sql
-- 初始化数据库：工作区、任务表及索引

CREATE TABLE IF NOT EXISTS workspaces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    type TEXT DEFAULT 'todo' CHECK(type IN ('todo', 'reminder')),
    due_at DATETIME,
    remind_at DATETIME,
    is_completed BOOLEAN DEFAULT 0,
    start_time TEXT,
    end_time TEXT,
    interval_value INTEGER,
    interval_unit TEXT CHECK(interval_unit IN ('minute', 'hour')),
    repeat_mode TEXT CHECK(repeat_mode IN ('daily', 'workday', 'weekly', 'monthly')),
    weekdays TEXT,
    month_day INTEGER,
    next_trigger_at DATETIME,
    paused_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tasks_workspace_id ON tasks(workspace_id);
CREATE INDEX IF NOT EXISTS idx_tasks_is_completed ON tasks(is_completed);
CREATE INDEX IF NOT EXISTS idx_tasks_next_trigger_at ON tasks(next_trigger_at);

-- 默认工作区（MVP 预设）
INSERT INTO workspaces (name) VALUES
    ('家庭'),
    ('工作'),
    ('个人');
