package controllers

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/astaxie/beego"
)

// Operations about Users
type HookController struct {
	beego.Controller
}

func execGitShell(shell string) {
	file := beego.AppConfig.String("gitShellDir") + shell + ".sh"
	cmd := exec.Command("bash", "-c", file)

	//显示运行的命令
	beego.Info(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		beego.Error(err)
		return
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		beego.Info(line)
	}
	cmd.Wait()
}

func (this *HookController) Hook() {

	body := this.Ctx.Input.RequestBody

	//header := this.Ctx.Input.
	signature := this.Ctx.Input.Header("X-Hub-Signature")
	beego.Info("X-Hub-Signature:" + signature)
	beego.Info("body:" + string(body))

	if signature == "" {
		this.Data["json"] = map[string]interface{}{"sMsg": "No signature", "iRet": -1001}
		this.ServeJSON()
		return
	}

	arr := strings.Split(signature, "=")
	if len(arr) != 2 {
		this.Data["json"] = map[string]interface{}{"sMsg": "Signature is not right", "iRet": -1002}
		this.ServeJSON()
		return
	}

	//hmac ,use sha1
	key := []byte(beego.AppConfig.String("MySiteSecret"))
	mac := hmac.New(sha1.New, key)
	mac.Write(body)
	//fmt.Printf("%x\n", mac.Sum(nil))
	hmacStr := fmt.Sprintf("%x", mac.Sum(nil))
	beego.Info("hmacStr = " + hmacStr)
	if hmacStr != arr[1] {
		this.Data["json"] = map[string]interface{}{"sMsg": "Check signature faild", "iRet": -1003}
		this.ServeJSON()
		return
	}
	beego.Info("Git pull start...")
	execGitShell("GitMySite")
	beego.Info("Git pull end.")

	this.Data["json"] = map[string]interface{}{"sMsg": "OK", "iRet": 0}
	this.ServeJSON()
}
