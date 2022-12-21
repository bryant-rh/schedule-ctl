package cmd

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/bryant-rh/schedule-ctl/pkg/save"
	"github.com/bryant-rh/schedule-ctl/pkg/schedule"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	version  string
	members  string
	fileName string

	outputFileName string
	outputDefault  = fmt.Sprintf("%s.xlsx", time.Now().Format("2006-01-02"))
)

var (
	memberList []string
)

// versionString returns the version prefixed by 'v'
// or an empty string if no version has been populated by goreleaser.
// In this case, the --version flag will not be added by cobra.
func versionString() string {
	if len(version) == 0 {
		return ""
	}
	return "v" + version
}

func NewCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:           "schedulectl",
		Short:         "schedulectl is a tool for scheduling",
		Version:       versionString(),
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if members == "" && fileName == "" {
				cmd.Help()
				return

			}
			//校验参数
			if members != "" && fileName != "" {
				pterm.Error.Println("cannot specify --member and --filename at the same time")
				return
			}
			if fileName != "" {
				// 读取文件解析
				bytes, err := ioutil.ReadFile(fileName)
				if err != nil {
					pterm.Error.Printfln("文件: %s 读取失败.", fileName)
					return
				}
				str := string(bytes)
				str = strings.Replace(str, "\r", "", -1)
				str = strings.Trim(str, " ")
				memberList = strings.Split(str, "\n")
				if len(memberList) == 0 {
					pterm.Error.Println("文件内容为空")
					return
				}
			}

			if members != "" {
				memberList = strings.Split(members, ",")

			}
			// 排除掉不符合条件的成员
			var m []string
			for _, v := range memberList {
				if len(v) > 0 {
					m = append(m, v)
				}
			}
			var totalM int = len(m)
			if totalM < 2 {
				notice := fmt.Sprintf("成员数量是:%s, 不符合要求", strconv.Itoa(totalM))
				pterm.Info.Println(notice)
				return
			}
			//fmt.Println(memberList)

			// 用户确认是否正确
			confirmMemberMsg := fmt.Sprintf("当前一共有%s个成员,列表如下:", strconv.Itoa(totalM))
			pterm.NewRGB(15, 199, 209).Println(confirmMemberMsg)
			i := 1
			for _, v := range m {
				pterm.NewRGB(178, 44, 199).Println(strconv.Itoa(i) + ": " + v)
				i++
			}

			confirmMember, _ := pterm.DefaultInteractiveConfirm.Show()
			pterm.Println() // Blank line

			if !confirmMember {
				pterm.Info.Printfln(pterm.Red("Bye-bye"))
				return
			}

			// 用户选择每天值班人数
			numsOneDay := setNumsOneDay(totalM)

			// 执班总天数
			totalDay := setTotalDay(numsOneDay, totalM)

			// 用户确认
			numsOneDayNotice := fmt.Sprintf("你输入的每天值班人数是:%s, 总值班天数是:%s", strconv.Itoa(numsOneDay), strconv.Itoa(totalDay))
			fmt.Println(pterm.Green(numsOneDayNotice))
			pterm.Println()
			fmt.Println(pterm.Yellow("正在生成排班计划,请稍等..."))

			schedule := schedule.Schedule{}
			result := schedule.Create(m, numsOneDay, totalDay)

			for _, v := range result {
				fmt.Println(v)
			}

			err := save.SaveExcel(result, outputFileName)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("数据成功保存到:" + outputFileName + ", 再会")
			}

			//endNotice()

		},
	}

	rootCmd.PersistentFlags().StringVarP(&members, "member", "m", "", "指定成员,多个用逗号分割")
	rootCmd.PersistentFlags().StringVarP(&fileName, "filename", "f", "", "指定成员文件")
	rootCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", outputDefault, "指定输出文件名")
	return rootCmd
}

func setNumsOneDay(max int) int {
	confirm := false
	var result int
	var num int

	for !confirm {
		// Text input with multi line enabled
		data, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("请输入每天值班人数,需大于0小于成员数")
		// Blank line
		pterm.Println()
		nums, err := strconv.ParseInt(data, 10, 64)
		num = int(nums)
		if err != nil || num < 1 || num >= max {
			continue
		} else {
			confirm, result = true, num
			break
		}

	}

	return result

}

//func totalDay
func setTotalDay(numsOneDay, totalM int) int {
	confirm := false

	var result int
	var total int

	for !confirm {
		// Text input with multi line enabled
		data, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("请输入值班总天数")
		// Blank line
		pterm.Println()

		totalT, err := strconv.ParseInt(data, 10, 64)
		total = int(totalT)
		ava := (total * numsOneDay) / totalM //平均每人要值班几天
		if err == nil && ava > 0 {
			result = (ava * totalM) / numsOneDay
			confirm = true
			break
		} else {
			continue
		}
	}

	if total != result {
		noticeMsg := fmt.Sprintf("由于你输入的天数%s不能保证每人值班天数相同,已自动设为%s天", strconv.Itoa(total), strconv.Itoa(result))
		fmt.Println(pterm.Red(noticeMsg))
	}

	return result
}
