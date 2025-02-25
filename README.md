# rank
## 模块
单例启动一个rank模块，模块内维护一个有序player切片，一个玩家id与切片下标的map映射。
1. 玩家对榜单的更新操作，会放入一个cmd队列，模块每5秒（根据业务情况更改）执行队列中所有的命令。
2. player切片使用二分排序（再优化的话，可以考虑引用第三方库，用红黑树，go没有原生的树结构）
3. player信息中维护一个排名，而不是直接用切片下标作为排名，每次排序完成后，更新排名字段，以实现密集排名。

## 用户
简单的基于http的req resp，用户的请求，会发送到rank模块的cmd通道中执行。
![截屏2025-02-26 00.29.16.png](..%2F..%2F..%2Fvar%2Ffolders%2F1p%2Fkpklbfs107q0n83925zv9z500000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_vYvyWQ%2F%E6%88%AA%E5%B1%8F2025-02-26%2000.29.16.png)