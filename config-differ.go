package main

/*
Config Differ - Go版本
高性能配置文件对比工具

编译:
  go mod init config-differ
  go get gopkg.in/yaml.v3
  go build -o config-differ

使用:
  ./config-differ file1.yaml file2.yaml
  ./config-differ --ignore-order prod.json test.json
  ./config-differ -d /etc/nginx/prod /etc/nginx/test
*/

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/yaml.v3"
)

// ANSI颜色
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

type Config struct {
	IgnoreComments bool
	IgnoreBlank    bool
	IgnoreOrder    bool
	Context        int
}

func main() {
	var (
		file1          string
		file2          string
		isDir          bool
		pattern        string
		ignoreOrder    bool
		noIgnoreComments bool
		noIgnoreBlank  bool
		context        int
	)

	flag.BoolVar(&isDir, "d", false, "对比目录")
	flag.BoolVar(&isDir, "directory", false, "对比目录")
	flag.StringVar(&pattern, "p", "*", "文件模式")
	flag.StringVar(&pattern, "pattern", "*", "文件模式")
	flag.BoolVar(&ignoreOrder, "ignore-order", false, "忽略键顺序")
	flag.BoolVar(&noIgnoreComments, "no-ignore-comments", false, "不忽略注释")
	flag.BoolVar(&noIgnoreBlank, "no-ignore-blank", false, "不忽略空行")
	flag.IntVar(&context, "c", 3, "上下文行数")
	flag.IntVar(&context, "context", 3, "上下文行数")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s [选项] <file1> <file2>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s prod.yaml test.yaml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --ignore-order prod.json test.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -d /etc/nginx/prod /etc/nginx/test\n", os.Args[0])
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	file1 = args[0]
	file2 = args[1]

	config := Config{
		IgnoreComments: !noIgnoreComments,
		IgnoreBlank:    !noIgnoreBlank,
		IgnoreOrder:    ignoreOrder,
		Context:        context,
	}

	if isDir {
		compareDirs(file1, file2, pattern, config)
	} else {
		same, diff := compareFiles(file1, file2, config)
		fmt.Println(diff)
		if !same {
			os.Exit(1)
		}
	}
}

func detectType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".ini", ".conf", ".cfg":
		return "ini"
	default:
		return "text"
	}
}

func readFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func normalizeText(content string, fileType string, config Config) []string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		// 去除空行
		if config.IgnoreBlank && strings.TrimSpace(line) == "" {
			continue
		}

		// 去除注释
		if config.IgnoreComments {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "#") {
				continue
			}
			// 去除行尾注释
			if idx := strings.Index(line, "#"); idx > 0 {
				line = strings.TrimRight(line[:idx], " \t")
			}
		}

		result = append(result, line)
	}

	return result
}

func parseYAML(content string, ignoreOrder bool) (string, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		return "", err
	}

	if ignoreOrder {
		data = sortMapRecursive(data)
	}

	out, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func parseJSON(content string, ignoreOrder bool) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		return "", err
	}

	if ignoreOrder {
		data = sortMapRecursive(data)
	}

	var out []byte
	out, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func sortMapRecursive(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = sortMapRecursive(val)
		}
		return result
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			result[fmt.Sprint(k)] = sortMapRecursive(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = sortMapRecursive(val)
		}
		return result
	default:
		return v
	}
}

func compareFiles(file1, file2 string, config Config) (bool, string) {
	fileType := detectType(file1)

	fmt.Printf("%s文件类型: %s%s\n", ColorCyan, fileType, ColorReset)
	fmt.Printf("%s文件1: %s%s\n", ColorCyan, file1, ColorReset)
	fmt.Printf("%s文件2: %s%s\n", ColorCyan, file2, ColorReset)
	fmt.Printf("%s%s%s\n\n", ColorCyan, strings.Repeat("=", 60), ColorReset)

	content1, err := readFile(file1)
	if err != nil {
		return false, fmt.Sprintf("%s错误: %v%s", ColorRed, err, ColorReset)
	}

	content2, err := readFile(file2)
	if err != nil {
		return false, fmt.Sprintf("%s错误: %v%s", ColorRed, err, ColorReset)
	}

	var str1, str2 string

	switch fileType {
	case "yaml":
		str1, err = parseYAML(content1, config.IgnoreOrder)
		if err != nil {
			fmt.Printf("%s警告: YAML解析失败，回退到文本对比: %v%s\n", ColorYellow, err, ColorReset)
			return compareText(content1, content2, file1, file2, fileType, config)
		}
		str2, err = parseYAML(content2, config.IgnoreOrder)
		if err != nil {
			fmt.Printf("%s警告: YAML解析失败，回退到文本对比: %v%s\n", ColorYellow, err, ColorReset)
			return compareText(content1, content2, file1, file2, fileType, config)
		}

	case "json":
		str1, err = parseJSON(content1, config.IgnoreOrder)
		if err != nil {
			fmt.Printf("%s警告: JSON解析失败，回退到文本对比: %v%s\n", ColorYellow, err, ColorReset)
			return compareText(content1, content2, file1, file2, fileType, config)
		}
		str2, err = parseJSON(content2, config.IgnoreOrder)
		if err != nil {
			fmt.Printf("%s警告: JSON解析失败，回退到文本对比: %v%s\n", ColorYellow, err, ColorReset)
			return compareText(content1, content2, file1, file2, fileType, config)
		}

	default:
		return compareText(content1, content2, file1, file2, fileType, config)
	}

	return unifiedDiff(str1, str2, file1, file2, config.Context)
}

