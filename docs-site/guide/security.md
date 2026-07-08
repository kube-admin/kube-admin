# 用户 / 权限 / 审计

Kube Admin 内置完整的认证、授权与审计能力。

## 用户管理

进入「系统设置 → 用户管理」（仅 `admin` 可见）：

- 创建 / 编辑 / 删除用户
- 修改用户角色与邮箱
- 重置密码（bcrypt 加密存储）

::: warning 默认账号
首次启动自动创建 `admin / admin123`，**请立即修改密码**。
:::

## 角色与权限（RBAC）

| 角色 | 资源读写 | 终端/日志/YAML | 用户/集群管理 | 审计查看 |
|---|---|---|---|---|
| `admin` | ✅ | ✅ | ✅ | ✅ |
| `operator` / `user` | ✅ | ✅ | ❌ | ❌ |
| `viewer`（预留） | 只读 | ❌ | ❌ | ❌ |

权限在两个层面生效：

- **接口层**：`RequireRole` 限制用户/集群/审计管理仅 `admin`；`WriteAuth` 限制所有写操作需 `admin`/`operator`/`user`。
- **界面层**：`v-permission` 指令按角色隐藏管理类按钮（如「创建 Clusters」仅 admin 可见）。

## 审计日志

所有写操作（`POST` / `PUT` / `DELETE` / `PATCH`）都会被记录：

- 字段：操作人、方法、路径、响应状态、客户端 IP、UserAgent、时间
- 仅 `admin` 可在「审计日志」接口（`GET /api/v1/audit/logs`）查看

## 集群凭据安全

- 接入集群的 `Token` 与 `kubeconfig` 内容写入数据库前 **AES-256-GCM 加密**
- API 响应对敏感字段脱敏（仅返回「是否已配置」标志）
- TLS 证书默认校验，仅 `TLS_SKIP_VERIFY=true` 时放行（仅限受控环境）

## 安全建议

1. 生产环境设置强随机 `JWT_SECRET` 与 `ENCRYPT_KEY`
2. 启用 HTTPS（通过反向代理终止 TLS）
3. 定期轮换集群凭据与用户密码
4. 为不同团队分配最小权限角色
5. 关注审计日志中的异常写操作

## 下一步

- [部署指南](./deploy.md)
- [API 速查](./api.md)
