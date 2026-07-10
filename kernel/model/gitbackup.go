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

package model

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHTTP "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/siyuan-note/logging"
	"github.com/siyuan-note/siyuan/kernel/util"
)

// gitBackupLock 保证同一时刻只有一次备份在执行
var gitBackupLock = sync.Mutex{}

// PerformGitBackup 将工作空间 data/ 目录提交并推送到配置的 Git 远程仓库（单向备份）。
// 首次执行时会在 data/ 下初始化仓库并绑定远程；随后每次提交全部变更并强制以本地为准推送。
func PerformGitBackup() (err error) {
	conf := Conf.GitBackup
	if nil == conf || !conf.Enabled {
		return errors.New("Git-бэкап не включён")
	}
	if "" == strings.TrimSpace(conf.RepoURL) {
		return errors.New("Не указан адрес репозитория")
	}

	if !gitBackupLock.TryLock() {
		return errors.New("Резервное копирование уже выполняется")
	}
	defer gitBackupLock.Unlock()

	branch := strings.TrimSpace(conf.Branch)
	if "" == branch {
		branch = "main"
	}

	// 备份 data/ 目录（笔记、资源等），先落盘再提交
	FlushTxQueue()
	dataDir := util.DataDir

	repo, err := git.PlainOpen(dataDir)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		repo, err = git.PlainInit(dataDir, false)
		if err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("git open failed: %w", err)
	}

	// 绑定/更新远程地址
	if err = ensureGitRemote(repo, conf.RepoURL); err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("git worktree failed: %w", err)
	}

	// 暂存全部变更（含新增、修改、删除）
	if err = worktree.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("git status failed: %w", err)
	}

	authorName := strings.TrimSpace(conf.Username)
	if "" == authorName {
		authorName = "SiYuan"
	}
	authorEmail := strings.TrimSpace(conf.Email)
	if "" == authorEmail {
		authorEmail = "siyuan@localhost"
	}
	signature := &object.Signature{Name: authorName, Email: authorEmail, When: time.Now()}

	// 仅在有变更时提交，避免空提交
	if !status.IsClean() {
		msg := "SiYuan backup " + time.Now().Format("2006-01-02 15:04:05")
		if _, err = worktree.Commit(msg, &git.CommitOptions{Author: signature, Committer: signature}); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}
	}

	// 将当前 HEAD 指向目标分支引用，保证推送到期望分支
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("git head failed: %w", err)
	}
	localRef := plumbing.NewBranchReferenceName(branch)
	if head.Name() != localRef {
		if err = repo.Storer.SetReference(plumbing.NewHashReference(localRef, head.Hash())); err != nil {
			return fmt.Errorf("git set branch failed: %w", err)
		}
	}

	auth := &gitHTTP.BasicAuth{Username: gitAuthUsername(conf.Username), Password: conf.Token}
	refSpec := config.RefSpec("+" + localRef.String() + ":" + localRef.String())
	pushErr := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		RefSpecs:   []config.RefSpec{refSpec},
		Force:      true, // 单向备份：始终以本地为准
	})
	if pushErr != nil && !errors.Is(pushErr, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("git push failed: %w", pushErr)
	}

	logging.LogInfof("git backup pushed to [%s] branch [%s]", conf.RepoURL, branch)
	return nil
}

// ensureGitRemote 确保仓库存在名为 origin 且地址正确的远程
func ensureGitRemote(repo *git.Repository, url string) error {
	remote, err := repo.Remote("origin")
	if errors.Is(err, git.ErrRemoteNotFound) {
		_, err = repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{url}})
		return err
	}
	if err != nil {
		return err
	}
	// 地址变化时重建远程
	if 1 > len(remote.Config().URLs) || remote.Config().URLs[0] != url {
		if err = repo.DeleteRemote("origin"); err != nil {
			return err
		}
		_, err = repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{url}})
		return err
	}
	return nil
}

// gitAuthUsername 令牌鉴权时用户名可为任意非空值，缺省用 x-access-token（GitHub/GitLab 均接受）
func gitAuthUsername(username string) string {
	if "" != strings.TrimSpace(username) {
		return username
	}
	return "x-access-token"
}

var lastGitBackupMinute = -1

// GitBackupJob 由 cron 每分钟调用，按配置的间隔自动推送备份
func GitBackupJob() {
	conf := Conf.GitBackup
	if nil == conf || !conf.Enabled || !conf.AutoEnabled {
		return
	}
	interval := conf.AutoMinutes
	if 5 > interval {
		interval = 5
	}
	// 以“分钟数整除间隔”为触发点，并去重同一分钟内的多次触发
	minutes := time.Now().Hour()*60 + time.Now().Minute()
	if 0 != minutes%interval || minutes == lastGitBackupMinute {
		return
	}
	lastGitBackupMinute = minutes
	if err := PerformGitBackup(); err != nil {
		logging.LogWarnf("auto git backup failed: %s", err)
	}
}
