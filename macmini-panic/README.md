# Mac mini 连续崩溃日志与分析

2026 年 11 月 9 日的这 5 次 kernel panic 指向同一种故障机制：`kernelmanagerd` 等待 `syspolicyd`，后者等待 `fseventsd`；与此同时，`fseventsd` 的工作线程在同一条 APFS 调用路径中进入不可中断等待，最终 `kernelmanagerd` 超过 60 秒没有报到，由 watchdog 触发 panic。

**这 5 次不只是报错标题相似，跨进程等待关系、线程状态和 APFS 栈偏移也一致，因此「同一种故障反复发生」有较强证据支持。** 但现有日志还不能把最终根因锁定为 macOS/APFS 缺陷、文件系统状态异常、外接 SSD、硬盘盒、线材或 USB 连接中的某一个。

与仓库中[此前的 3 份日志](../three-panic-logs-2026-10-24-25.txt)相比，2026 年 10 月 24 日的 `kernelmanagerd` panic 也出现了相同的服务等待链和相应的文件系统等待。它是与这 5 次最接近的历史案例；不能据此把此前 3 次全部认定为同一个最终根因。

## 日志清单与时间线

以下时间均为 2026 年 11 月 9 日，北京时间 UTC+08:00。文件名和报告时间晚于内核记录的 `Epoch Time / Calendar`，两者分别列出，避免把报告时间直接当作 panic 发生时间。时间差的具体原因尚未确定。

| 原始报告 | 报告时间 | 内核记录的 panic 时间 | 启动至 panic |
| --- | --- | --- | --- |
| [panic-full-2026-11-09-230114.0002.panic](panic-full-2026-11-09-230114.0002.panic) | 23:01:14 | 22:58:49 | 113 秒 |
| [panic-full-2026-11-09-230327.0002.panic](panic-full-2026-11-09-230327.0002.panic) | 23:03:27 | 23:02:15 | 93 秒 |
| [panic-full-2026-11-09-230514.0002.panic](panic-full-2026-11-09-230514.0002.panic) | 23:05:14 | 23:04:28 | 122 秒 |
| [panic-full-2026-11-09-231401.0002.panic](panic-full-2026-11-09-231401.0002.panic) | 23:14:01 | 23:06:15 | 96 秒 |
| [panic-full-2026-11-09-232258.0002.panic](panic-full-2026-11-09-232258.0002.panic) | 23:22:58 | 23:15:02 | 93 秒 |

启动至 panic 的时长按内核 `Boot` 和 `Calendar` 字段的整数秒计算。5 份报告的 `Sleep` 和 `Wake` 都为零，说明这些启动后的重复故障没有记录到中间睡眠、唤醒过程。

本目录保留 5 份未经修改的 `.panic` 原文件，共 2,802,818 字节。[SHA256SUMS](SHA256SUMS) 记录各文件的 SHA-256。诊断元数据、进程名、UUID 和地址保持原样。

## 共同的报错与等待链

5 份日志都包含以下 panic 原文：

```text
IOKit Daemon (kernelmanagerd) stall[0], (60s): registry root held busy, kernelmanagerd has not checked in @IOService.cpp:6140
```

5 次的 `Panicked task` 都是 PID 125 的 `watchdogd`。它检测到服务未及时报到并触发 panic，不能仅因它出现在这里就把它认定为根因。

通过 `processByPid` 中的 `turnstileInfo`、`waitInfo` 和线程栈，可以还原出以下关系：

1. `kernelmanagerd` 主线程等待 `com.apple.security.syspolicy.kext` 服务。
2. 该请求阻塞在 `syspolicyd` 主线程；后者又等待 FSEvents 服务。
3. `fseventsd` 的 `mach_channel_queue` 线程等待条件变量；同一进程的工作线程阻塞在 APFS 中。
4. `kernelmanagerd` 未在时限内完成报到，触发 60 秒超时和 kernel panic。

其中，跨进程的等待对象由日志直接记录；FSEvents 条件变量与 APFS 工作线程之间的依赖，是根据同时出现的线程状态推断的。日志没有完整标出该条件变量的拥有者，也没有解析到 APFS 内部具体的锁或 I/O 依赖。

| 报告后缀 | kernelmanagerd 线程 | syspolicyd 线程 | FSEvents 请求线程 | 阻塞在 APFS 的 FSEvents 工作线程 |
| --- | --- | --- | --- | --- |
| 230114 | 2435 | 3230 | 2610 | 2478 |
| 230327 | 2413 | 3252 | 2600 | 2463 |
| 230514 | 2449 | 3190 | 2604 | 2501 |
| 231401 | 2430 | 3420 | 2581 | 2478 |
| 232258 | 2440 | 3156 | 2524 | 2477 |

上述 5 个工作线程的状态均包含 `TH_WAIT, TH_UNINT`，APFS 栈帧偏移也完全一致，按日志记录顺序为：

```text
1398348 -> 1408024 -> 947424 -> 531676
```

这些帧均引用各自报告中的 `binaryImages[8]`，UUID 为 `4102CB18-9E03-3523-8F8D-011ACB52F131`。随后产生的本机 watchdog spin 报告将该 UUID 映射为 `com.apple.filesystems.apfs`，版本为 `3288.1.3`。当前日志未将这些 APFS 地址符号化为具体函数名。

**共同的等待链、不可中断等待状态和相同 APFS 偏移，是判断这 5 次属于同一种故障机制的主要依据。**

## 本机环境与补充证据

