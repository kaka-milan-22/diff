#!/usr/bin/env python3
"""
Config Differ - 智能配置文件对比工具

支持多种格式:
- YAML/YML
- JSON
- INI/CONF
- 纯文本

功能:
- 自动识别文件类型
- 结构化对比（YAML/JSON）
- 忽略注释和空行
- 彩色输出
- 支持远程文件
- 多环境批量对比

作者: DevOps Team
"""

import argparse
import json
import yaml
import sys
import os
import subprocess
import tempfile
from pathlib import Path
from typing import Dict, List, Tuple, Optional
from difflib import unified_diff, HtmlDiff
from colorama import init, Fore, Style
import configparser

init(autoreset=True)


class ConfigDiffer:
    def __init__(self, ignore_comments=True, ignore_blank=True, 
                 ignore_order=False, context_lines=3):
        self.ignore_comments = ignore_comments
        self.ignore_blank = ignore_blank
        self.ignore_order = ignore_order
        self.context_lines = context_lines
    
    def detect_type(self, filepath: str) -> str:
        """检测文件类型"""
        ext = Path(filepath).suffix.lower()
        
        type_map = {
            '.yaml': 'yaml',
            '.yml': 'yaml',
            '.json': 'json',
            '.ini': 'ini',
            '.conf': 'ini',
            '.cfg': 'ini',
            '.properties': 'properties',
            '.env': 'env',
        }
        
        return type_map.get(ext, 'text')
    
    def read_file(self, filepath: str) -> str:
        """读取文件内容"""
        if filepath.startswith('http://') or filepath.startswith('https://'):
            # 远程文件
            import urllib.request
            with urllib.request.urlopen(filepath) as response:
                return response.read().decode('utf-8')
        
        if ':' in filepath and not filepath.startswith('/'):
            # SSH远程文件 (host:path)
            host, path = filepath.split(':', 1)
            result = subprocess.run(
                ['ssh', host, f'cat {path}'],
                capture_output=True,
                text=True
            )
            if result.returncode != 0:
                raise Exception(f"无法读取远程文件: {result.stderr}")
            return result.stdout
        
        # 本地文件
        with open(filepath, 'r', encoding='utf-8') as f:
            return f.read()
    
    def normalize_text(self, content: str, file_type: str) -> List[str]:
        """标准化文本（去除注释、空行等）"""
        lines = content.split('\n')
        result = []
        
        for line in lines:
            # 去除空行
            if self.ignore_blank and not line.strip():
                continue
            
            # 去除注释
            if self.ignore_comments:
                if file_type in ['yaml', 'ini', 'properties', 'env']:
                    # #注释
                    if line.strip().startswith('#'):
                        continue
                    # 行尾注释
                    if '#' in line:
                        line = line.split('#')[0].rstrip()
                
                elif file_type == 'json':
                    # JSON标准不支持注释，但有些工具支持
                    if line.strip().startswith('//'):
                        continue
            
            result.append(line)
        
        return result
    
    def parse_yaml(self, content: str) -> dict:
        """解析YAML并排序（可选）"""
        data = yaml.safe_load(content)
        if self.ignore_order and isinstance(data, dict):
            return self._sort_dict_recursive(data)
        return data
    
    def parse_json(self, content: str) -> dict:
        """解析JSON并排序（可选）"""
        data = json.loads(content)
        if self.ignore_order and isinstance(data, dict):
            return self._sort_dict_recursive(data)
        return data
    
    def _sort_dict_recursive(self, obj):
        """递归排序字典"""
        if isinstance(obj, dict):
            return {k: self._sort_dict_recursive(v) for k, v in sorted(obj.items())}
        elif isinstance(obj, list):
            return [self._sort_dict_recursive(item) for item in obj]
        return obj
    
    def diff_structured(self, file1: str, file2: str, file_type: str) -> str:
        """结构化对比（YAML/JSON）"""
        content1 = self.read_file(file1)
        content2 = self.read_file(file2)
        
        try:
            if file_type == 'yaml':
                data1 = self.parse_yaml(content1)
                data2 = self.parse_yaml(content2)
                # 转为格式化YAML字符串
                str1 = yaml.dump(data1, default_flow_style=False, sort_keys=True)
                str2 = yaml.dump(data2, default_flow_style=False, sort_keys=True)
            else:  # json
                data1 = self.parse_json(content1)
                data2 = self.parse_json(content2)
                # 转为格式化JSON字符串
                str1 = json.dumps(data1, indent=2, sort_keys=True)
                str2 = json.dumps(data2, indent=2, sort_keys=True)
            
            return self._unified_diff(str1, str2, file1, file2)
        
        except Exception as e:
            print(f"{Fore.YELLOW}警告: 结构化解析失败，回退到文本对比: {e}")
            return self.diff_text(file1, file2, file_type)
    
    def diff_text(self, file1: str, file2: str, file_type: str) -> str:
        """文本对比"""
        content1 = self.read_file(file1)
        content2 = self.read_file(file2)
        
        lines1 = self.normalize_text(content1, file_type)
        lines2 = self.normalize_text(content2, file_type)
        
        return self._unified_diff('\n'.join(lines1), '\n'.join(lines2), file1, file2)
    
    def _unified_diff(self, str1: str, str2: str, file1: str, file2: str) -> str:
        """生成unified diff"""
        lines1 = str1.splitlines(keepends=True)
        lines2 = str2.splitlines(keepends=True)
        
        diff = unified_diff(
            lines1,
            lines2,
            fromfile=file1,
            tofile=file2,
            n=self.context_lines
        )
        
        return self._colorize_diff(diff)
    
    def _colorize_diff(self, diff) -> str:
        """彩色化diff输出"""
        result = []
        for line in diff:
            if line.startswith('---') or line.startswith('+++'):
                result.append(f"{Fore.CYAN}{Style.BRIGHT}{line}{Style.RESET_ALL}")
            elif line.startswith('@@'):
                result.append(f"{Fore.MAGENTA}{line}{Style.RESET_ALL}")
            elif line.startswith('+'):
                result.append(f"{Fore.GREEN}{line}{Style.RESET_ALL}")
            elif line.startswith('-'):
                result.append(f"{Fore.RED}{line}{Style.RESET_ALL}")
            else:
                result.append(line)
        
        return ''.join(result)
    
    def compare(self, file1: str, file2: str) -> Tuple[bool, str]:
        """对比两个文件"""
        file_type = self.detect_type(file1)
        
        print(f"{Fore.CYAN}文件类型: {file_type}")
        print(f"{Fore.CYAN}文件1: {file1}")
        print(f"{Fore.CYAN}文件2: {file2}")
        print(f"{Fore.CYAN}{'='*60}{Style.RESET_ALL}\n")
        
        if file_type in ['yaml', 'json']:
            diff_output = self.diff_structured(file1, file2, file_type)
        else:
            diff_output = self.diff_text(file1, file2, file_type)
        
        if not diff_output.strip():
            return True, f"{Fore.GREEN}✓ 文件完全相同{Style.RESET_ALL}"
        else:
            return False, diff_output
    
    def compare_dirs(self, dir1: str, dir2: str, pattern: str = "*") -> Dict[str, str]:
        """对比两个目录"""
        from pathlib import Path
        import fnmatch
        
        results = {}
        path1 = Path(dir1)
        path2 = Path(dir2)
        
        # 获取所有匹配文件
        files1 = set()
        for f in path1.rglob(pattern):
            if f.is_file():
                files1.add(f.relative_to(path1))
        
        files2 = set()
        for f in path2.rglob(pattern):
            if f.is_file():
                files2.add(f.relative_to(path2))
        
        # 只在dir1中
        only_in_1 = files1 - files2
        if only_in_1:
            print(f"\n{Fore.YELLOW}只在 {dir1} 中:{Style.RESET_ALL}")
            for f in sorted(only_in_1):
                print(f"  {f}")
        
        # 只在dir2中
        only_in_2 = files2 - files1
        if only_in_2:
            print(f"\n{Fore.YELLOW}只在 {dir2} 中:{Style.RESET_ALL}")
            for f in sorted(only_in_2):
                print(f"  {f}")
        
        # 两边都有的文件
        common = files1 & files2
        
        print(f"\n{Fore.CYAN}对比共有文件 ({len(common)} 个):{Style.RESET_ALL}\n")
        
        for rel_path in sorted(common):
            f1 = path1 / rel_path
            f2 = path2 / rel_path
            
            print(f"\n{Fore.CYAN}{'='*60}")
            print(f"对比: {rel_path}")
            print(f"{'='*60}{Style.RESET_ALL}")
            
            same, diff = self.compare(str(f1), str(f2))
            results[str(rel_path)] = 'same' if same else 'different'
            
            if same:
                print(diff)
            else:
                print(diff)
        
        return results


