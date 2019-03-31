package controllers

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/astaxie/beego"
)

type repository struct {
	Name string `json:"name"`
}

type HookReq struct {
	Repository repository `json:"repository"`
}

// Operations about Users
type HookController struct {
	beego.Controller
}

func execGitShell(shell string) {
	beego.Info("Do exec shell start...")
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
	beego.Info("Do exec shell end.")
}

func (this *HookController) Hook() {

	body := this.Ctx.Input.RequestBody

	var req HookReq
	if ok := json.Unmarshal(body, &req); ok != nil {
		beego.Info("Hook req to json failed.")
	}

	if req.Repository.Name == "" {
		this.Data["json"] = map[string]interface{}{"sMsg": "Get Repository failde", "iRet": -1000}
		this.ServeJSON()
		return
	}
	beego.Info("Repository = " + req.Repository.Name)

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
	key := []byte(beego.AppConfig.String(req.Repository.Name + "Secret"))
	mac := hmac.New(sha1.New, key)
	mac.Write(body)
	hmacStr := fmt.Sprintf("%x", mac.Sum(nil))
	beego.Info("hmacStr = " + hmacStr)
	if hmacStr != arr[1] {
		this.Data["json"] = map[string]interface{}{"sMsg": "Check signature faild", "iRet": -1003}
		this.ServeJSON()
		return
	}

	beego.Info("Git pull start...")
	go execGitShell(req.Repository.Name + "Git")
	beego.Info("Git pull end.")

	this.Data["json"] = map[string]interface{}{"sMsg": "OK", "iRet": 0}
	this.ServeJSON()
}