- 机型：Mac mini，型号标识 `Mac16,10`，内存 24 GiB。
- 系统：5 份日志均为 macOS 27.0.1，构建号 `26A434`，内核 Darwin 27.0.0。
- 启动盘：当前系统从 USB 外接的 `SSD 980 500GB` 启动。根目录 `/` 对应快照 `disk5s2s1`，其 APFS 容器为 `disk5`，物理存储为 `disk4s2`；数据卷 `disk5s1` 也在同一外接盘上。
- 可用空间：检查时外接 APFS 容器约剩余 154.4 GB。panic 原文中的内存压缩和交换空间指标也显示正常，因此这些报告不支持将磁盘容量耗尽或压缩内存耗尽视为直接触发因素。
- 硬盘健康信息：`diskutil info disk4` 返回 `SMART Status: Not Supported`。内置 Apple SSD 的 `Verified` 状态属于另一块盘，不能用于证明外接启动盘健康。
- 最后一次重启后，23:24—23:25 左右还产生了 `logd`、`configd` 和 `WindowServer` 的 watchdog 超时报告。`configd` 的 spin 报告也出现 APFS 阻塞栈，与系统多个服务受阻的现象一致，但不等于已经证明物理硬盘损坏。
- 本轮连续重启之前，还有一份 `panic-base+socd-2026-11-09-225745.000.panic`，记录了 `initproc exited -- exit reason namespace 2 subcode 0xa`，出事进程是 `launchd`，内核记录的 panic 时间为 22:56:45。该报告缺少完整 stackshot，无法建立同一等待链；它不属于本次要求上传的 5 份日志。

外接盘承担系统和数据卷，使 USB、硬盘盒、线材和 SSD 路径成为优先排查对象。但目前的栈不能定位到具体受影响的 APFS 卷，也无法区分 APFS 锁等待和更底层的 I/O 故障。

## 与 10 月日志的对比

| 旧报告时间 | 系统版本 | 对比结果 |
| --- | --- | --- |
| 10 月 24 日 17:01:47 | 26.5.2，`25F84` | `launchd` 退出的 namespace/subcode 与 11 月 9 日单独的第一份 panic 一致。仅凭退出原因相同，不能认定根因相同。 |
| 10 月 24 日 17:40:33 | 26.5.2，`25F84` | 同为 `kernelmanagerd` 60 秒超时，报错位置为 `IOService.cpp:5986`。线程 2474 等待 `syspolicyd` 线程 3370，后者等待 `fseventsd` 线程 2653；FSEvents 工作线程也处于 `TH_WAIT, TH_UNINT`，并出现相应文件系统栈。这是与本次 5 份日志最强的匹配。 |
| 10 月 25 日 19:50:56 | 27.0.1，`26A434` | 报错为 `watchdogd` 连续 91 秒没有报到，直接触发条件不同。尚未证明它具有同一条完整等待链；日志中的 `magpie` 任务名也不能单独作为其导致崩溃的证据。 |

不同系统构建的源码行号和栈偏移会变化，不能直接用绝对地址是否相同判断故障是否相同。这里主要比较服务依赖关系、等待状态和调用路径。

## 后续如何定位最终根因

以下为建议的对照验证，尚未执行：

1. **先备份重要数据，再进行文件系统修复或反复启动测试。**
2. 完全关机后，更换外接启动盘的线材、接口或硬盘盒，优先测试直接连接、不经过扩展坞的状态。每次只改变一个变量，记录是否复现。系统运行期间不要拔掉当前启动盘。
3. 如果内置盘上有已确认可正常使用的 macOS，可在关机后断开外接盘，从内置系统启动进行对照。若只有外接系统复现，可缩小到外接系统或存储路径，但仍不能直接区分软件与硬件。
4. 进入 macOS 恢复模式，用「磁盘工具」的「急救」依次检查外接盘的卷、容器和物理设备，并记录实际检查结果。当前栈本身不能证明文件系统已经损坏。
5. 测试安全模式及适用的 macOS 更新。如果仍复现，将这些日志和外接启动盘连接方式一并提供给 Apple 或硬盘盒、SSD 厂商。

官方参考：[Mac 意外重启的排查步骤](https://support.apple.com/en-my/102382)、[Apple 芯片 Mac 的恢复模式、安全模式与存储修复](https://support.apple.com/en-ie/guide/mac-help/mchl82829c17/27/mac/27)。访问日期：2026 年 11 月 9 日。

## 证据来源与验证结果

分析日期：2026 年 11 月 9 日。本机证据来自 `/Library/Logs/DiagnosticReports/`，主要为本目录中的 5 份日志，以及 `configd_2026-11-09-232514_shownJis-Mac-mini.userspace_watchdog_timeout.spin`。当前存储布局通过 `diskutil info /`、`diskutil list`、`diskutil apfs list` 和 `diskutil info disk4` 核对，系统版本与内存通过 `sw_vers` 和 `sysctl` 核对。旧日志因网页查看器无法完整展示，改用 GitHub API 读取全文后比较。

已解析每份 panic 文件的两段 JSON，核对共同报错、跨进程线程引用、APFS 镜像 UUID 和栈偏移，并验证原文件与上传副本的 SHA-256 一致。本次没有执行文件系统修复、压力测试、硬件替换或启动配置变更。此次新增的是诊断资料，不涉及应用构建、代码测试覆盖率或性能测试。

在本目录执行以下命令，可验证 5 份日志是否完整：

```sh
shasum -a 256 -c SHA256SUMS
```

仓库变更可通过回退新增本目录的提交撤销，不影响系统中的原始日志；无迁移操作。
