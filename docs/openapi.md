将aiConverter中的AI模型提取出来，以便后续其他模块也可以使用。
ai供应商除了ANTHROPIC,还需要提供openai的能力。
设计要求：
1. 使用设计模式
2. 使用优雅简洁
3. 便于后续的复用
4. 与config实现联动，可以配置ai供应商
可以参看代码：https://github.com/nextlevelbuilder/goclaw/blob/main/internal/providers/openai.go