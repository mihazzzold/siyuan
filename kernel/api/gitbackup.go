// SiYuan - Refactor your thinking
// Copyright (c) 2020-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"net/http"

	"github.com/88250/gulu"
	"github.com/gin-gonic/gin"
	"github.com/siyuan-note/siyuan/kernel/conf"
	"github.com/siyuan-note/siyuan/kernel/model"
	"github.com/siyuan-note/siyuan/kernel/util"
)

func getGitBackup(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	gitBackup := model.Conf.GitBackup
	if nil == gitBackup {
		gitBackup = conf.NewGitBackup()
	}

	// 不回传令牌明文，仅告知是否已设置
	masked := *gitBackup
	tokenSet := "" != masked.Token
	masked.Token = ""

	ret.Data = map[string]any{
		"gitBackup": masked,
		"tokenSet":  tokenSet,
	}
}

func setGitBackup(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}

	param, err := gulu.JSON.MarshalJSON(arg)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	gitBackup := conf.NewGitBackup()
	if err = gulu.JSON.UnmarshalJSON(param, gitBackup); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	// 令牌为空表示不修改，沿用已保存的令牌
	if "" == gitBackup.Token && nil != model.Conf.GitBackup {
		gitBackup.Token = model.Conf.GitBackup.Token
	}

	model.Conf.GitBackup = gitBackup
	model.Conf.Save()

	masked := *gitBackup
	masked.Token = ""
	ret.Data = map[string]any{
		"gitBackup": masked,
		"tokenSet":  "" != gitBackup.Token,
	}
}

func performGitBackup(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	if err := model.PerformGitBackup(); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		ret.Data = map[string]any{"closeTimeout": 7000}
		return
	}

	ret.Msg = "Отправлено в Git"
}
