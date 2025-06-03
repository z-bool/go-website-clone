# 功能更新说明

## 新增功能

### 1. ConfigID 配置标识符
- **功能**: 在`goclone.Config`结构体中新增`ConfigID`字段
- **类型**: `string`
- **用途**: 使用UUID作为项目文件夹名称，替代原来基于域名的命名方式
- **特性**: 
  - 如果`ConfigID`为空，系统会自动生成UUID
  - 可以手动指定ConfigID
  - 确保每次克隆的项目都有唯一的标识符

### 2. MaxFolderSize 文件夹大小限制
- **功能**: 在`goclone.Config`结构体中新增`MaxFolderSize`字段
- **类型**: `int64` (字节数)
- **用途**: 限制克隆项目文件夹的最大大小
- **特性**:
  - 设置为0表示无大小限制
  - 在下载过程中实时检查文件夹大小
  - 超过限制时跳过后续文件下载
  - 提供详细的大小监控日志

## 代码结构优化

### 1. 新增文件管理功能
- `file.CreateProjectWithID()`: 使用指定ID创建项目目录
- `file.GetFolderSize()`: 计算文件夹总大小
- `file.CheckFolderSizeLimit()`: 检查文件夹大小是否超限

### 2. 爬虫功能增强
- `CrawlConfig`接口: 定义配置接口，避免循环依赖
- `CollectorWithSizeLimit()`: 带大小限制的资源收集器
- 实时大小检查机制

### 3. 接口设计改进
- 通过接口解耦，提高代码可维护性
- 支持可配置的爬取行为
- 更好的错误处理和日志输出

## 使用示例

```go
config := &goclone.Config{
    URLs:          []string{"https://example.com"},
    ConfigID:      "", // 自动生成UUID，或手动指定
    MaxFolderSize: 50 * 1024 * 1024, // 50MB限制
    UserAgent:     "Mozilla/5.0...",
    // 其他配置...
}

result := goclone.Clone(ctx, config)
fmt.Printf("项目保存在: %s, 配置ID: %s\n", result.FirstProject, config.ConfigID)
```

## 兼容性说明
- 现有的API保持兼容
- 新字段为可选配置
- 原有的`QuickClone`函数仍然可用 