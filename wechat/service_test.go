package wechat

import (
	"os"
	"testing"

	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
)

var htmlContent string = `
<div style="max-width: 100%; margin: 0 auto; padding: 24px 20px 48px 20px; font-family: -apple-system, BlinkMacSystemFont, &quot;Segoe UI&quot;, Roboto, &quot;Helvetica Neue&quot;, Arial, sans-serif; font-size: 16px; line-height: 1.7 !important; color: #2d4a3e !important; background-color: #f0faf5 !important; word-wrap: break-word;"><h1 style="; font-size: 32px; font-weight: 700; color: #1a7a5a !important; line-height: 1.3 !important; margin: 38px 0 16px; letter-spacing: -0.015em;">企业AI团队必备的8大LLM开发技能</h1>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">当组织讨论采用大语言模型时，对话通常从模型选择开始：GPT还是Claude？开源还是闭源？更大还是更便宜？</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">但在真实的企业系统中，这种关注点是错位的。LLM的生产成功更多取决于架构纪律，而非模型本身。将脆弱的演示与弹性、可治理的系统区分开来的，是对少数核心工程技能的掌握。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">这些技能决定了模型如何被指令、如何被锚定、如何部署、如何被观测，以及如何随时间演进。本文将从构建真实系统（而非在笔记本中实验）的角度，讨论这八项技能。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">1. 提示工程（Prompt Engineering）</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879730-1769965621587.png" alt="提示工程，结构化指令流" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">提示工程是任何LLM系统的基础层。它将人类意图转化为模型可以可靠执行的精确、结构化指令。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">在生产环境中，提示词不是手写的字符串。它们使用模板、角色、约束、示例和安全规则程序化组装。强大的提示工程可以减少幻觉、提高一致性，通常还能延迟对微调或Agent等更复杂方法的需求。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">相反，糟糕的提示词会放大变异性，迫使团队用脆弱的下游逻辑来弥补。在成熟系统中，提示词像应用代码一样被版本控制、测试和审查。这种纪律使团队能够在不重写业务逻辑的情况下更换模型。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">2. 上下文工程（Context Engineering）</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879731-1769965657911.png" alt="上下文工程，确定性上下文组装" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">上下文工程决定了模型在推理时看到什么信息。系统不再将所有内容塞进单一提示词，而是从内存存储、结构化数据库、文档和API动态组装相关上下文。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">这是企业可靠性的真正起点。上下文工程是确定性和可审计的。你可以解释模型为什么那样响应，因为你确切知道给了它什么数据。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">跳过这一步的团队往往依赖模型推断缺失信息。这种方法在演示中可能有效，但在监管审查或运营规模下会失败。上下工程将LLM从概率猜测者转变为受控的推理组件。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">3. 微调（Fine-Tuning）</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879732-1769965911027.png" alt="微调，受控的模型适配" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">微调修改模型本身，使期望的行为被内化而非被反复指令。这种方法在相同任务大规模重复时最有效，如分类、提取或领域特定推理。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">代价是灵活性。微调后的模型更难改变，需要严格的数据治理。训练数据必须经过筛选、版本控制，并审查偏见和漂移。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">在企业环境中，微调应该是一个深思熟虑的优化步骤，而非默认起点。许多团队过早微调，而提示工程和上下文工程本已足够。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">4. 检索增强生成（RAG）</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879733-1769965996926.png" alt="检索增强生成" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">检索增强生成（RAG）将模型输出锚定在外部知识中。系统不再信任模型记忆的内容，而是在运行时检索相关信息并注入提示词。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">这种模式主导了企业采用，因为它平衡了准确性、新鲜度和可解释性。知识可以更新而无需重新训练模型，响应可以追溯到源文档。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">设计良好的RAG系统将检索视为一等公民。分块策略、嵌入选择、排序逻辑和上下文都会显著影响结果质量。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">5. Agent智能体</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879734-1769966072776.png" alt="Agent架构" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">Agent引入了自主性。Agent不只是响应输入——它推理、规划、调用工具、评估结果并迭代直到达成目标。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">这种能力强大但若误用则危险。Agent最适合多步骤分析、编排和决策支持等工作流。它们不适合事实检索或合规敏感的输出。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">常见失败模式包括无限循环、工具幻觉、成本失控和不可预测的行为。在企业系统中，Agent必须用明确目标、步骤限制、工具白名单和强可观测性来约束。没有护栏的自主性不是智能，而是风险。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">6. LLM部署</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879735-1769966178360.png" alt="LLM部署" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">部署将模型转化为可靠的服务。这一层处理路由、可扩展性、认证、授权和版本控制。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">干净的部署架构让团队能够在不强制应用变更的情况下更换模型。在企业环境中，部署还定义了安全边界——数据流向何处、请求如何记录、故障如何隔离。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">将LLM视为普通的API依赖是错误的。它们是需要谨慎暴露和生命周期管理的概率系统。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">7. LLM优化</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879738-1769966266899.png" alt="LLM优化" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">优化确保规模化的性能和成本效率。包括缓存频繁响应、压缩上下文、将请求路由到不同模型，以及应用量化等技术。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">优化对终端用户通常不可见，但对可持续性至关重要。没有它，即使设计良好的系统也会随着使用增长变得过于昂贵。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">团队应将优化视为持续纪律而非一次性工作。使用模式在演进，优化策略也应如此。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">8. LLM可观测性</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879739-1769966354797.png" alt="LLM可观测性、治理和监控" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">可观测性提供对提示词、响应、延迟、成本和失败模式的可见性。没有它，LLM系统实际上是无法治理的。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">在受监管行业，可观测性不是可选项。团队必须能够追踪输出、审计决策、检测漂移或滥用。有效的可观测性结合了追踪、指标和结构化日志，使团队能够调试行为、执行策略并持续改进系统质量。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">企业LLM参考架构</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><img src="https://dz2cdn1.dzone.com/storage/temp/18879742-1769967188158.png" alt="企业LLM参考架构" style="; max-width: 100%; height: auto; display: block; margin: 24px auto; border-radius: 8px;; display:block; width:100%; max-width:100%; height:auto; margin:30px auto !important; padding:8px !important; border-radius:14px !important; box-sizing:border-box; box-shadow:0 16px 34px rgba(15,23,42,0.22), 0 4px 10px rgba(15,23,42,0.12); border:1px solid rgba(15,23,42,0.12);"></p>
<h3 style="; font-size: 21px; font-weight: 600; color: #2d4a3e !important; line-height: 1.4 !important; margin: 28px 0 14px;">参考架构说明</h3>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">1. 提示工程</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">系统的基础是提示工程，用户意图被转化为结构化指令。在生产中，提示词使用模板、系统角色、约束和示例程序化组装。LangChain和LlamaIndex等工具支持模块化、可复用的提示模板，提高确定性并减少幻觉。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">2. 上下文工程</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">上下文工程确保模型在推理时看到正确信息。系统从多个来源动态组装相关上下文，包括：内存数据库（Redis、DynamoDB）、文档存储（S3、Blob Storage）、结构化企业数据（Postgres、Snowflake）、向量存储（Pinecone、Weaviate、FAISS）。上下文构建器对数据进行排名和过滤，为模型提供确定性、可审计的输入。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">3. 微调</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">微调定制模型以内化重复任务的行为。这对分类、提取或大规模推理等领域特定任务至关重要。使用HuggingFace或SageMaker等平台实现微调模型，以灵活性为代价换取一致性。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">4. 检索增强生成</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">RAG确保输出锚定在外部知识而非仅依赖模型记忆。系统从向量存储、文档仓库、企业数据库和内存层检索相关信息并嵌入提示词。这平衡了准确性、新鲜度和可解释性，是企业可靠性的关键组成部分。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">5. Agent和工具</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">Agent编排自主推理和任务执行。它们决策、迭代并调用外部工具以达成目标。企业工具包括搜索、SQL查询、Python脚本或API。Agent在原始推理之上提供结构化工作流层，在保持控制和可审计性的同时实现复杂多步操作。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">6. 模型和推理</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">这一层管理基础模型和微调模型的执行。模型路由器根据成本、延迟或其他标准选择适当模型。基础模型（GPT 4.x、Claude、Gemini）处理通用任务，微调模型执行领域特定操作。这一层将模型转化为可扩展、可演进且无需更改应用逻辑的可靠服务。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">7. 优化</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">优化确保性能和成本效率。技术包括：使用Redis进行响应缓存、通过摘要和分块进行上下文压缩、INT8或INT4量化。这些优化对用户不可见但对规模化的可持续性至关重要。</p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><strong style="; font-weight: 700; color: #1a7a5a !important;">8. 可观测性和治理</strong></p>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">最后一层提供可见性、可追溯性和监控。OpenTelemetry和LangSmith追踪提示和模型活动。Prometheus或Datadog处理指标和成本追踪。ELK或CloudWatch收集提示和模型输出日志，Grafana等仪表盘为工程师和决策者提供全面视图。可观测性实现治理、审计和运营可靠性。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<h2 style="; font-size: 26px; font-weight: 600; color: #1a7a5a !important; line-height: 1.35 !important; margin: 32px 0 16px;">总结</h2>
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;">当你理解这八项技能及其组合方式，你不再以模型思维思考，而是开始以系统思维思考。这种转变将LLM采用从实验转变为真正的工程领导力。</p>
<hr style="; margin: 36px auto; border: none; height: 1px; background-color: #c8e6d8 !important; width: 100%;">
<p style="; margin: 18px 0 !important; line-height: 1.7 !important; color: #2d4a3e !important;"><em style="; font-style: italic; color: #5a8a72 !important;">本文翻译整理自 DZone</em></p>
</div>
`

func Test_CreateDraft(t *testing.T) {
	cfg := &config.Config{
		WechatAppID:     os.Getenv("WECHAT_APP_ID"),
		WechatAppSecret: os.Getenv("WECHAT_APP_SECRET"),
	}

	service := NewService(cfg)
	imagePath, err := DownloadFile("https://dz2cdn1.dzone.com/storage/temp/18879732-1769965911027.png")
	if err != nil || imagePath == "" {
		t.Fatal(err)
	}
	ur, err := service.UploadMaterial(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("media :%s, %s \n", ur.MediaID, ur.WechatURL)

	cr, err := service.CreateDraft([]*draft.Article{
		{
			Title:              "测试",
			Author:             "梦朝思夕",
			Digest:             "文章摘要",
			ThumbMediaID:       ur.MediaID,
			Content:            htmlContent,
			ShowCoverPic:       1,
			NeedOpenComment:    1,
			OnlyFansCanComment: 0,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("draft :%s ,%s\n", cr.MediaID, cr.DraftURL)

}
