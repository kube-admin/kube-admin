.PHONY: help dev backend frontend build test test-backend test-frontend docker docker-up docker-down docker-logs k8s-apply helm-install helm-template clean docs-install docs-dev docs-build docs-preview

help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

dev: ## 启动后端与前端开发服务（后端前台，前端需另开终端或改用 dev-bg）
	cd backend && go run cmd/main.go

backend: ## 启动后端
	cd backend && go run cmd/main.go

backend-build: ## 编译后端二进制到 backend/bin
	cd backend && go build -o bin/server cmd/main.go

frontend: ## 启动前端开发服务
	cd frontend && npm run dev

frontend-build: ## 构建前端生产包
	cd frontend && npm run build

build: backend-build frontend-build ## 构建前后端

test: test-backend test-frontend ## 运行全部测试

test-backend: ## 运行后端测试
	cd backend && go test ./...

test-frontend: ## 运行前端测试
	cd frontend && npm run test

docker: docker-up ## 通过 deploy.sh 构建并启动（= docker-up）

docker-up: ## 构建并启动（含密钥校验 + 健康等待）
	bash deploy/deploy.sh up

docker-down: ## 停止服务（保留数据卷）
	bash deploy/deploy.sh down

docker-logs: ## 跟踪日志（可加 backend/frontend）
	bash deploy/deploy.sh logs

k8s-apply: ## 应用 K8s 原生清单（需先编辑 deploy/k8s/secret.yaml 替换密钥）
	kubectl apply -f deploy/k8s/

helm-install: ## Helm 安装（密钥用强随机值覆盖）
	helm install kube-admin deploy/helm/kube-admin --create-namespace -n kube-admin \
		--set secrets.jwtSecret="$$(openssl rand -base64 32)" \
		--set secrets.encryptKey="$$(openssl rand -base64 32)"

helm-template: ## Helm 渲染验证（不安装）
	helm template kube-admin deploy/helm/kube-admin -n kube-admin

clean: ## 清理构建产物
	rm -rf backend/bin frontend/dist docs-site/.vitepress/dist docs-site/.vitepress/cache

docs-install: ## 安装文档站依赖
	cd docs-site && npm install

docs-dev: ## 启动文档站本地开发服务
	cd docs-site && npm run docs:dev

docs-build: ## 构建文档站静态产物
	cd docs-site && npm run docs:build

docs-preview: ## 预览文档站构建产物
	cd docs-site && npm run docs:preview
