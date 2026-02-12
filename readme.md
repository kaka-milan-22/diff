# 配置文件 Diff 命令行工具指南

## 目录
- [基础工具](#基础工具)
- [高级工具](#高级工具)
- [实用脚本](#实用脚本)
- [多环境对比](#多环境对比)
- [自动化集成](#自动化集成)

---

## 基础工具

### 1. diff - 经典工具
```bash
# 基础对比
diff file1.conf file2.conf

# 并排显示（推荐）
diff -y file1.conf file2.conf

# 彩色输出
diff --color=always file1.conf file2.conf

# 忽略空白字符
diff -w file1.conf file2.conf

# 上下文模式（显示周围3行）
diff -u file1.conf file2.conf

# 递归对比目录
diff -r /etc/nginx/prod /etc/nginx/test

# 只显示差异文件名
diff -qr /etc/nginx/prod /etc/nginx/test
```

**使用场景**: 快速对比两个文件，Linux自带无需安装

---

### 2. vimdiff - 可视化编辑
```bash
# 并排对比（Vim用户最爱）
vimdiff prod.yaml test.yaml

# 三向对比
vimdiff base.conf prod.conf test.conf

# 水平分割
vim -d -o prod.conf test.conf

# 垂直分割
vim -d -O prod.conf test.conf
```

**快捷键**:
- `]c` - 跳到下一个差异
- `[c` - 跳到上一个差异
- `do` - (diff obtain) 从另一个文件获取差异
- `dp` - (diff put) 将差异放到另一个文件
- `:diffupdate` - 更新差异高亮
- `zo` / `zc` - 展开/折叠相同内容

**使用场景**: 需要边对比边编辑，可视化操作

---

### 3. colordiff - diff彩色版
```bash
# 安装
sudo apt-get install colordiff  # Ubuntu/Debian
sudo yum install colordiff      # CentOS/RHEL
brew install colordiff          # macOS

# 使用（语法同diff）
colordiff -u prod.conf test.conf

# 结合less分页
colordiff -u prod.conf test.conf | less -R

# 并排彩色
colordiff -y prod.conf test.conf | less -R

# 设置别名
alias diff='colordiff'
```

**使用场景**: 让diff输出更易读，彩色标记差异

---

### 4. git diff - 不仅仅是Git
```bash
# 即使不是Git仓库也能用
git diff --no-index prod.conf test.conf

# 彩色、词级差异
git diff --no-index --color-words prod.conf test.conf

# 并排显示
git diff --no-index --word-diff=color prod.conf test.conf

# 统计差异
git diff --no-index --stat prod.conf test.conf

# 忽略空白
git diff --no-index -w prod.conf test.conf
```

**使用场景**: 最佳的彩色diff，词级高亮

---

## 高级工具

### 5. delta - 现代化diff查看器
```bash
# 安装
brew install git-delta
cargo install git-delta

# Ubuntu/Debian
wget https://github.com/dandavison/delta/releases/download/0.16.5/git-delta_0.16.5_amd64.deb
sudo dpkg -i git-delta_0.16.5_amd64.deb

# 使用
delta prod.conf test.conf

# 配合git
git diff --no-index prod.conf test.conf | delta

# 并排模式
delta --side-by-side prod.conf test.conf

# 配置为git默认
git config --global core.pager delta
git config --global interactive.diffFilter "delta --color-only"
git config --global delta.navigate true
git config --global delta.side-by-side true
```

**特性**:
- 语法高亮
- 行内词级差异
- Git集成
- 主题支持

**使用场景**: 需要最佳可读性的diff查看

---

### 6. icdiff - 改进的并排diff
```bash
# 安装
pip install icdiff --break-system-packages

# 基础使用
icdiff prod.conf test.conf

# 全屏并排
icdiff --whole-file prod.conf test.conf

# 指定列宽
icdiff --cols=160 prod.conf test.conf

# 高亮模式
icdiff --highlight prod.conf test.conf

# 递归对比目录
icdiff -r /etc/nginx/prod /etc/nginx/test
```

**使用场景**: 需要清晰的并排对比，Python环境

---

### 7. dyff - YAML/JSON专用
```bash
# 安装
brew install dyff
go install github.com/homeport/dyff/cmd/dyff@latest

# YAML对比
dyff between prod.yaml test.yaml

# JSON对比
dyff between prod.json test.json

# 只显示差异路径
dyff between --output brief prod.yaml test.yaml

# 忽略顺序
dyff between --ignore-order-changes prod.yaml test.yaml

# 输出格式
dyff between --output human prod.yaml test.yaml     # 人类可读
dyff between --output json prod.yaml test.yaml      # JSON格式
dyff between --output yaml prod.yaml test.yaml      # YAML格式
```

**特性**:
- 理解YAML/JSON结构
- 忽略注释和顺序
- 显示路径层级
- 彩色输出

**使用场景**: Kubernetes配置、微服务配置对比

---

### 8. jd - JSON专用diff
```bash
# 安装
npm install -g jd-cli
go install github.com/josephburnett/jd@latest

# 使用
jd prod.json test.json

# 输出格式
jd -f patch prod.json test.json     # JSON patch格式
jd -f merge prod.json test.json     # JSON merge patch
jd -f diff prod.json test.json      # 自定义格式

# 只显示差异
jd -set prod.json test.json
```

**使用场景**: API响应、配置文件JSON对比

---

### 9. yq + diff - YAML处理
```bash
# 安装yq
brew install yq
wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/local/bin/yq

# 提取特定字段对比
diff <(yq eval '.database' prod.yaml) <(yq eval '.database' test.yaml)

# 排序后对比（忽略顺序）
diff <(yq eval 'sort_keys(..)' prod.yaml) <(yq eval 'sort_keys(..)' test.yaml)

# 转JSON后对比
diff <(yq eval -o=json prod.yaml) <(yq eval -o=json test.yaml)

# 合并差异
yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' prod.yaml test.yaml
```

**使用场景**: 复杂YAML结构对比、Kubernetes清单

---

### 10. daff - 表格数据diff
```bash
# 安装
npm install -g daff
pip install daff --break-system-packages

# CSV对比
daff prod.csv test.csv

# 高亮差异单元格
daff --output diff.csv prod.csv test.csv

# HTML输出
daff --output diff.html prod.csv test.csv
```

**使用场景**: CSV、数据库导出文件对比

---

## 实用脚本

### 11. 智能配置对比脚本
```bash
#!/bin/bash
# smart-diff.sh - 根据文件类型选择最佳diff工具

FILE1="$1"
FILE2="$2"

if [ -z "$FILE1" ] || [ -z "$FILE2" ]; then
    echo "用法: $0 <file1> <file2>"
    exit 1
fi

# 检测文件类型
EXT1="${FILE1##*.}"
EXT2="${FILE2##*.}"

# 根据扩展名选择工具
case "$EXT1" in
    yaml|yml)
        if command -v dyff &>/dev/null; then
            dyff between "$FILE1" "$FILE2"
        elif command -v yq &>/dev/null; then
            diff -u <(yq eval 'sort_keys(..)' "$FILE1") <(yq eval 'sort_keys(..)' "$FILE2")
        else
            git diff --no-index --color "$FILE1" "$FILE2"
        fi
        ;;
    
    json)
        if command -v jd &>/dev/null; then
            jd "$FILE1" "$FILE2"
        elif command -v jq &>/dev/null; then
            diff -u <(jq -S . "$FILE1") <(jq -S . "$FILE2")
        else
            git diff --no-index --color "$FILE1" "$FILE2"
        fi
        ;;
    
    csv)
        if command -v daff &>/dev/null; then
            daff "$FILE1" "$FILE2"
        else
            diff -y "$FILE1" "$FILE2" | less
        fi
        ;;
    
    ini|conf|cfg)
        if command -v icdiff &>/dev/null; then
            icdiff "$FILE1" "$FILE2"
        elif command -v delta &>/dev/null; then
            delta "$FILE1" "$FILE2"
        else
            colordiff -y "$FILE1" "$FILE2" | less -R
        fi
        ;;
    
    *)
        # 默认使用git diff或delta
        if command -v delta &>/dev/null; then
            git diff --no-index "$FILE1" "$FILE2" | delta
        elif command -v icdiff &>/dev/null; then
            icdiff "$FILE1" "$FILE2"
        else
            git diff --no-index --color "$FILE1" "$FILE2" | less -R
        fi
        ;;
esac
```

---

### 12. 多环境批量对比
```bash
#!/bin/bash
# multi-env-diff.sh - 对比多个环境的配置

ENVIRONMENTS=("dev" "test" "prod")
CONFIG_DIR="/etc/nginx/sites-enabled"
CONFIGS=("api.conf" "web.conf" "admin.conf")

echo "=== 多环境配置对比报告 ==="
echo "时间: $(date)"
echo ""

for config in "${CONFIGS[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "配置: $config"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # dev vs test
    echo "▶ dev vs test:"
    if diff -q "${CONFIG_DIR}/dev/${config}" "${CONFIG_DIR}/test/${config}" &>/dev/null; then
        echo "  ✓ 相同"
    else
        echo "  ✗ 差异:"
        diff -u --color=always "${CONFIG_DIR}/dev/${config}" "${CONFIG_DIR}/test/${config}" | head -20
    fi
    echo ""
    
    # test vs prod
    echo "▶ test vs prod:"
    if diff -q "${CONFIG_DIR}/test/${config}" "${CONFIG_DIR}/prod/${config}" &>/dev/null; then
        echo "  ✓ 相同"
    else
        echo "  ✗ 差异:"
        diff -u --color=always "${CONFIG_DIR}/test/${config}" "${CONFIG_DIR}/prod/${config}" | head -20
    fi
    echo ""
done
```

---

### 13. 配置变更追踪
```bash
#!/bin/bash
# config-watch.sh - 监控配置文件变化

CONFIG_FILE="/etc/nginx/nginx.conf"
SNAPSHOT_DIR="/var/backups/config-snapshots"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "$SNAPSHOT_DIR"

# 创建快照
cp "$CONFIG_FILE" "${SNAPSHOT_DIR}/nginx.conf.${TIMESTAMP}"

# 查找最近的快照
LAST_SNAPSHOT=$(ls -t "${SNAPSHOT_DIR}"/nginx.conf.* | head -2 | tail -1)

if [ -n "$LAST_SNAPSHOT" ]; then
    echo "=== 配置变更检测 ==="
    echo "当前: $CONFIG_FILE"
    echo "对比: $LAST_SNAPSHOT"
    echo ""
    
    if diff -q "$CONFIG_FILE" "$LAST_SNAPSHOT" &>/dev/null; then
        echo "✓ 无变化"
    else
        echo "✗ 发现变更:"
        diff -u --color=always "$LAST_SNAPSHOT" "$CONFIG_FILE"
        
        # 可选：发送告警
        # echo "配置文件已变更" | mail -s "Config Change Alert" ops@company.com
    fi
fi

# 清理旧快照（保留最近30个）
ls -t "${SNAPSHOT_DIR}"/nginx.conf.* | tail -n +31 | xargs rm -f
```

---

### 14. 远程配置对比
```bash
#!/bin/bash
# remote-diff.sh - 对比本地和远程服务器配置

LOCAL_FILE="$1"
REMOTE_HOST="$2"
REMOTE_FILE="$3"

if [ -z "$LOCAL_FILE" ] || [ -z "$REMOTE_HOST" ] || [ -z "$REMOTE_FILE" ]; then
    echo "用法: $0 <local_file> <remote_host> <remote_file>"
    echo "示例: $0 nginx.conf prod-web-01 /etc/nginx/nginx.conf"
    exit 1
fi

# 获取远程文件
TEMP_FILE=$(mktemp)
scp "${REMOTE_HOST}:${REMOTE_FILE}" "$TEMP_FILE"

echo "=== 本地 vs 远程配置对比 ==="
echo "本地: $LOCAL_FILE"
echo "远程: ${REMOTE_HOST}:${REMOTE_FILE}"
echo ""

# 对比
if command -v delta &>/dev/null; then
    git diff --no-index "$LOCAL_FILE" "$TEMP_FILE" | delta
elif command -v icdiff &>/dev/null; then
    icdiff "$LOCAL_FILE" "$TEMP_FILE"
else
    diff -u --color=always "$LOCAL_FILE" "$TEMP_FILE"
fi

# 清理
rm -f "$TEMP_FILE"
```

---

### 15. 批量服务器配置对比
```bash
#!/bin/bash
# cluster-diff.sh - 对比集群内所有服务器的配置一致性

CONFIG_FILE="/etc/nginx/nginx.conf"
SERVERS=("web-01" "web-02" "web-03")
TEMP_DIR=$(mktemp -d)

echo "=== 集群配置一致性检查 ==="
echo "配置: $CONFIG_FILE"
echo ""

# 下载所有服务器的配置
for server in "${SERVERS[@]}"; do
    echo "获取 $server 配置..."
    scp "${server}:${CONFIG_FILE}" "${TEMP_DIR}/${server}.conf" 2>/dev/null
done

# 计算MD5
echo ""
echo "MD5校验:"
for server in "${SERVERS[@]}"; do
    if [ -f "${TEMP_DIR}/${server}.conf" ]; then
        md5=$(md5sum "${TEMP_DIR}/${server}.conf" | awk '{print $1}')
        echo "$server: $md5"
    fi
done

# 对比差异
echo ""
echo "差异对比:"
REFERENCE="${SERVERS[0]}"
for server in "${SERVERS[@]:1}"; do
    echo ""
    echo "━━━ $REFERENCE vs $server ━━━"
    if diff -q "${TEMP_DIR}/${REFERENCE}.conf" "${TEMP_DIR}/${server}.conf" &>/dev/null; then
        echo "✓ 配置一致"
    else
        echo "✗ 发现差异:"
        diff -u --color=always "${TEMP_DIR}/${REFERENCE}.conf" "${TEMP_DIR}/${server}.conf" | head -30
    fi
done

# 清理
rm -rf "$TEMP_DIR"
```

---

## 多环境对比

### 16. 结构化对比（YAML/JSON）
```bash
# Kubernetes ConfigMap对比
diff <(kubectl get cm my-config -n prod -o yaml | yq eval 'del(.metadata)' -) \
     <(kubectl get cm my-config -n test -o yaml | yq eval 'del(.metadata)' -)

# 只对比data部分
diff <(kubectl get cm my-config -n prod -o jsonpath='{.data}' | jq -S .) \
     <(kubectl get cm my-config -n test -o jsonpath='{.data}' | jq -S .)

# 对比Secret（解码后）
diff <(kubectl get secret my-secret -n prod -o json | jq -r '.data | map_values(@base64d)') \
     <(kubectl get secret my-secret -n test -o json | jq -r '.data | map_values(@base64d)')
```

---

### 17. Nginx配置对比
```bash
# 对比并忽略注释
diff <(grep -v '^#' /etc/nginx/prod.conf | grep -v '^$') \
     <(grep -v '^#' /etc/nginx/test.conf | grep -v '^$')

# 对比有效配置（测试语法）
diff <(nginx -T -c /etc/nginx/prod.conf 2>/dev/null | grep -v '^#') \
     <(nginx -T -c /etc/nginx/test.conf 2>/dev/null | grep -v '^#')

# 提取特定块对比
diff <(sed -n '/server {/,/}/p' /etc/nginx/prod.conf) \
     <(sed -n '/server {/,/}/p' /etc/nginx/test.conf)
```

---

### 18. 数据库配置对比
```bash
# MySQL配置
diff <(grep -v '^#' /etc/mysql/prod.cnf | sort) \
     <(grep -v '^#' /etc/mysql/test.cnf | sort)

# PostgreSQL
diff <(grep -v '^#' /etc/postgresql/prod/postgresql.conf | sort) \
     <(grep -v '^#' /etc/postgresql/test/postgresql.conf | sort)

# Redis
diff <(redis-cli -h prod-redis CONFIG GET '*' | paste - -) \
     <(redis-cli -h test-redis CONFIG GET '*' | paste - -)
```

---

## 自动化集成

### 19. Git钩子自动对比
```bash
# .git/hooks/pre-commit
#!/bin/bash

# 检查配置变更
CHANGED_CONFIGS=$(git diff --cached --name-only | grep -E '\.(yaml|json|conf)$')

if [ -n "$CHANGED_CONFIGS" ]; then
    echo "=== 配置文件变更检测 ==="
    for file in $CHANGED_CONFIGS; do
        echo ""
        echo "━━━ $file ━━━"
        git diff --cached --color "$file" | delta
    done
    
    read -p "确认提交这些变更? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi
```

---

### 20. CI/CD配置验证
```yaml
# .gitlab-ci.yml
config-diff:
  stage: test
  script:
    - |
      echo "对比生产环境配置差异"
      for file in configs/*.yaml; do
        filename=$(basename $file)
        if kubectl get cm ${filename%.yaml} -n prod &>/dev/null; then
          echo "检查: $filename"
          diff <(yq eval 'sort_keys(..)' $file) \
               <(kubectl get cm ${filename%.yaml} -n prod -o yaml | yq eval '.data' - | yq eval 'sort_keys(..)' -)
        fi
      done
  only:
    changes:
      - configs/**/*.yaml
```

---

### 21. Prometheus告警配置对比
```bash
#!/bin/bash
# 对比Prometheus配置

PROD_URL="http://prometheus-prod:9090"
TEST_URL="http://prometheus-test:9090"

# 获取告警规则
curl -s "${PROD_URL}/api/v1/rules" | jq -S '.data.groups' > /tmp/prod-rules.json
curl -s "${TEST_URL}/api/v1/rules" | jq -S '.data.groups' > /tmp/test-rules.json

# 对比
jd /tmp/prod-rules.json /tmp/test-rules.json
```

---

## 工具选择决策树

```
配置文件类型?
│
├─ YAML/JSON
│  ├─ Kubernetes? → dyff + kubectl diff
│  ├─ 复杂嵌套? → dyff
│  └─ 简单对比? → yq + diff / jq + diff
│
├─ INI/CONF
│  ├─ 需要编辑? → vimdiff
│  ├─ 快速查看? → icdiff / delta
│  └─ 自动化? → diff / git diff
│
├─ CSV/表格
│  └─ daff
│
└─ 其他
   ├─ 最佳可读性? → delta
   ├─ 简单快速? → colordiff
   └─ 原生工具? → diff
```

---

## 常用别名配置

```bash
# ~/.bashrc 或 ~/.zshrc

# 基础别名
alias diff='colordiff'
alias vdiff='vimdiff'

# 智能diff
alias sdiff='smart-diff.sh'  # 使用上面的脚本

# YAML对比
alias ydiff='dyff between'
alias ydiff-sort='diff <(yq eval "sort_keys(..)" $1) <(yq eval "sort_keys(..)" $2)'

# JSON对比
alias jdiff='jd'
alias jdiff-sort='diff <(jq -S . $1) <(jq -S . $2)'

# 并排对比
alias diffside='icdiff'

# Git diff增强
alias gd='git diff --no-index --color-words'

# 远程对比
rdiff() {
    local file=$1
    local host=$2
    local remote_file=${3:-$file}
    diff <(cat $file) <(ssh $host cat $remote_file)
}

# K8s配置对比
kdiff() {
    local resource=$1
    local name=$2
    diff <(kubectl get $resource $name -n prod -o yaml | yq eval 'del(.metadata)' -) \
         <(kubectl get $resource $name -n test -o yaml | yq eval 'del(.metadata)' -)
}
```

---

## 推荐工具组合

### 日常使用（轻量）
```bash
brew install colordiff git-delta
pip install icdiff --break-system-packages
```

### 运维专业（全能）
```bash
# 安装所有工具
brew install colordiff git-delta dyff yq jq
pip install icdiff daff --break-system-packages
go install github.com/josephburnett/jd@latest
```

### 云原生场景
```bash
brew install dyff yq
# 配合kubectl diff使用
```

---

## 性能对比

| 工具 | 大文件 | 大目录 | 实时性 | 可读性 |
|------|--------|--------|--------|--------|
| diff | ★★★★★ | ★★★★★ | ★★★★★ | ★★☆☆☆ |
| colordiff | ★★★★☆ | ★★★★☆ | ★★★★☆ | ★★★★☆ |
| git diff | ★★★★☆ | ★★★★☆ | ★★★★☆ | ★★★★★ |
| delta | ★★★☆☆ | ★★★☆☆ | ★★★☆☆ | ★★★★★ |
| icdiff | ★★★☆☆ | ★★☆☆☆ | ★★★☆☆ | ★★★★☆ |
| dyff | ★★★★☆ | N/A | ★★★★☆ | ★★★★★ |
| vimdiff | ★★☆☆☆ | N/A | ★★★☆☆ | ★★★★☆ |

---

**更新日期**: 2024-02-12  
**适用场景**: Linux/macOS运维  
**推荐组合**: colordiff + delta + dyff + 自定义脚本
