# schedule-ctl
一个用于自动生成排班表的工具。

# Usage

```Bash
schedulectl is a tool for scheduling

Usage:
  schedulectl [flags]

Flags:
  -f, --filename string   通过文件指定值班人员
  -h, --help              help for schedulectl
  -m, --member string     直接通过命令行指定值班人员,多个用逗号分割
  -o, --output string     指定输出文件名 (default "2022-12-21.xlsx")
```
# Demo
## 1.直接通过命令行指定值班人员,多个用逗号分割
![demo1](image/demo-member.gif)

## 2.通过文件指定值班人员
![demo2](image/demo-file.gif)
