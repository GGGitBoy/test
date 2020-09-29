### 背景:
广发在secret additional-scrape-configs.yaml中添加一个自定义监控Job，监控升级时，自定义添加的Job会被抹除

### 原因:
* additional-scrape-configs.yaml的配置会被operator controller刷掉。监控配置的统一入口应在 **集群->工具->监控**

* 例子

如要添加一个job:
```
 - job_name: "prometheus"
   static_configs:
   - targets:
     - "localhost:9090"
```

目前在监控页面的应答需要添加的形式:

```
prometheus.additionalScrapeConfigs[0].job_name  =  prometheus
prometheus.additionalScrapeConfigs[0].static_configs[0].targets[0]  = localhost:9090
```

### 需求:

* 以上方式当自定义内容较多时，比较繁琐，不够友好。所以基于以上情况想要在启用监控页面，支持yaml编辑。

### 存在的问题:

1. 支持了valuesYaml的配置后，应用商店监控会出现 answer 和 yaml 展示内容不一致的问题（前端在判断valuesYaml有内容后，会将内容直接填充到 yaml 的部分，而不转换为 answer）
2. 用户可能会通过 应用商店 去修改helm部署的监控组件


注: 监控的helm部署会通过cluster handler去控制, 与其他系统组件相比比较特殊。
在用户不清楚我们监控的部署流程时，可能会去应用商店修改cluster-monitoring或monitoring-operator的内容。这样的修改, 在集群进行变动时就会被覆盖
  

### 觉得可执行的方案
1. 将cluster-monitoring和monitoring-operator的upgrade禁用
2. 以文档的形式说明监控需要在哪做相应的配置

**Relate issue**: https://github.com/cnrancher/pandaria/issues/1054