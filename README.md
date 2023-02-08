# GC-KDM

## 现状与背景

目前 `github.com/rancher/kontainer-driver-metadata`(下称 KDM) 的数据来源主要有 3 块：

- RKE 部分数据使用代码进行维护
  - **Service Options** 组件启动选项(Linux/Windows) -- 以大版本为中心进行维护，并支持根据小版本进行补丁，Windows 环境下的组件参数单独维护。
  - **RKE System Images** -- 以 rke k8s 版本为中心维护。
  - **Addons template** -- 以 addon 的类型+模版版本为中心进行维护，最后会使用部分系统配置（如 RBAC 是否启用）以及 RKE system images 的数据内容进行渲染；每个类型与版本的组合会对应一个适用的 rke k8s version 范围。
  - **Docker 版本支持情况** -- 以 rke k8s 的主版本号为中心维护，维护一个支持的 Docker 版本范围。
  - **Rancher Server 版本是配情况以及默认版本** -- 每个 rke k8s version 都有一个适配的 Rancher 版本范围；并配置了 Rancher 版本对应默认的 rke k8s 版本以及默认情况下的版本号
  - **CIS 检查配置** -- 维护 CIS 版本与其支持的 rke k8s 版本范围
- K3s 部分数据使用 channels.yaml 进行维护
  - channels.yaml 用于 K3s 对外的 upgrade channel 服务
  - 默认 channel 下的 K3s 的版本信息
  - Rancher 各个版本中，K3s 的默认版本信息
  - 各个 K3s 版本信息以及其对应支持的 Rancher server 版本，以及 K3s Server/Agent 运行参数
- RKE2 部分数据使用 channels-rke2.yaml 进行维护
  - channels.yaml 用于 RKE2 对外的 upgrade channel 服务
  - 默认 channel 下的 RKE2 的版本信息
  - Rancher 各个版本中，RKE2 的默认版本信息
  - 各个 RKE2 版本信息以及其对应支持的 Rancher server 版本，以及 RKE2 Server/Agent 运行参数，以及 RKE2 charts 运行 chart 名称以及参数

项目 KDM 会使用上面的数据，聚合生成一个 data.json 文件，后续被以下功能进行使用

- Rancher Server 构建时，下载一份指定 commit 的 data.json 文件到容器目录中，在离线部署时被使用
- Rancher 代码中，使用类型 `github.com/rancher/rke/types/kdm.Data`对 data.json 进行结构化解码，以及 `github.com/rancher/rancher/pkg/channelserver` 包中使用在线/本地的 data.json 进行 k3s/rke2 的版本以及相关数据存储
- RKE 会从代码层面使用 `github.com/rancher/kontainer-driver-metadata` 包中的数据

### 上游 KDM 维护方式

KDM 声明周期通常使用以下两个渠道进行管理：

