| 资源类型                              | 云厂商 | 操作系统               | 配置   | 带宽 | QPS            | 单条请求数据量 | 传输速率     |
|-----------------------------------|-----|--------------------|------|----|----------------|---------|----------|
| RUM-Session/View/Resource         | 阿里云 | Ubuntu 20.04.6 LTS | 2C4G |    | 519 request/s  | 15.7Kb  | 8Mb/s    |
| RUM-Session/View/Resource         | 阿里云 | Ubuntu 20.04.6 LTS | 4C8G |    | 1130 request/s | 15.7Kb  | 17Mb/s   |
| RUM-Error（no-sourcemap）           | 阿里云 | Ubuntu 20.04.6 LTS | 4C8G |    | 2862 request/s | 6.6Kb   | 18.6Mb/s |
| RUM-Error（no-sourcemap）           | 阿里云 | Ubuntu 20.04.6 LTS | 4C8G |    | 2862 request/s | 6.6Kb   | 18.6Mb/s |
| RUM-Error (sourcemap: android）    | 阿里云 | Ubuntu 20.04.6 LTS | 4C8G |    | 2.5 request/s  | 6.6Kb   | 904Kb/s  |
| RUM-Error (sourcemap: javascript） | 阿里云 | Ubuntu 20.04.6 LTS | 4C8G |    | 8276 request/s | 1.22Kb  | 9.86Mb/s |
| RUM-Error (sourcemap: ios）        | 阿里云 | Ubuntu 20.04.6 LTS | 4C8G |    | 52 request/s   | 3.1Kb   | 161Kb/s  |