func compareText(content1, content2, file1, file2, fileType string, config Config) (bool, string) {
	lines1 := normalizeText(content1, fileType, config)
	lines2 := normalizeText(content2, fileType, config)

	str1 := strings.Join(lines1, "\n")
	str2 := strings.Join(lines2, "\n")

	return unifiedDiff(str1, str2, file1, file2, config.Context)
}

func unifiedDiff(str1, str2, file1, file2 string, context int) (bool, string) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(str1),
		B:        difflib.SplitLines(str2),
		FromFile: file1,
		ToFile:   file2,
		Context:  context,
	}

	result, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return false, fmt.Sprintf("%s错误: %v%s", ColorRed, err, ColorReset)
	}

	if result == "" {
		return true, fmt.Sprintf("%s✓ 文件完全相同%s", ColorGreen, ColorReset)
	}

	// 彩色化输出
	colored := colorizeDiff(result)
	return false, colored
}

func colorizeDiff(diff string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(diff))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			result.WriteString(ColorCyan + ColorBold + line + ColorReset + "\n")
		} else if strings.HasPrefix(line, "@@") {
			result.WriteString(ColorYellow + line + ColorReset + "\n")
		} else if strings.HasPrefix(line, "+") {
			result.WriteString(ColorGreen + line + ColorReset + "\n")
		} else if strings.HasPrefix(line, "-") {
			result.WriteString(ColorRed + line + ColorReset + "\n")
		} else {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}

func compareDirs(dir1, dir2, pattern string, config Config) {
	files1, _ := filepath.Glob(filepath.Join(dir1, pattern))
	files2, _ := filepath.Glob(filepath.Join(dir2, pattern))

	// 提取文件名
	names1 := make(map[string]bool)
	for _, f := range files1 {
		names1[filepath.Base(f)] = true
	}

	names2 := make(map[string]bool)
	for _, f := range files2 {
		names2[filepath.Base(f)] = true
	}

	// 找出共有文件
	var common []string
	for name := range names1 {
		if names2[name] {
			common = append(common, name)
		}
	}
	sort.Strings(common)

	fmt.Printf("\n%s对比共有文件 (%d 个):%s\n\n", ColorCyan, len(common), ColorReset)

	sameCount := 0
	diffCount := 0

	for _, name := range common {
		file1 := filepath.Join(dir1, name)
		file2 := filepath.Join(dir2, name)

		fmt.Printf("\n%s%s%s\n", ColorCyan, strings.Repeat("=", 60), ColorReset)
		fmt.Printf("%s对比: %s%s\n", ColorCyan, name, ColorReset)
		fmt.Printf("%s%s%s\n\n", ColorCyan, strings.Repeat("=", 60), ColorReset)

		same, diff := compareFiles(file1, file2, config)
		fmt.Println(diff)

		if same {
			sameCount++
		} else {
			diffCount++
		}
	}

	fmt.Printf("\n%s%s%s\n", ColorCyan, strings.Repeat("=", 60), ColorReset)
	fmt.Printf("%s汇总%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s%s%s\n", ColorCyan, strings.Repeat("=", 60), ColorReset)
	fmt.Printf("%s相同: %d%s\n", ColorGreen, sameCount, ColorReset)
	fmt.Printf("%s差异: %d%s\n", ColorRed, diffCount, ColorReset)
	fmt.Printf("总计: %d\n", sameCount+diffCount)
}

/*
go.mod:

module config-differ

go 1.21

require (
	github.com/pmezard/go-difflib v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)
*/
