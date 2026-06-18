---
name: dever-source
description: Use when 修改 Dever source 组件，包括资源中心、资源、频道、分类、来源、资源状态、source page JSON、Provider hook、选项加载、权限和迁移行为。
version: 0.1.0
---

# Source 组件

本组件 skill 必须和 `shemic-dever` 一起使用。先遵守 Dever 框架规则，再按这里的 source 组件边界修改。

## 事实来源

- 组件源码：`backend/package/source`
- 组件声明：`backend/package/source/dever.json`
- Model：`model`
- Provider hook：`service/hook.go`
- 后台页面：`front/page`

## 硬规则

- 资源、频道、分类、来源、状态这些普通后台维护页优先使用 `Model + package/front + page JSON`。
- 不为普通 CRUD 新增 API 或 Service。
- 校验、归一化、保存前后生命周期统一放在 `service.SourceHook`，不要散落到页面 action。
- 不手改生成文件、编译产物或项目级菜单来补 source 菜单。
- 不在 `dever.json` 写 `apiRoots`；API 扫描由 Dever 按组件自动处理。
- 菜单分组归属在 `dever.json`：`source-center`、`source-resource`、`source-config`。

## Page 规则

- 后台页面路径继续放在 `front/page/admin/...`。
- 标准 list/update 页应复用 front 自动推导 model 和 action。
- 左分类右列表只在右侧列表需要刷新时刷新数据，不为分类切换强行跳转 URL。
- Options 和 Relations 优先从 model comment、Options、Relations 生成，不在 page JSON 重复写死。

## Service 规则

- `service/hook.go` 只承担 source 组件保存生命周期、字段规范化和 option 支撑。
- 需要跨表一致性或保存前校验时扩展现有 hook。
- 不新增空 passthrough Provider，不新增只包一层 model CRUD 的 Service。

## 常见检查

- 菜单不显示：先检查 `dever.json`、`module/source/main.go` 和生成后的组件注册。
- option 报错：先检查 model Options/Relations 和页面上下文，不要硬编码模型名。
- 权限报错：先检查 page action 是否走 front 标准权限，不要临时放开通配权限。