def main():
    parser = argparse.ArgumentParser(
        description='智能配置文件对比工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:

  # 对比两个本地文件
  %(prog)s prod.yaml test.yaml

  # 对比本地和远程
  %(prog)s local.conf server1:/etc/nginx/nginx.conf

  # 对比两个目录
  %(prog)s -d /etc/nginx/prod /etc/nginx/test

  # YAML对比（忽略顺序）
  %(prog)s --ignore-order prod.yaml test.yaml

  # 包含注释和空行
  %(prog)s --no-ignore-comments prod.conf test.conf

  # 指定文件模式
  %(prog)s -d /etc/prod /etc/test -p "*.conf"
        """
    )
    
    parser.add_argument('file1', nargs='?', help='第一个文件或目录')
    parser.add_argument('file2', nargs='?', help='第二个文件或目录')
    
    parser.add_argument('-d', '--directory', action='store_true',
                       help='对比目录而非文件')
    parser.add_argument('-p', '--pattern', default='*',
                       help='目录对比时的文件模式 (默认: *)')
    
    parser.add_argument('--ignore-order', action='store_true',
                       help='忽略YAML/JSON的键顺序')
    parser.add_argument('--no-ignore-comments', action='store_true',
                       help='不忽略注释')
    parser.add_argument('--no-ignore-blank', action='store_true',
                       help='不忽略空行')
    
    parser.add_argument('-c', '--context', type=int, default=3,
                       help='上下文行数 (默认: 3)')
    
    parser.add_argument('--html', metavar='OUTPUT',
                       help='输出HTML格式报告')
    
    args = parser.parse_args()
    
    if not args.file1 or not args.file2:
        parser.print_help()
        sys.exit(1)
    
    differ = ConfigDiffer(
        ignore_comments=not args.no_ignore_comments,
        ignore_blank=not args.no_ignore_blank,
        ignore_order=args.ignore_order,
        context_lines=args.context
    )
    
    try:
        if args.directory:
            results = differ.compare_dirs(args.file1, args.file2, args.pattern)
            
            # 统计
            same = sum(1 for v in results.values() if v == 'same')
            diff = sum(1 for v in results.values() if v == 'different')
            
            print(f"\n{Fore.CYAN}{'='*60}")
            print(f"汇总")
            print(f"{'='*60}{Style.RESET_ALL}")
            print(f"{Fore.GREEN}相同: {same}")
            print(f"{Fore.RED}差异: {diff}")
            print(f"总计: {same + diff}")
        else:
            same, diff = differ.compare(args.file1, args.file2)
            print(diff)
            sys.exit(0 if same else 1)
    
    except Exception as e:
        print(f"{Fore.RED}错误: {e}{Style.RESET_ALL}", file=sys.stderr)
        sys.exit(2)


if __name__ == '__main__':
    main()
