#!/usr/bin/env bash
# Kube Admin 统一部署脚本（Docker Compose 单机部署的薄封装）。
# 命令：up | down | restart | logs [backend|frontend] | build | ps
# 选项：--dev 跳过密钥强校验（仅本地开发）
#
# macOS / Linux bash 兼容；依赖 docker（compose v2 优先，v1 回退）与 curl。
set -euo pipefail

# ----- 定位目录（不依赖 readlink -f，兼容 macOS） -----
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/docker/docker-compose.yaml"
ENV_FILE="$SCRIPT_DIR/.env"
BACKEND_PORT="${PORT:-8080}"

# ----- 颜色（tput 跨平台，无 tput 则降级） -----
if command -v tput >/dev/null 2>&1 && [ -t 1 ]; then
  RED="$(tput setaf 1)"; GREEN="$(tput setaf 2)"; YELLOW="$(tput setaf 3)"; BOLD="$(tput bold)"; RESET="$(tput sgr0)"
else
  RED=""; GREEN=""; YELLOW=""; BOLD=""; RESET=""
fi

info()  { printf "%s==>%s %s\n" "$GREEN" "$RESET" "$*"; }
warn()  { printf "%s⚠ %s%s\n" "$YELLOW" "$*" "$RESET"; }
err()   { printf "%s✗ %s%s\n" "$RED" "$*" "$RESET" >&2; }

# ----- 选择 docker compose 命令 -----
detect_compose() {
  if docker compose version >/dev/null 2>&1; then
    COMPOSE=(docker compose)
  elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE=(docker-compose)
  else
    err "未找到 docker compose（v2）或 docker-compose（v1），请先安装。"
    exit 1
  fi
}

# ----- 密钥校验 -----
check_env() {
  local dev_mode=0
  for arg in "$@"; do [ "$arg" = "--dev" ] && dev_mode=1; done

  if [ ! -f "$ENV_FILE" ]; then
    err "未找到 $ENV_FILE，请先：cp .env.example .env 并填写密钥"
    exit 1
  fi

  # shellcheck disable=SC1090
  set -a; . "$ENV_FILE"; set +a

  local missing=0 weak=0
  [ -z "${JWT_SECRET:-}" ] && { err "JWT_SECRET 为空"; missing=1; }
  [ -z "${ENCRYPT_KEY:-}" ] && { err "ENCRYPT_KEY 为空"; missing=1; }
  case "${JWT_SECRET:-}" in  dev-only-*|your-secret-key-change-in-production|"") warn "JWT_SECRET 仍是弱默认值"; weak=1;; esac
  case "${ENCRYPT_KEY:-}" in  dev-only-*|"") warn "ENCRYPT_KEY 仍是弱默认值"; weak=1;; esac

  if [ "$missing" = "1" ]; then
    err "密钥缺失，拒绝启动。请在 $ENV_FILE 设置 JWT_SECRET / ENCRYPT_KEY"
    err "生成强随机值：openssl rand -base64 32"
    exit 1
  fi
  if [ "$weak" = "1" ] && [ "$dev_mode" != "1" ]; then
    err "检测到弱密钥且未使用 --dev，拒绝启动。本地测试请加 --dev"
    exit 1
  fi
}

# ----- 健康等待 -----
wait_health() {
  info "等待后端健康检查 http://localhost:$BACKEND_PORT/healthz ..."
  local i
  for i in $(seq 1 60); do
    if curl -sf "http://localhost:$BACKEND_PORT/healthz" >/dev/null 2>&1; then
      printf "%s✓ 后端已就绪（耗时约 %ss）%s\n" "$GREEN" "$i" "$RESET"
      printf "\n  %s前端：%shttp://localhost%s\n" "$BOLD" "$RESET" "${BACKEND_PORT:+ ($BACKEND_PORT 后端 / 80 前端)}"
      printf "  %s后端：%shttp://localhost:%s/healthz%s\n" "$BOLD" "$RESET" "$BACKEND_PORT" "$RESET"
      printf "  %s登录：%sadmin / admin123（首次登录后请立即改密）%s\n\n" "$BOLD" "$RESET" "$RESET"
      return 0
    fi
    sleep 1
  done
  warn "后端 60s 内未就绪，最近日志："
  "${COMPOSE[@]}" -f "$COMPOSE_FILE" logs --tail=30 backend || true
  return 1
}

# ----- 命令分发 -----
usage() {
  cat <<EOF
用法: bash deploy.sh [--dev] <命令>

命令:
  up                构建并启动（前台构建，后台运行），等待健康检查
  down              停止并删除容器（保留数据卷）
  restart           重启服务
  logs [service]    跟踪日志（service: backend | frontend，默认全部）
  build             仅构建镜像，不启动
  ps                查看容器状态

选项:
  --dev             跳过弱密钥强校验（仅本地开发）
EOF
}

main() {
  detect_compose
  case "${1:-}" in
    up)
      check_env "$@"
      info "构建并启动（compose file: $COMPOSE_FILE）"
      "${COMPOSE[@]}" -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --build
      wait_health
      ;;
    down)
      info "停止服务"
      "${COMPOSE[@]}" -f "$COMPOSE_FILE" down
      ;;
    restart)
      info "重启服务"
      "${COMPOSE[@]}" -f "$COMPOSE_FILE" restart
      ;;
    logs)
      shift || true
      "${COMPOSE[@]}" -f "$COMPOSE_FILE" logs -f "${1:-}" || "${COMPOSE[@]}" -f "$COMPOSE_FILE" logs -f
      ;;
    build)
      check_env "$@"
      "${COMPOSE[@]}" -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build
      ;;
    ps)
      "${COMPOSE[@]}" -f "$COMPOSE_FILE" ps
      ;;
    -h|--help|""|"help")
      usage
      ;;
    *)
      err "未知命令: $1"
      usage
      exit 1
      ;;
  esac
}

main "$@"