- 在 rancher/rancher 中使用独立的 Milestone 管理发布，如：[https://github.com/rancher/rancher/milestone/304](https://github.com/rancher/rancher/milestone/304)
- 合并在 Rancher release milestone 中的 KDM issue，如：[https://github.com/rancher/rancher/issues/40210](https://github.com/rancher/rancher/issues/40210)

代码主要使用以下分支进行维护：

- dev-\[major version\] 分支维护各个 Rancher 主版本对应的开发中 KDM 数据
- release-\[major version\] 分支维护各个 Rancher 主版本对应的发布的 KDM 数据

1. 在开发阶段，主要在 dev 分支进行提交，将添加完成 RKE 代码后，使用 `go generate` 生成新的 data.json 文件，一并提交到 commit 中
1. 发布时从对应的 release 分支 checkout 出新分支，并 cherry-pick 需要的 commit 到新分支后，再提交 PR 到 release 分支进行合并
1. release 分支合并后，触发 CI，将生成的 data.json 文件上传到 cloudfront 进行分发。
1. 在 Rancher 构建时，将使用对应大版本的 release 分支 commit 拉取 data.json

## GC-KDM 的目标

GC-KDM 需要解决以下几个问题：

- 需要在 data.json 中增加 rke k8s 自定义版本，需要同时兼容 Rancher Prime 以及 Rancher 开源版的使用
- 需要支持 RFO（Rancher for openEuler）发行版，维护方式与 rke2 一致
- 需要支持生成 image-tool 的镜像 mirror list，以便在后续维护中不需要再对开源版的 Rancher 镜像进行调整

## GC-KDM 设计

### data.json 文件修改

针对 data.json 文件，我们需要增加 `rfo` 段落保存 rfo 发行版的信息，其余结构均无需改变。

### Rancher Prime RKE Kubernetes Distro

我们将引入后缀为 `-rancher-entX-X` RKE K8s 版本号作为 GC-KDM 中 Rancher Prime 使用的 k8s 版本号，该版本号符合 semver versioning 定义，并能很好的与开源版进行兼容。

以下面版本号为例子，按照 semver 排序后如下表：

1. 1.24.9-rancher1-1
1. 1.24.9-rancher1-2
1. 1.24.9-rancher2-2
1. 1.24.9-rancherprime1-1
1. 1.24.9-rancherprime1-2
1. 1.24.9-rancherprime2-1
1. 1.24.10-rancher1-1

在相同的 k8s 小版本的情况下，semver compare 场景 `-rancher-entX-X` > `-rancherX-X`，这样能直接适配开源 data.json 中 `Addons template` 以及 `Service Options` 选项，减少维护工作量。

### rancher-prime-patch.yaml 文件设计

对于我们新增维护的 RKE K8s 版本，我们引入一个新的配置文件进行维护：`rancher-prime-patch.yaml`。该文件 schema 参考 `https://github.com/rancher/channelserver/blob/master/pkg/model/config.go` 以及 k3s/rke2 的 channels.yaml 进行设计和扩展。文件结构说明如下：

- appDefaults(array): 如 channels.yaml 中，appDefaults 定义，用于保存不同 Rancher 版本下的默认版本，会合并到 data.json 中 `RancherDefaultK8sVersion` 段落
  - appName: 默认为 rancher
  - defaults(array): 数组中记录 rancher 版本号与 RKE K8s 默认版本关系
- mirrorOptions(array): 新增属性，用于生成 image mirror list 文件，在 `releases` 段落中维护的 system image 将会根据这个 mirrorOptions 进行调整
  - targetRegistry: 转换后的镜像仓库
  - regex(array): 转换使用的正则表达式列表，包含匹配和替换两个属性，如果需要调整 repo 名称，则可以在这里进行配置
    - regexMatch
    - regexReplacement
- releases(array): 如 channels.yaml 中，releases 定义，并对列表中的结构定义进行了扩展。
  - version: RKE K8s version 名称。
  - minChannelServerVersion: 当前版本 K8s 最低支持的 Rancher 版本。
  - maxChannelServerVersion: 当前版本 K8s 最高支持的 Rancher 版本。
  - serviceOptions(新增参数): 这里与 data.json 中，每个版本的 serviceOption 一致，优先查找其中大版本参数，如果指定了参数，则使用覆盖大版本参数中的对应参数或添加该参数
  - windowsServiceOptions(新增参数): 同 `serviceOptions` 定义，与 data.json 中 windowsServiceOptions 一致。
  - systemImages(新增参数): 复用 github.com/rancher/rke/types.RKESystemImages 结构，保存该版本使用的 System image，最终经过 mirror options 后注入到 data.json 中.
  - addonTemplates(map/新增参数): 这个 map 以 addon 种类为 key，value 为 addon template 内容。生成到 kdm data.json 中时，会根据重复的内容进行聚合并创建 addon 代号，并对支持的 k8s 版本进行范围聚合后输出到 K8sVersionedTemplates 中。

其中针对 `addonTemplates` 中，map 的 key 枚举如下，可以讨论是否支持对 addon 进行扩展，如新增 macvlan addon：

- aci
- calico
- canal
- coreDNS
- flannel
- kubeDNS
- metricsServer
- nginxIngress
- nodelocal
- weave

以下为维护文件例子：

```yaml
appDefaults:
  - appName: rancher
    defaults:
      - appVersion: ">= 2.5.0-0 < 2.6.0-0"
        defaultVersion: "1.20.x"
mirrorOptions:
  - targetRegistry: "" # 这里使用空表示mirror到dockerhub
    regex:
      - regexMatch: "^rancher/(.+)" # 将rancher下的镜像
        regexReplacement: "cnrahcner/$1" # 改为 cnrancher下的镜像
releases:
  - version: v1.20.15-rancherprime1-1
    minChannelServerVersion: "2.5.11-rc0"
    maxChannelServerVersion: "2.5.99"
    serviceOptions: null # 没有options需要更新
    windowsServiceOptions: null
    systemImages:
      etcd: rancher/mirrored-coreos-etcd:v3.4.15-rancher1
      alpine: cnrancher/cn-rke-tools:v0.1.80-ent2
      nginxProxy: cnrancher/cn-rke-tools:v0.1.80-ent2
      certDownloader: cnrancher/cn-rke-tools:v0.1.80-ent2
      kubernetesServicesSidecar: cnrancher/cn-rke-tools:v0.1.80-ent2
      kubedns: rancher/mirrored-k8s-dns-kube-dns:1.15.10
      dnsmasq: rancher/mirrored-k8s-dns-dnsmasq-nanny:1.15.10
      kubednsSidecar: rancher/mirrored-k8s-dns-sidecar:1.15.10
      kubednsAutoscaler: rancher/mirrored-cluster-proportional-autoscaler:1.8.1
      coredns: rancher/mirrored-coredns-coredns:1.8.0
      corednsAutoscaler: rancher/mirrored-cluster-proportional-autoscaler:1.8.1
      nodelocal: rancher/mirrored-k8s-dns-node-cache:1.15.13
      kubernetes: rancher/hyperkube:v1.20.15-rancher2
      flannel: rancher/mirrored-coreos-flannel:v0.15.1
      flannelCni: rancher/flannel-cni:v0.3.0-rancher6
      calicoNode: cnrancher/mirrored-calico-node:v3.17.2
      calicoCni: cnrancher/mirrored-calico-cni:v3.17.2
      calicoControllers: cnrancher/mirrored-calico-kube-controllers:v3.17.2
      calicoCtl: cnrancher/mirrored-calico-ctl:v3.17.2
      calicoFlexVol: cnrancher/mirrored-calico-pod2daemon-flexvol:v3.17.2
      canalNode: cnrancher/mirrored-calico-node:v3.17.2
      canalCni: cnrancher/mirrored-calico-cni:v3.17.2
      canalControllers: cnrancher/mirrored-calico-kube-controllers:v3.17.2
      canalFlannel: rancher/mirrored-coreos-flannel:v0.15.1
      canalFlexVol: cnrancher/mirrored-calico-pod2daemon-flexvol:v3.17.2
      weaveNode: weaveworks/weave-kube:2.8.1
      weaveCni: weaveworks/weave-npc:2.8.1
      podInfraContainer: rancher/mirrored-pause:3.6
      ingress: rancher/nginx-ingress-controller:nginx-1.2.1-rancher1
      ingressBackend: rancher/mirrored-nginx-ingress-controller-defaultbackend:1.5-rancher1
      ingressWebhook: rancher/mirrored-ingress-nginx-kube-webhook-certgen:v1.1.1
      metricsServer: rancher/mirrored-metrics-server:v0.5.0
      windowsPodInfraContainer: rancher/mirrored-pause:3.6
      aciCniDeployContainer: noiro/cnideploy:5.1.1.0.1ae238a
      aciHostContainer: noiro/aci-containers-host:5.1.1.0.1ae238a
      aciOpflexContainer: noiro/opflex:5.1.1.0.1ae238a
      aciMcastContainer: noiro/opflex:5.1.1.0.1ae238a
      aciOvsContainer: noiro/openvswitch:5.1.1.0.1ae238a
      aciControllerContainer: noiro/aci-containers-controller:5.1.1.0.1ae238a
      aciGbpServerContainer: noiro/gbp-server:5.1.1.0.1ae238a
      aciOpflexServerContainer: noiro/opflex-server:5.1.1.0.1ae238a
    addonTemplates:
      macvlan: |-
        {{ .DemoTemplate }}
        test: test
        abcd: abcd
      nginxIngress: |-
        test: ingress template
```

### 维护方式

GC-KDM 需要在维护流程上支持以下内容：

- GC-KDM 项目个分支中维护使用上游 KDM 的 commit id，分支设计上与 KDM 保持一致
  - dev-\[major version\]分支维护开发中的 KDM 数据；
  - release-\[major version\]分支维护已发布的 KDM 数据；
  - GC-KDM 的发布时对 release 分支进行 tag 操作；
  - 在 release 分支触发 tag 事件时，触发 CI，将 data.json 数据上传到 `https://releases.rancher.cn` 对应的 OSS Bucket 中，使用子路径 `kontainer-driver-metadata/<branch>/` 对外提供服务；
- 在 GC-KDM dev 分支合并后，更新结束后，需要输出新旧 GC-KDM data.json 差异
- 在 GC-KDM 发版时构建最新的 KDM nginx 镜像，使用 GC-KDM tag 名称作为镜像 tag，tag 的规则与场景如下：
  - 与 Rancher Prime 发版保持版本同步的 tag，如 `cnrancher/kdm:v2.7.2`
  - GC-KDM 单独发版的 tag，如 `cnrancher/kdm:v2.7-build20230210`，使用 `<branch name> + -build<date>[01]` 的格式进行命名
- 根据 GC-KDM 版本生成的 image mirror list，进行镜像 mirror 操作（nice to have）
