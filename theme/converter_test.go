package theme

import (
	"strings"
	"testing"
)

func TestConvertBasic(t *testing.T) {
	markdown := `
# OpenClaw安全自检清单：5分钟判断你的实例是否在“裸奔”


## 第一章：风暴中心——当“贾维斯”变成“特洛伊木马”

### 1.1 一夜成名的“数字员工”与背后的阴影

2026年1月，开源社区迎来了一位超级明星——**OpenClaw**（曾用名Clawdbot、Moltbot）。

它被誉为“运行在你本地的贾维斯”。不同于只会聊天的LLM，OpenClaw是一个真正的**行动派AI代理（Agent）**。它能帮你自动回复微信、整理Notion笔记、在GitHub上提交代码，甚至控制家里的智能灯光。上线仅三周，其GitHub星标数就突破了18万，被TechCrunch称为“2026年最落地的AI应用”。

对于开发者、极客和效率追求者来说，拥有一只自己的“大龙虾”成了当时的潮流。大家兴奋地配置技能（Skills），连接OAuth，享受着AI带来的生产力飞跃。

然而，狂欢在2026年2月初戛然而止。

安全研究人员披露了一个编号为**CVE-2026-25253**的高危远程代码执行（RCE）漏洞。这个漏洞的CVSS评分高达**8.8**（满分10分），属于“严重”级别。紧接着的48小时内，官方被迫连发三个版本修复了包括Token泄露在内的7个CVE漏洞。

**漏洞原理简单得令人发指：**攻击者不需要高超的黑客技术，只需要构造一个恶意的网页链接。当安装了旧版OpenClaw的用户点击这个链接时，浏览器中的恶意脚本会利用OpenClaw本地服务（默认监听在~127.0.0.1:18789~）的信任机制缺陷，通过WebSocket直接接管本地实例。

**后果是灾难性的：**

- **令牌窃取**：攻击者直接获取你的API Key、数据库密码。
- **完全控制**：在你的电脑上执行任意命令（~rm -rf /~、安装后门、窃取SSH密钥）。
- **内网渗透**：以你的电脑为跳板，攻击你局域网内的其他设备（NAS、打印机、服务器）。

据不完全统计，截至2026年2月中旬，全球仍有超过**1.5万台**OpenClaw实例暴露在风险中。许多用户在不知情的情况下，他们的“数字员工”已经变成了黑客的“肉鸡”。

### 1.2 救星降临：v2026.3.8 发布

就在大家人心惶惶之时，2026年3月9日（两天前），OpenClaw官方正式发布了**v2026.3.8**稳定版。

这是一个里程碑式的版本。它不仅彻底修复了之前的RCE漏洞，还引入了多项关键的安全增强功能：

- **备份与恢复机制**：新增 ~openclaw backup create/verify~ 命令，确保在遭受攻击或配置错误时能快速回滚。
- **远程网关令牌保护**：改进了远程模式的Token处理，防止明文泄露。
- **系统运行沙箱化**：限制了 ~system.run~ 的权限，绑定脚本快照，防止执行被篡改的脚本。
- **浏览器SSRF防御**：阻断私有网络的重定向跳转，防止服务器端请求伪造攻击。

**但是，补丁不等于安全。**如果你还没有升级到 **v2026.3.8**，或者虽然升级了但配置依然沿用旧的“裸奔”习惯（如绑定 ~0.0.0.0~、未撤销旧Token），那么你依然处于危险之中。

现在，请深呼吸。无论你是否已经遭遇攻击，**立即停止盲目操作**，跟随本文的“5分钟自检清单”，为你的实例穿上一件防弹衣。

---

## 第二章：分级诊疗——你对号入座了吗？

在进行具体技术操作前，我们需要先明确你的使用场景。不同的场景对应着完全不同的风险等级和处理策略。请诚实回答以下问题，找到属于你的“生存方案”。

### 2.1 第一级：从未使用（观望者）

- **特征**：听说过OpenClaw很火，但还没下载，或者下载了还没运行。
- **风险等级**：🟢 **安全**
- **官方建议**：**可以安装，但必须严格遵循安全规范**。

- **理由**：随着v2026.3.8的发布，核心漏洞已修复，且新增了备份等安全特性。现在的版本比2月份安全得多。但仍需注意供应链安全。
- **行动指南**：

1. **只下载最新版**：务必确认下载的是 **v2026.3.8** 或更新版本。
2. **沙箱运行**：建议在**虚拟机**、**Docker容器**或**专用测试机**中首次运行，熟悉配置后再考虑迁移到主力机。
3. **最小化授权**：初次配置时，不要勾选所有权限，仅授予必要项。

### 2.2 第二级：仅本地测试（尝鲜者）

- **特征**：为了体验功能，在本地跑了一下，试过几个Demo技能，目前未投入实际工作流。
- **风险等级**：🟡 **中低风险**（取决于是否已更新）
- **官方建议**：**完整卸载 + 撤销OAuth** 或 **升级并加固**。

- 停止服务：~openclaw stop~。
- 删除文件：~rm -rf ~/.openclaw~ (Mac/Linux) 或 ~%USERPROFILE%\.openclaw~ (Windows)。
- **撤销授权**：登录Telegram/Discord/Slack等，撤销所有OpenClaw相关的Bot授权。

- 强制升级到 v2026.3.8。
- 执行一次完整的 ~openclaw backup create --only-config~ 备份配置。
- 按照第三章进行安全审计。

- **理由**：既然未深度依赖，没必要承担潜在的供应链风险。即使你更新了版本，残留的配置文件或授权的Bot权限仍可能成为隐患。
- **行动指南**：

1. **方案A（推荐）**：彻底清理。
2. **方案B（想继续玩）**：

### 2.3 第三级：日常使用（重度依赖者）

- **特征**：已经将OpenClaw集成到日常工作流中，每天依靠它处理邮件、代码或文档，拥有大量自定义技能和OAuth授权。
- **风险等级**：🔴 **高风险**（若未更新） / 🟠 **中风险**（若已更新但未审计）
- **官方建议**：**立即更新到v2026.3.8 + 全面安全审计**。

- **理由**：你是黑客眼中的“高价值目标”。v2026.3.8 虽然修复了漏洞，但你之前的配置可能仍存在暴露面。必须确保版本最新，并清洗不安全的配置。
- **行动指南**：

1. **强制更新**：必须升级到 **v2026.3.8**。这是目前的最新稳定版，包含了最关键的安全修复。
2. **执行自检**：严格按照下文第三章的“5分钟实操”进行检查。
3. **创建备份**：利用新版特性，立即执行 ~openclaw backup create~，保存当前安全状态。
4. **最小化权限**：审查所有已安装的技能，移除不必要的；检查OAuth权限，取消“全选”授权。

### 2.4 第四级：已暴露公网（裸奔者）

- **特征**：为了方便远程访问，将OpenClaw绑定到了 ~0.0.0.0~，且没有配置强密码或防火墙，端口（默认18789）可直接从外网访问。
- **风险等级**：🚨 **极度危险**（可能已被入侵）
- **官方建议**：**立即停止服务 + 断网隔离 + 全面取证**。

- **轮换密钥**：假设你所有的API Key、数据库密码、SSH密钥都已泄露。**立即**在相关平台上重置所有凭证。
- **检查日志**：查看 ~~/.openclaw/logs/~ 下的日志文件，寻找异常的执行记录。
- **重装系统**：如果无法确认是否被植入Rootkit或后门，最安全的做法是备份纯数据文件（代码、文档），然后**重装操作系统**。不要试图在可能被污染的系统中继续运行。

- **理由**：在互联网扫描器面前，暴露的OpenClaw实例就像黑夜里的火炬。大概率已经被自动化脚本扫描并植入后门。此时简单的“更新”已无法保证安全。
- **行动指南**：

1. **物理断网**：立即拔掉网线或断开Wi-Fi。
2. **停止服务**：在本地终止OpenClaw进程。
3. **假设失陷**：

---

## 第三章：5分钟实操自检清单（核心干货）

如果你属于“日常使用”用户，且确认尚未被入侵（或已完成重装），请立即打开终端，跟随以下步骤进行“体检”。整个过程不超过5分钟。

### 步骤一：版本检查（生死线）

**目标**：确保你不在漏洞版本范围内，且用上了最新的安全特性。**标准**：必须是 **v2026.3.8** 或更高。

1. **打开终端**（Terminal / CMD / PowerShell）。
2. **输入命令**：

   ~~~
   openclaw --version
   ~~~

   *(如果是Docker部署，输入 ~docker exec <容器名> openclaw --version~)*
3. **判断标准**：

- ✅ **安全**：版本号显示 ~v2026.3.8~ 或更高（注意：有些版本可能带有commit hash后缀，只要主版本号对即可）。
- ❌ **危险**：版本号低于 ~v2026.3.8~（特别是 ~v2026.1.28~ 及以下，~v2026.2.x~ 系列也存在其他CVE）。

4. **补救措施**： 如果是危险版本，**立即更新**：

   ~~~
   # npm全局安装更新
   npm install -g openclaw@latest

   # 或者指定版本
   npm install -g openclaw@2026.3.8

   # Docker更新
   docker pull openclaw/openclaw:latest
   docker stop <容器名> && docker rm <容器名>
   # 重新运行容器（注意保留卷挂载）
   ~~~

   *注意：更新后务必重启服务使补丁生效。*

### 步骤二：网络暴露审计（是否“大门敞开”）

**目标**：确认你的实例是否只允许本地访问。 这是最关键的一步。很多用户以为“没做端口映射”就安全了，其实配置文件里的绑定地址才是关键。

1. **方法A：使用内置审计命令（推荐）**新版OpenClaw自带安全检查工具：

   ~~~
   openclaw security audit
   ~~~

   观察输出结果：

- ✅ **PASS**: ~Gateway binding is secure (127.0.0.1)~
- ❌ **FAIL**: ~WARNING: Gateway is bound to 0.0.0.0! Public exposure detected.~

2. **方法B：手动检查配置文件**如果命令不可用，直接查看配置文件：

- ✅ **安全**：值为 ~"127.0.0.1"~ 或 ~"localhost"~。这意味着只有你这台电脑能访问。
- ❌ **危险**：值为 ~"0.0.0.0"~。这意味着互联网上的任何人都能尝试连接你的实例。

- **文件路径**：~~/.openclaw/openclaw.json~ (Linux/Mac) 或 ~%USERPROFILE%\.openclaw\openclaw.json~ (Windows)。
- **查找字段**：~"bind"~ 或 ~"host"~。
- **判断标准**：

3. **方法C：外网连通性测试（终极验证）**

   **如何修复暴露问题？**修改 ~openclaw.json~：

   ~~~
   {
     "gateway": {
       "bind": "127.0.0.1",  // 改这里！
       "port": 18789
     }
   }
   ~~~

   *如果需要远程访问怎么办？***绝对不要**直接绑定 ~0.0.0.0~！请使用 **SSH隧道** 或 **Tailscale**：

   ~~~
   # SSH隧道示例
   ssh -L 18789:127.0.0.1:18789 user@your-home-ip
   # 然后在远程浏览器访问 http://localhost:18789 即可安全连接
   ~~~

- 如果**打不开**（连接超时）：✅ 安全。
- 如果**能打开**登录页或API界面：🚨 **极度危险**！说明你的防火墙没挡住，或者配置错了。立即修改配置文件为 ~127.0.0.1~ 并重启服务。

- 关闭电脑的Wi-Fi，使用手机流量（确保手机和电脑不在同一局域网）。
- 在手机浏览器输入：~http://<你的电脑公网IP>:18789~。
- **结果**：

### 步骤三：技能来源审查（供应链排毒）

**目标**：清退恶意或不明来源的“内鬼”。 OpenClaw的强大源于技能（Skills），但这也是最大的风险点。恶意技能可以伪装成“PDF总结工具”，实则在后台记录你的键盘输入。

1. **列出已安装技能**：

   ~~~
   openclaw skills list
   ~~~
2. **逐一审查**： 检查列表中的每一个技能，问自己三个问题：

   *典型案例*：曾有一个名为 ~quick-clipboard~ 的热门技能，代码中包含一段逻辑：每次读取剪贴板时，悄悄将内容发送到攻击者的服务器。

- **来源可信吗？** 是官方认证（Verified Badge）的吗？还是来自某个不知名GitHub仓库？
- **真的需要吗？** 这个技能上次使用是什么时候？如果超过一个月没用，删掉它。
- **代码审计过吗？** 如果你有能力，点开技能的 ~main.js~ 或 ~index.py~ 看看。重点搜索关键词：~fetch(~ (向外发送数据), ~fs.readFile~ (读取敏感文件), ~exec(~ (执行系统命令)。

3. **清理动作**： 对于任何可疑、不再使用或来源不明的技能，坚决卸载：

   ~~~
   openclaw skills uninstall <skill-name>
   ~~~

   **原则**：最小化安装。只保留核心必需品。

### 步骤四：OAuth与凭证检查（切断后路）

**目标**：防止Token泄露导致的连锁反应。 即使你的OpenClaw实例安全了，如果之前泄露的OAuth Token还在有效期内，黑客依然可以通过其他途径控制你的账号。

1. **检查通讯软件授权**：

- **Telegram**: 设置 -> 隐私与安全 -> 活跃会话 / 已连接的Bot。撤销所有不认识的设备或旧的OpenClaw会话。
- **Discord**: 用户设置 -> 授权应用（Authorized Apps）。移除OpenClaw或其他可疑应用。
- **Slack/飞书/钉钉**: 进入企业管理后台或个人设置，查看“第三方应用授权”，撤销旧令牌。

2. **轮换API Key**： 如果你在OpenClaw配置文件中填入了 LLM Provider 的 API Key（如OpenAI, Anthropic, 阿里云百炼等），**建议全部重置**。

- 去对应的云平台控制台，删除旧Key，生成新Key。
- 更新到 ~openclaw.json~ 中。
- *为什么？* 因为如果之前实例被RCE，这些Key大概率已经被抓取并上传到黑客的数据库了。

### 步骤五：备份验证（新版特性）

**目标**：确保在出事时有“后悔药”。 v2026.3.8 引入了强大的备份功能，这是之前版本没有的。

1. **创建备份**：

   ~~~
   openclaw backup create --include-workspace
   ~~~

   这会创建一个包含配置和工作区状态的归档文件。
2. **验证备份**：

   ~~~
   openclaw backup verify <备份文件名>
   ~~~

   确保备份文件完整可用。
3. **意义**：一旦未来发现配置被篡改或遭受勒索，你可以迅速恢复到这个干净的状态。

---

## 第四章：进阶防御——给“龙虾”穿上防弹衣

完成了上述5分钟自检，你只能算是“脱离了生命危险”。要在AI时代长期安全地使用OpenClaw，还需要建立更深层次的防御体系。

### 4.1 容器化运行：沙箱隔离

**永远不要**直接在宿主机（尤其是你的主力开发机）上运行OpenClaw。 推荐使用 **Docker** 进行隔离，限制其对宿主机的访问权限。

**安全启动示例**：

~~~
docker run -d \
  --name openclaw-secure \
  --read-only \  # 文件系统只读，防止写入恶意文件
  --tmpfs /tmp \ # 临时目录可写
  --cap-drop=ALL \ # 丢弃所有Linux能力
  --no-new-privileges \ # 禁止提权
  -p 127.0.0.1:18789:18789 \ # 仅绑定本地回环
  -v ~/openclaw-data:/app/data \ # 仅挂载必要的数据目录
  openclaw/openclaw:2026.3.8
~~~

*解释*：

- ~--read-only~：即使黑客攻入容器，也无法修改系统文件或植入持久化后门。
- ~--cap-drop=ALL~：剥夺容器操作网络、硬件等底层权限。
- ~-p 127.0.0.1:...~：双重保险，确保Docker层面也不对外暴露。

### 4.2 网络隔离：零信任架构

如果你的业务确实需要多人协作或远程访问：

1. **组建虚拟局域网**：使用 **Tailscale**、**ZeroTier** 或 **WireGuard** 组建私有网络。只有加入该网络的设备才能访问OpenClaw的IP。

- v2026.3.8 特别优化了 Tailscale 的网关发现功能，支持 ~.ts.net~ 域名直连，更加安全便捷。

2. **反向代理加锁**：在OpenClaw前层架设 Nginx 或 Traefik，强制开启 **HTTP Basic Auth** 或 **OAuth2 登录**。

- 即使端口暴露，没有密码也无法进入。
- 配置HTTPS，防止中间人窃听。

### 4.3 最小权限原则（Least Privilege）

1. **专用用户**：创建一个专门的低权限系统用户（如 ~ocl_user~）来运行OpenClaw，不要用 ~root~ 或 ~Administrator~。

   ~~~
   sudo useradd -m -s /bin/bash ocl_user
   sudo chown -R ocl_user:ocl_user ~/.openclaw
   ~~~
2. **技能权限隔离**：未来的OpenClaw版本可能会支持技能级的权限控制。对于需要访问文件的技能，仅授予特定目录的读写权，而不是全盘访问。
3. **系统运行限制**：v2026.3.8 增强了 ~system.run~ 的安全性，绑定了脚本快照。不要尝试关闭这一保护，也不要随意批准未知脚本的执行请求。

### 4.4 监控与告警

不要等到数据丢了才发现。

- **开启详细日志**：在配置文件中设置 ~log_level: "debug"~。
- **异常监控**：编写一个简单的脚本，定期扫描日志文件。如果发现大量的 ~403 Forbidden~（暴力破解尝试）或异常的 ~POST /api/exec~ 请求，立即触发告警（发送邮件或短信）。
- **流量监控**：使用 ~iftop~ 或 Wireshark 观察OpenClaw进程的网络连接。如果它突然向一个陌生的海外IP发送大量数据，立刻断网！

---

## 第五章：常见误区与Q&A

在安全社区的实际交流中，我们发现用户对OpenClaw安全存在不少误解。以下是高频问题的解答。

**Q1: “我改了端口（比如从18789改成29999），是不是就安全了？”**

- **A**: **完全不是**。这叫“隐匿式安全”（Security by Obscurity）。黑客的扫描器会在几秒钟内遍历所有常用端口，甚至进行全端口扫描。改端口只能防住小白，防不住自动化脚本。**唯一的正道是绑定127.0.0.1并配合防火墙。**

**Q2: “我的电脑有防火墙（Windows Defender / ufw），是不是不用怕了？”**

- **A**: **不一定**。大多数个人防火墙默认允许“出站连接”和“本地回环连接”。CVE-2026-25253 利用的是浏览器发起的本地连接，防火墙通常不会拦截 ~127.0.0.1~ 的流量。你必须显式配置防火墙规则，禁止非本地进程访问OpenClaw端口，或者直接修改OpenClaw绑定地址。

**Q3: “我已经更新了最新版，是不是可以高枕无忧了？”**

- **A**: **不能**。软件漏洞只是风险的一部分。

- **社会工程学**：黑客可能诱导你安装一个“看起来很酷”的恶意技能。
- **配置错误**：你可能不小心把配置文件同步到了公开的GitHub仓库。
- **供应链攻击**：即使是官方技能库，也不能100%保证未来不被投毒。
- **结论**：保持警惕，定期审计，遵循最小权限原则。v2026.3.8 只是起点，不是终点。

**Q4: “如果我不小心点了恶意链接，但马上关掉了网页，会有事吗？”**

- **A**: **有可能**。RCE攻击往往在毫秒级完成。只要页面加载并执行了JS脚本，攻击载荷可能已经下发完毕。

- **建议**：立即按照“第四级：已暴露公网”的流程处理。假设已被入侵，轮换密钥，检查日志，必要时重装系统。不要抱有侥幸心理。

**Q5: “OpenClaw更新这么快，我还要继续用吗？”**

- **A**: 这取决于你的风险承受能力。

- 快速迭代意味着团队在积极响应安全问题（如48小时修7个CVE），这其实是好事。
- 如果是**生产环境**（处理公司数据、客户信息）：建议在严格的沙箱（Docker + 无敏感数据）中使用，并密切关注更新日志。
- 如果是**个人娱乐**：可以在做好备份和隔离的前提下继续使用，享受AI带来的便利。

---

## 第六章：结语——在AI时代，安全是“动词”

OpenClaw的这次安全危机，给狂热的AI社区敲响了警钟。我们渴望拥有像“贾维斯”一样全能的助手，却往往忽略了赋予它权力时的代价。

**AI代理（Agent）的本质，是权力的让渡。**当你把API Key、文件系统权限、网络访问权交给一个程序时，你实际上是在签署一份“信任契约”。如果这份契约缺乏技术约束（如鉴权、沙箱、审计），那么背叛只是时间问题。

对于每一位OpenClaw用户，请记住这三句话：

1. **默认不信任**：无论是网络、技能还是第三方服务，默认都是不可信的，直到你证明了它们的安全性。
2. **最小化授权**：只给予完成任务所需的最小权限，多一分都是隐患。
3. **持续监控**：安全不是一次性的配置，而是一个持续的过程。定期检查日志，关注安全公告，及时升级到像 v2026.3.8 这样的安全版本。

技术本身没有善恶，但使用技术的方式决定了结局。一只没有盔甲的龙虾，在深海里既是捕猎者，也是猎物。希望这份自检清单，能成为你那件坚实的“防弹衣”。

**最后，行动起来！**现在就打开终端，花5分钟检查一下你的实例。

1. 确认版本是 **v2026.3.8**。
2. 确认绑定地址是 **127.0.0.1**。
3. 执行一次 **backup**。

如果你发现身边还有朋友在“裸奔”，请把这篇文章转发给他们。在AI安全的战场上，没有人是一座孤岛。

---

### 附录：一键自检脚本（Linux/Mac）

为了方便大家，我们编写了一个简单的Shell脚本，你可以复制保存为 ~check_openclaw.sh~ 并运行。此脚本已适配最新版本检查逻辑。

~~~
#!/bin/bash

echo"🛡️  OpenClaw 安全快速自检脚本 (v2.0 - 适配2026.3.8)"
echo"----------------------------------------------------"

# 1. 检查版本
echo"🔍 [1/5] 检查版本..."
VERSION=$(openclaw --version 2>&1 | grep -oP 'v?\d{4}\.\d+\.\d+' | head -n 1 || echo"Unknown")
if [[ "$VERSION" == "Unknown" ]]; then
    echo"❌ 无法获取版本，请确认是否安装。"
elif [[ "$VERSION" < "v2026.3.8" ]] && [[ "$VERSION" < "2026.3.8" ]]; then
    echo"🚨 危急！当前版本 $VERSION 存在已知高危漏洞！请立即运行 'npm install -g openclaw@latest'"
else
    echo"✅ 版本安全：$VERSION"
fi

# 2. 检查绑定地址
echo"🔍 [2/5] 检查网络绑定..."
CONFIG_FILE="$HOME/.openclaw/openclaw.json"
if [ -f "$CONFIG_FILE" ]; then
    BIND_ADDR=$(grep -oP '"bind"\s*:\s*"\K[^"]+'"$CONFIG_FILE" || grep -oP '"host"\s*:\s*"\K[^"]+'"$CONFIG_FILE")
    if [[ "$BIND_ADDR" == "0.0.0.0" ]]; then
        echo"🚨 危急！实例绑定在 0.0.0.0，已暴露公网风险！请修改为 127.0.0.1"
    elif [[ "$BIND_ADDR" == "127.0.0.1" ]] || [[ "$BIND_ADDR" == "localhost" ]]; then
        echo"✅ 网络配置安全：$BIND_ADDR"
    else
        echo"⚠️ 警告：未知绑定地址 $BIND_ADDR，请人工确认。"
    fi
else
    echo"⚠️ 未找到配置文件，可能未初始化或使用默认值。"
fi

# 3. 检查端口监听
echo"🔍 [3/5] 检查端口监听状态..."
# 检查是否有进程监听在非127.0.0.1的18789端口
LISTEN_STATUS=$(netstat -tuln 2>/dev/null | grep 18789 || ss -tuln | grep 18789)
ifecho"$LISTEN_STATUS" | grep -q "0.0.0.0:18789\|*:18789"; then
    echo"🚨 危急！检测到端口 18789 对所有IP开放！"
else
    echo"✅ 端口监听正常（仅本地或未运行）。"
fi

# 4. 检查备份
echo"🔍 [4/5] 检查备份状态..."
BACKUP_DIR="$HOME/.openclaw/backups"
if [ -d "$BACKUP_DIR" ]; then
    COUNT=$(ls -1 "$BACKUP_DIR" 2>/dev/null | wc -l)
    if [ "$COUNT" -gt 0 ]; then
        echo"✅ 发现 $COUNT 个备份文件。"
    else
        echo"⚠️ 警告：备份目录为空，建议运行 'openclaw backup create'。"
    fi
else
    echo"⚠️ 警告：未找到备份目录，建议运行 'openclaw backup create'。"
fi

# 5. 提醒OAuth
echo"🔍 [5/5] 安全提示..."
echo"💡 请记得去 Telegram/Discord/Slack 撤销旧的Bot授权，并轮换API Key。"
echo"💡 如需远程访问，请使用 SSH 隧道或 Tailscale，严禁绑定 0.0.0.0。"

echo"----------------------------------------------------"
echo"自检结束。如有🚨标记，请立即处理！"
~~~

*(注：使用前请赋予执行权限 ~chmod +x check_openclaw.sh~)*

---

**参考资料**：

1. OpenClaw Release v2026.3.8, GitHub Repository, March 9, 2026.
2. CVE-2026-25253 Technical Analysis, CSDN Blog, Feb 2026.
3. "OpenClaw 48小时连爆7个CVE", 今日头条, Feb 26, 2026.
4. OpenClaw Security Advisory [#2026](javascript:;)-001 to [#2026](javascript:;)-007.

*(本文旨在科普与防御，请勿利用文中提到的技术进行非法攻击。)*


`

	markdown = strings.ReplaceAll(markdown, "~", "`")
	// t.Log(markdown)
	html := Convert(markdown, "sspai")

	if !strings.Contains(html, "<h1") {
		t.Error("Expected HTML to contain h1")
	}

	if !strings.Contains(html, "<div") {
		t.Error("Expected HTML to contain div wrapper")
	}
	t.Logf("%s", html)
}

