# Debug History

## 问题记录

### 问题 1: Projects 页面表格表头渲染错误

**发现时间**: 2026-03-13

**问题描述**:
在 Projects 页面的表格中，表头 "负责人" 被错误渲染为 `&lt;负责人>负责人`。

**问题原因**:
在 `templates/projects.html` 第52行，模板语法错误，写成了：
```html
<负责人>负责人</th>
```
应该是：
```html
<th>负责人</th>
```

**测试结果**:
- 首页 `/`: 正常
- 会话页 `/sessions`: 正常
- 任务页 `/tasks`: 正常
- 项目页 `/projects`: **有错误** - 表格表头渲染错误
- 用量页 `/usage`: 正常
- `/health` API: 正常
- `/api/sessions` API: 正常
- `/api/tasks` API: 正常
- `/api/projects` API: 正常

**修复状态**: 已修复

**修复内容**:
将 `templates/projects.html` 第52行的 `<负责人>负责人</th>` 改为 `<th>负责人</th>`

**验证结果**:
- 修复后重新启动服务器
- 访问 `/projects` 页面，确认表头正确显示为 "负责人"
- 所有页面（/、/sessions、/tasks、/projects、/usage）均返回 200 状态码
- 项目页面不再有HTML编码错误
