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

package conf

// GitBackup 将工作空间 data/ 目录单向备份（提交并推送）到用户自建的 Git 远程仓库
type GitBackup struct {
	Enabled     bool   `json:"enabled"`     // 是否启用 Git 备份
	RepoURL     string `json:"repoURL"`     // 远程仓库 HTTPS 地址，如 https://github.com/user/notes.git
	Branch      string `json:"branch"`      // 推送分支，默认 main
	Token       string `json:"token"`       // 访问令牌（个人访问令牌 PAT），用于 HTTPS 鉴权
	Username    string `json:"username"`    // 提交作者名及 HTTPS 用户名（可留空，令牌鉴权时通常不校验用户名）
	Email       string `json:"email"`       // 提交作者邮箱
	AutoEnabled bool   `json:"autoEnabled"` // 是否按间隔自动推送
	AutoMinutes int    `json:"autoMinutes"` // 自动推送间隔（分钟），最小 5
}

func NewGitBackup() *GitBackup {
	return &GitBackup{
		Enabled:     false,
		Branch:      "main",
		AutoEnabled: false,
		AutoMinutes: 30,
	}
}