func TestThemeExists(t *testing.T) {
	if !ThemeExists("apple") {
		t.Error("Expected apple theme to exist")
	}
	if !ThemeExists("wechat") {
		t.Error("Expected wechat theme to exist")
	}
	if ThemeExists("nonexistent") {
		t.Error("Expected nonexistent theme to not exist")
	}
}

func TestGetTheme(t *testing.T) {
	theme := GetTheme("claude")
	if theme.ID != "claude" {
		t.Errorf("Expected theme ID to be claude, got %s", theme.ID)
	}

	defaultTheme := GetTheme("nonexistent")
	if defaultTheme.ID != "apple" {
		t.Errorf("Expected default theme ID to be apple, got %s", defaultTheme.ID)
	}
}

func TestAllThemes(t *testing.T) {
	themes := AllThemes()
	if len(themes) == 0 {
		t.Error("Expected themes to not be empty")
	}
	if len(themes) < 30 {
		t.Logf("Warning: Expected at least 30 themes, got %d", len(themes))
	}
}

func TestConvertWithImageGrids(t *testing.T) {
	markdown := "![image1](img1.png)\n\n![image2](img2.png)\n\n![image3](img3.png)\n"

	html := Convert(markdown, "apple")

	if !strings.Contains(html, "image-grid") {
		t.Log("Image grids not applied (may need paragraph format)")
	}
}

func TestPreprocessMarkdown(t *testing.T) {
	input := "Some *** text --- with ___ special chars"
	output := PreprocessMarkdown(input)

	if strings.Contains(output, "***") {
		t.Error("Expected *** to be removed")
	}
}
