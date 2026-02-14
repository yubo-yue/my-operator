# my-operator (我的Operator工作区)

一个包含多个operator项目的Kubernetes operator开发工作区，用于开发和测试各种Kubernetes Operators。

## 概述

本仓库作为使用 [Kubebuilder](https://kubebuilder.io/) 开发Kubernetes operators的工作空间。它包含示例operators和本地开发测试的实用工具。

## 前置要求

- [Go](https://golang.org/doc/install) (v1.21+)
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) (Kubernetes in Docker)
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html#installation) (v4.0+)
- [VSCode](https://code.visualstudio.com/) 配合 [GitHub Copilot](https://github.com/features/copilot) (推荐)

## 项目结构

```
.
├── .vscode/              # VSCode配置
├── first-operator/       # 示例operator项目
├── scripts/              # 实用脚本
├── kind-config.yaml      # Kind集群配置
├── README.md             # 英文文档
└── README_zh.md          # 中文文档
```

## 快速开始

### 1. 创建本地Kubernetes集群

使用Kind创建用于本地开发的集群：

```bash
./scripts/create-kind-cluster.sh
```

这将创建一个名为 `operator-dev` 的Kind集群，包含一个control-plane节点和一个worker节点。

### 2. 构建和运行 first-operator

进入first-operator目录：

```bash
cd first-operator
```

#### 方式A: 本地运行（集群外）

```bash
# 将CRD安装到集群
make install

# 本地运行operator（连接到当前kubeconfig上下文）
make run
```

#### 方式B: 部署到集群

```bash
# 构建并部署到Kind集群
cd ..
./scripts/deploy-operator.sh
```

### 3. 测试Operator

创建一个示例FirstApp自定义资源：

```bash
cd first-operator
kubectl apply -f config/samples/apps_v1alpha1_firstapp.yaml
```

检查创建的资源：

```bash
# 查看FirstApp资源
kubectl get firstapps

# 查看operator创建的deployment
kubectl get deployments

# 查看pods
kubectl get pods
```

## first-operator

`first-operator` 是一个示例Kubernetes operator，管理名为 `FirstApp` 的自定义资源。它根据FirstApp规范自动创建和管理Kubernetes Deployments。

### FirstApp自定义资源

FirstApp CRD允许您定义：

- **size**: 副本数量（1-10）
- **image**: 要部署的容器镜像（必需）
- **port**: 要暴露的端口（1-65535）

示例：

```yaml
apiVersion: apps.example.com/v1alpha1
kind: FirstApp
metadata:
  name: firstapp-sample
spec:
  size: 2
  image: nginx:latest
  port: 80
```

### 功能特性

- 自动为FirstApp资源创建Deployments
- 管理副本数量
- 当FirstApp规范变更时更新deployment
- 在状态中跟踪运行中的副本数

## 开发

### VSCode配置

本仓库包含VSCode配置和推荐扩展：

- Go扩展与语言服务器
- GitHub Copilot
- Kubernetes工具
- YAML支持

在VSCode中打开工作区，并在提示时安装推荐的扩展。

### 调试

使用VSCode调试器配置来调试operator：

1. 在代码中设置断点
2. 按F5或使用"Debug First Operator"启动配置
3. operator将以调试器附加的方式在本地运行

### 创建新的Operator

要在此工作区中创建新的operator：

```bash
# 创建新目录
mkdir my-new-operator
cd my-new-operator

# 使用kubebuilder初始化
kubebuilder init --domain example.com --repo github.com/yubo-yue/my-operator/my-new-operator

# 创建API
kubebuilder create api --group <group> --version <version> --kind <Kind>
```

## 清理

删除Kind集群：

```bash
./scripts/delete-kind-cluster.sh
```

从集群卸载operator：

```bash
cd first-operator
make undeploy
```

## 常用命令

### first-operator命令

```bash
cd first-operator

# 安装CRD
make install

# 卸载CRD
make uninstall

# 构建operator二进制文件
make build

# 运行测试
make test

# 生成清单（CRD、RBAC等）
make manifests

# 构建和推送Docker镜像
make docker-build docker-push IMG=<registry>/first-operator:tag

# 部署到集群
make deploy IMG=<registry>/first-operator:tag

# 从集群卸载
make undeploy

# 运行代码检查
make lint
```

### Kubernetes命令

```bash
# 查看所有FirstApp资源
kubectl get firstapps

# 描述一个FirstApp
kubectl describe firstapp firstapp-sample

# 查看operator日志（当部署到集群时）
kubectl logs -n first-operator-system deployment/first-operator-controller-manager -f

# 查看事件
kubectl get events --sort-by='.lastTimestamp'
```

## 快速演示

运行完整的演示（创建集群、部署operator、创建示例）：

```bash
make demo
```

## 资源

- [Kubebuilder文档](https://book.kubebuilder.io/)
- [Kubernetes Operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Kind文档](https://kind.sigs.k8s.io/)

## 详细文档

- [快速开始指南](QUICKSTART.md) - 5分钟入门
- [开发指南](DEVELOPMENT.md) - 详细的开发者文档

## 许可证

详见 [LICENSE](LICENSE) 文件。
