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
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/88250/gulu"
	"github.com/88250/lute/ast"
	"github.com/gin-gonic/gin"
	"github.com/siyuan-note/filelock"
	"github.com/siyuan-note/logging"
	"github.com/siyuan-note/siyuan/kernel/av"
	"github.com/siyuan-note/siyuan/kernel/conf"
	"github.com/siyuan-note/siyuan/kernel/sql"
	"github.com/siyuan-note/siyuan/kernel/treenode"
	"github.com/siyuan-note/siyuan/kernel/util"
)

type PublishAccessItem struct {
	ID       string `json:"id"`
	Visible  bool   `json:"visible"`  // 是否发布可见
	Password string `json:"password"` // 密码，为空字符串时表示无密码
	Disable  bool   `json:"disable"`  // 是否禁止发布
	SelfOnly bool   `json:"selfOnly"` // 规则仅作用于本文档，不继承给子文档
}

type PublishAccess []*PublishAccessItem

var (
	publishAccessLastModified int64
	publishAccess             PublishAccess
	publishAccessLock         = sync.Mutex{}
)

func GetPublishAccess() (ret PublishAccess) {
	ret = PublishAccess{}
	now := time.Now().UnixMilli()
	if now-publishAccessLastModified < 30*1000 {
		return publishAccess
	}

	publishAccessLock.Lock()
	defer publishAccessLock.Unlock()

	publishAccessLastModified = now

	publishAccessPath := filepath.Join(util.DataDir, ".siyuan", "publishAccess.json")
	err := os.MkdirAll(filepath.Dir(publishAccessPath), 0755)
	if err != nil {
		return
	}
	if !filelock.IsExist(publishAccessPath) {
		if err = filelock.WriteFile(publishAccessPath, []byte("[]")); err != nil {
			logging.LogErrorf("create publishAccess.json [%s] failed: %s", publishAccessPath, err)
			return
		}
	}
	data, err := os.ReadFile(publishAccessPath)
	if err != nil {
		logging.LogErrorf("read publishAccess.json [%s] failed: %s", publishAccessPath, err)
		return
	}
	if err = gulu.JSON.UnmarshalJSON(data, &publishAccess); err != nil {
		logging.LogWarnf("unmarshal publishAccess.json failed: %s", err)
		return
	}
	ret = publishAccess
	return
}

func SetPublishAccess(inputPublishAccess PublishAccess) (err error) {
	now := time.Now().UnixMilli()
	publishAccessLock.Lock()
	defer publishAccessLock.Unlock()
	publishAccessLastModified = now
	publishAccess = inputPublishAccess

	publishAccessPath := filepath.Join(util.DataDir, ".siyuan", "publishAccess.json")
	err = os.MkdirAll(filepath.Dir(publishAccessPath), 0755)
	if err != nil {
		msg := fmt.Sprintf("create dir for publishAccess.json [%s] failed: %s", publishAccessPath, err)
		logging.LogErrorf(msg)
		err = errors.New(msg)
		return
	}

	data, err := gulu.JSON.MarshalJSON(inputPublishAccess)
	if err != nil {
		logging.LogErrorf("marshal publishAccess.json [%s] failed: %s", publishAccessPath, err)
		return
	}

	err = filelock.WriteFile(publishAccessPath, data)
	if err != nil {
		msg := fmt.Sprintf("write publishAccess.json [%s] failed: %s", publishAccessPath, err)
		logging.LogErrorf(msg)
		err = errors.New(msg)
		return
	}
	return
}

func PurgePublishAccess() {
	publishAccess := GetPublishAccess()
	IDs := []string{}
	for _, item := range publishAccess {
		IDs = append(IDs, item.ID)
	}

	boxes, err := ListNotebooks()
	if err != nil {
		return
	}
	// 必须在所有笔记本都打开的情况下才能执行清除工作，否则会把关闭的笔记本里文档的发布访问控制状态清除
	for _, box := range boxes {
		if box.Closed {
			return
		}
	}

	checkResult := treenode.ExistBlockTrees(IDs)
	tempPublishAccess := PublishAccess{}
	for i, ID := range IDs {
		if exists, ok := checkResult[ID]; ok && exists {
			tempPublishAccess = append(tempPublishAccess, publishAccess[i])
		} else {
			for _, box := range boxes {
				if box.ID == ID {
					tempPublishAccess = append(tempPublishAccess, publishAccess[i])
					break
				}
			}
		}
	}
	SetPublishAccess(tempPublishAccess)
	return
}

// PublishVerdict 某文档路径的发布访问裁决结果
type PublishVerdict struct {
	Visible    bool   // 是否在目录树、搜索等列表中可见
	Disable    bool   // 是否禁止访问内容
	PasswordID string // 提供密码规则的文档 ID
	Password   string // 访问密码，为空表示无需密码
}

// ResolvePublishAccess 从文档自身向祖先方向查找最近的显式规则并以其为准；
// SelfOnly 规则只作用于规则所在文档本身；没有任何规则时按 Publish.DefaultDeny 决定默认裁决
func ResolvePublishAccess(box string, blockPath string, publishAccess PublishAccess) (verdict PublishVerdict) {
	byID := make(map[string]*PublishAccessItem, len(publishAccess))
	for _, item := range publishAccess {
		byID[item.ID] = item
	}

	self := true
	currentPath := blockPath
	for currentPath != "/" && currentPath != "." && currentPath != "" {
		currentID := strings.TrimSuffix(path.Base(currentPath), ".sy")
		if item, ok := byID[currentID]; ok && (self || !item.SelfOnly) {
			return PublishVerdict{
				Visible:    item.Visible && !item.Disable,
				Disable:    item.Disable,
				PasswordID: item.ID,
				Password:   item.Password,
			}
		}
		self = false
		currentPath = path.Dir(currentPath)
	}
	if item, ok := byID[box]; ok && !item.SelfOnly {
		return PublishVerdict{
			Visible:    item.Visible && !item.Disable,
			Disable:    item.Disable,
			PasswordID: item.ID,
			Password:   item.Password,
		}
	}

	if Conf.Publish.DefaultDeny {
		return PublishVerdict{Visible: false, Disable: true}
	}
	return PublishVerdict{Visible: true, Disable: false}
}

// CheckPathEnabledByPublishAccess 内容是否允许访问（禁止访问轴）
func CheckPathEnabledByPublishAccess(box string, blockPath string, publishAccess PublishAccess) bool {
	return !ResolvePublishAccess(box, blockPath, publishAccess).Disable
}

// CheckPathVisibleByPublishAccess 是否在列表（目录树、搜索、图谱等）中可见
func CheckPathVisibleByPublishAccess(box string, blockPath string, publishAccess PublishAccess) bool {
	verdict := ResolvePublishAccess(box, blockPath, publishAccess)
	return verdict.Visible && !verdict.Disable
}

func GetPathPasswordByPublishAccess(box string, blockPath string, publishAccess PublishAccess) (passwordID string, password string) {
	verdict := ResolvePublishAccess(box, blockPath, publishAccess)
	return verdict.PasswordID, verdict.Password
}

func CheckBlockIdAccessableByPublishAccess(c *gin.Context, publishAccess PublishAccess, blockID string) bool {
	bt := treenode.GetBlockTree(blockID)
	if bt == nil {
		return false
	}
	verdict := ResolvePublishAccess(bt.BoxID, bt.Path, publishAccess)
	return !verdict.Disable && (verdict.Password == "" || CheckPublishAuthCookie(c, verdict.PasswordID, verdict.Password))
}

func SetPublishAuthCookie(c *gin.Context, ID string, password string) {
	authCookie := util.SHA256Hash([]byte(ID + password))
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "publish-auth-" + ID,
		Value:    authCookie,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Secure:   util.SSL,
		HttpOnly: true,
	})
}

func CheckPublishAuthCookie(c *gin.Context, ID string, password string) bool {
	authCookie, err := c.Request.Cookie("publish-auth-" + ID)
	return err == nil && authCookie.Value == util.SHA256Hash([]byte(ID+password))
}

func CheckAbsPathAccessableByPublishAccess(c *gin.Context, absPath string, publishAccess PublishAccess) bool {
	absPath = filepath.Clean(absPath)

	if gulu.File.IsSubPath(util.HistoryDir, absPath) {
		return false
	}

	if gulu.File.IsSubPath(util.DataDir, absPath) {
		relPath, err := filepath.Rel(util.DataDir, absPath)
		if err != nil {
			return true
		}

		relPath = strings.ReplaceAll(relPath, "\\", "/")

		pathParts := strings.Split(relPath, "/")
		if len(pathParts) <= 1 {
			return true
		}

		if ast.IsNodeIDPattern(pathParts[0]) {
			box := pathParts[0]
			blockPath := "/" + strings.Join(pathParts[1:], "/")
			passwordID, password := GetPathPasswordByPublishAccess(box, blockPath, publishAccess)
			return CheckPathEnabledByPublishAccess(box, blockPath, publishAccess) && (password == "" || CheckPublishAuthCookie(c, passwordID, password))
		} else if pathParts[0] == "assets" {
			bts := treenode.GetBlockTreesByType("d")
			for _, bt := range bts {
				passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
				if CheckPathEnabledByPublishAccess(bt.BoxID, bt.Path, publishAccess) && (password == "" || CheckPublishAuthCookie(c, passwordID, password)) {
					assets, _ := DocAssets(bt.ID, false)
					for _, assetPath := range assets {
						if assetPath == relPath {
							return true
						}
					}
				}
			}
			return false
		}
	}
	return false
}

func FilterViewByPublishAccess(c *gin.Context, publishAccess PublishAccess, viewable av.Viewable) (ret av.Viewable) {
	ret = viewable

	switch ret.GetType() {
	case av.LayoutTypeTable:
		table := ret.(*av.Table)
		filteredRows := []*av.TableRow{}
		for _, row := range table.Rows {
			// 默认第一个属性是文档块
			var bt *treenode.BlockTree
			if len(row.Cells) > 0 {
				if row.Cells[0].Value.Block != nil {
					id := row.Cells[0].Value.Block.ID
					if id != "" {
						bt = treenode.GetBlockTree(id)
					}
				}
			}
			if bt != nil {
				// 不显示禁止文档
				if !CheckPathEnabledByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
					row = nil
				}
			}
			if row != nil {
				filteredRows = append(filteredRows, row)
			}
		}
		table.Rows = filteredRows
		if table.Groups != nil {
			for i, viewable := range table.Groups {
				table.Groups[i] = FilterViewByPublishAccess(c, publishAccess, viewable)
			}
		}
	case av.LayoutTypeGallery:
		gallery := ret.(*av.Gallery)
		filteredCards := []*av.GalleryCard{}
		for _, card := range gallery.Cards {
			// 默认第一个属性是文档块
			var bt *treenode.BlockTree
			if len(card.Values) > 0 {
				if card.Values[0].Value.Block != nil {
					id := card.Values[0].Value.Block.ID
					if id != "" {
						bt = treenode.GetBlockTree(id)
					}
				}
			}
			if bt != nil {
				// 替换封面
				newCoverContent := FilterContentByPublishAccess(c, publishAccess, bt.BoxID, bt.Path, card.CoverContent, true)
				if card.CoverContent != newCoverContent {
					card.CoverContent = newCoverContent
					card.CoverURL = ""
				}

				// 不显示禁止文档
				if !CheckPathEnabledByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
					card = nil
				}
			}
			if card != nil {
				filteredCards = append(filteredCards, card)
			}
		}
		gallery.Cards = filteredCards
		if gallery.Groups != nil {
			for i, viewable := range gallery.Groups {
				gallery.Groups[i] = FilterViewByPublishAccess(c, publishAccess, viewable)
			}
		}
	case av.LayoutTypeKanban:
		kanban := ret.(*av.Kanban)
		filteredCards := []*av.KanbanCard{}
		for _, card := range kanban.Cards {
			// 默认第一个属性是文档块
			var bt *treenode.BlockTree
			if len(card.Values) > 0 {
				if card.Values[0].Value.Block != nil {
					id := card.Values[0].Value.Block.ID
					if id != "" {
						bt = treenode.GetBlockTree(id)
					}
				}
			}
			if bt != nil {
				// 替换封面
				newCoverContent := FilterContentByPublishAccess(c, publishAccess, bt.BoxID, bt.Path, card.CoverContent, true)
				if card.CoverContent != newCoverContent {
					card.CoverContent = newCoverContent
					card.CoverURL = ""
				}

				// 不显示禁止文档
				if !CheckPathEnabledByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
					card = nil
				}
			}
			if card != nil {
				filteredCards = append(filteredCards, card)
			}
		}
		kanban.Cards = filteredCards
		kanban.CardCount = len(kanban.Cards)
		if kanban.Groups != nil {
			for i, viewable := range kanban.Groups {
				kanban.Groups[i] = FilterViewByPublishAccess(c, publishAccess, viewable)
			}
		}
	}
	return
}

func FilterBlockAttributeViewKeysByPublishAccess(c *gin.Context, publishAccess PublishAccess, blockAttributeViewKeys []*BlockAttributeViewKeys) (ret []*BlockAttributeViewKeys) {
	ret = []*BlockAttributeViewKeys{}
	for _, blockAttributeViewKey := range blockAttributeViewKeys {
		accessable := false
		bts := treenode.GetBlockTrees(blockAttributeViewKey.BlockIDs)
		for _, bt := range bts {
			passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
			if (password == "" || CheckPublishAuthCookie(c, passwordID, password)) && CheckPathEnabledByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
				accessable = true
				break
			}
		}
		if accessable {
			ret = append(ret, blockAttributeViewKey)
		}
	}
	return
}

func FilterBlockInfoByPublishAccess(c *gin.Context, publishAccess PublishAccess, info *BlockInfo) (ret *BlockInfo) {
	ret = info
	if info == nil {
		return
	}

	filteredAttrViews := []*AttrView{}
	avIDs := []string{}
	for _, attrView := range info.AttrViews {
		avBlocksAccessable := false
		if attrView.ID != "" {
			avBlockIDs := treenode.GetMirrorAttrViewBlockIDs(attrView.ID)
			avBlocks := treenode.GetBlockTrees(avBlockIDs)
			for _, avBlock := range avBlocks {
				passwordID, password := GetPathPasswordByPublishAccess(avBlock.BoxID, avBlock.Path, publishAccess)
				if (password == "" || CheckPublishAuthCookie(c, passwordID, password)) && CheckPathEnabledByPublishAccess(avBlock.BoxID, avBlock.Path, publishAccess) {
					avBlocksAccessable = true
					break
				}
			}
		}
		if avBlocksAccessable {
			filteredAttrViews = append(filteredAttrViews, attrView)
			avIDs = append(avIDs, attrView.ID)
		}
	}
	ret.AttrViews = filteredAttrViews
	ret.IAL[av.NodeAttrNameAvs] = strings.Join(avIDs, ",")

	bt := treenode.GetBlockTree(info.RootID)
	if bt != nil {
		passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
		if (password != "" && !CheckPublishAuthCookie(c, passwordID, password)) || !CheckPathEnabledByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
			ret.IAL["name"] = ""
			ret.IAL["alias"] = ""
			ret.IAL["memo"] = ""
			ret.IAL["bookmark"] = ""
			ret.IAL["tags"] = ""
			ret.RefCount = 0
			ret.RefIDs = []string{}
		}
	}
	return
}

func FilterContentByPublishAccess(c *gin.Context, publishAccess PublishAccess, box string, docPath string, content string, onlyIcon bool) (ret string) {
	ret = content

	// 密码访问
	passwordID, password := GetPathPasswordByPublishAccess(box, docPath, publishAccess)
	if password != "" {
		if !CheckPublishAuthCookie(c, passwordID, password) {
			if onlyIcon {
				passwordHTML := `<div class="protyle-password protyle-password--alert" data-node-id="%s">
	<span class="protyle-password__logo">🔒</span>
</div>`
				ret = fmt.Sprintf(passwordHTML, passwordID)
			} else {
				passwordHTML := `<div class="protyle-password" data-node-id="%s">
	<span class="protyle-password__logo">🔒</span>
	<label class="b3-form__icon protyle-password__content">
		<svg class="b3-form__icon-icon"><use xlink:href="#iconKey"></use></svg>
		<input type="text" class="b3-form__icon-input b3-text-field b3-form__icona-input" placeholder="%s"/>
		<svg class="protyle-password__button b3-form__icona-icon"><use xlink:href="#iconForward"></use></svg>
	</label>
</div>`
				ret = fmt.Sprintf(passwordHTML, passwordID, Conf.Language(283))
			}
		}
	}

	// 禁止访问
	ID := box
	if docPath != "/" {
		ID = strings.TrimSuffix(path.Base(docPath), ".sy")
	}
	if !CheckPathEnabledByPublishAccess(box, docPath, publishAccess) {
		if onlyIcon {
			forbiddenHTML := `<div class="protyle-password protyle-password--alert" data-node-id="%s">
	<span class="protyle-password__logo">🚫</span>
</div>`
			ret = fmt.Sprintf(forbiddenHTML, ID)
		} else {
			forbiddenHTML := `<div class="protyle-password protyle-password--forbidden" data-node-id="%s">
	<span class="protyle-password__logo">🚫</span>
	<div class="protyle-password__tip">%s</div>
</div>`
			ret = fmt.Sprintf(forbiddenHTML, ID, Conf.Language(284))
		}
	}
	return
}

func FilterEmbedBlocksByPublishAccess(c *gin.Context, publishAccess PublishAccess, embedBlocks []*EmbedBlock) (ret []*EmbedBlock) {
	ret = []*EmbedBlock{}
	for _, embedBlock := range embedBlocks {
		embedBlock.Block.Content = FilterContentByPublishAccess(c, publishAccess, embedBlock.Block.Box, embedBlock.Block.Path, embedBlock.Block.Content, false)
		ret = append(ret, embedBlock)
	}
	return
}

func FilterPathsByPublishAccess(c *gin.Context, publishAccess PublishAccess, paths []*Path) (ret []*Path) {
	ret = []*Path{}
	IDs := []string{}

	IDtoPathIndexMap := make(map[string]int)
	for i, path := range paths {
		IDs = append(IDs, path.ID)
		IDtoPathIndexMap[path.ID] = i
	}
	bts := treenode.GetBlockTrees(IDs)
	for _, bt := range bts {
		if bt == nil {
			continue
		}
		pathIndex := IDtoPathIndexMap[bt.ID]
		path := paths[pathIndex]
		passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
		if CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) && (password == "" || CheckPublishAuthCookie(c, passwordID, password)) {
			ret = append(ret, path)
		}
	}
	return
}

func FilterBlocksByPublishAccess(c *gin.Context, publishAccess PublishAccess, blocks []*Block) (ret []*Block) {
	ret = []*Block{}

	for _, block := range blocks {
		passwordID, password := GetPathPasswordByPublishAccess(block.Box, block.Path, publishAccess)
		if CheckPathVisibleByPublishAccess(block.Box, block.Path, publishAccess) && (c == nil || password == "" || CheckPublishAuthCookie(c, passwordID, password)) {
			ret = append(ret, block)
		}
	}
	return
}

func FilterBlockTreesByPublishAccess(publishAccess PublishAccess, bts map[string]*treenode.BlockTree) (ret map[string]*treenode.BlockTree) {
	ret = map[string]*treenode.BlockTree{}
	for id, bt := range bts {
		if CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
			ret[id] = bt
		}
	}
	return
}

func FilterRefDefsByPublishAccess(publishAccess PublishAccess, refDefs []*RefDefs) (retRefDefs []*RefDefs, originalRefBlockIDs map[string]string) {
	retRefDefs = []*RefDefs{}
	IDs := []string{}
	for _, refDef := range refDefs {
		IDs = append(IDs, refDef.RefID)
		IDs = append(IDs, refDef.DefIDs...)
	}
	IDs = gulu.Str.RemoveDuplicatedElem(IDs)
	bts := treenode.GetBlockTrees(IDs)
	bts = FilterBlockTreesByPublishAccess(publishAccess, bts)
	visibles := make(map[string]bool)
	for _, ID := range IDs {
		visibles[ID] = false
	}
	for _, bt := range bts {
		visibles[bt.ID] = true
	}
	for _, refDef := range refDefs {
		if !visibles[refDef.RefID] {
			continue
		}
		newDefIDs := []string{}
		for i, defID := range refDef.DefIDs {
			if visibles[defID] {
				newDefIDs = append(newDefIDs, refDef.DefIDs[i])
			}
		}
		refDef.DefIDs = newDefIDs
		if len(refDef.DefIDs) > 0 {
			retRefDefs = append(retRefDefs, refDef)
		}
	}
	originalRefBlockIDs = buildBacklinkListItemRefs(retRefDefs)
	return
}

func FilterConfByPublishAccess(publishAccess PublishAccess, appConf *AppConf) (ret *AppConf) {
	ret = appConf
	if appConf == nil {
		return
	}

	appConf.UILayout = FilterUILayoutByPublishAccess(publishAccess, appConf.UILayout)
	return
}

func FilterUILayoutByPublishAccess(publishAccess PublishAccess, uiLayout *conf.UILayout) (ret *conf.UILayout) {
	ret = uiLayout
	if uiLayout == nil {
		return
	}

	layout, ok := (*uiLayout)["layout"].(map[string]any)
	if !ok {
		return
	}
	layout = filterLayoutItemByPublishAccess(publishAccess, layout)
	(*ret)["layout"] = layout
	return
}

func filterLayoutItemByPublishAccess(publishAccess PublishAccess, item map[string]any) (ret map[string]any) {
	ret = item
	if item == nil {
		return
	}

	instanceItem, exists := item["instance"]
	if !exists {
		return
	}
	instance := instanceItem.(string)
	if instance == "Tab" {
		childrenItem, exists := item["children"]
		if !exists {
			return
		}
		children := childrenItem.(map[string]any)
		if children == nil {
			return
		}
		rootIdItem, exists := children["rootId"]
		if rootIdItem == nil {
			return
		}
		rootId := children["rootId"].(string)
		bt := treenode.GetBlockTree(rootId)
		if bt == nil {
			return
		}
		if !CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
			ret = nil
		}
	} else {
		childrenItem, exists := item["children"]
		if !exists {
			return
		}
		children := childrenItem.([]any)
		if children == nil {
			return
		}
		newChildren := []any{}
		updateTabs := false
		for _, childItem := range children {
			child := childItem.(map[string]any)
			if child == nil {
				return
			}
			child = filterLayoutItemByPublishAccess(publishAccess, child)
			if child != nil {
				newChildren = append(newChildren, child)
			} else {
				updateTabs = true
			}
		}
		if updateTabs {
			hasActive := false
			activeTimes := []int64{}
			for _, childItem := range newChildren {
				child := childItem.(map[string]any)
				activeTimeStr := child["activeTime"].(string)
				var activeTime int64
				if len(activeTimeStr) > 0 {
					activeTime, _ = strconv.ParseInt(activeTimeStr, 10, 64)
				}
				activeTimes = append(activeTimes, activeTime)
				if active, exists := child["active"]; exists && active.(bool) {
					hasActive = true
				}
			}
			if !hasActive && len(activeTimes) > 0 {
				// 如果原本激活的tab刚好被去掉了，就选择日期最新的一个tab激活
				maxIndex := 0
				for i, activeTime := range activeTimes {
					if activeTime > activeTimes[maxIndex] {
						maxIndex = i
					}
				}
				newChildren[maxIndex].(map[string]any)["active"] = true
			}
			if len(newChildren) == 0 {
				child := make(map[string]any)
				child["instance"] = "Tab"
				child["children"] = make(map[string]any)
				newChildren = append(newChildren, child)
			}
		}
		ret["children"] = newChildren
	}
	return
}

func FilterGraphByPublishAccess(publishAccess PublishAccess, nodes []*GraphNode, links []*GraphLink) (retNodes []*GraphNode, retLinks []*GraphLink) {
	retNodes = []*GraphNode{}
	retLinks = []*GraphLink{}
	ignoreNodeIDs := []string{}
	for _, node := range nodes {
		if CheckPathVisibleByPublishAccess(node.Box, node.Path, publishAccess) {
			retNodes = append(retNodes, node)
		} else {
			ignoreNodeIDs = append(ignoreNodeIDs, node.ID)
		}
	}
	for _, link := range links {
		ignore := false
		for _, ignoreNodeID := range ignoreNodeIDs {
			if ignoreNodeID == link.From || ignoreNodeID == link.To {
				ignore = true
				break
			}
		}
		if !ignore {
			retLinks = append(retLinks, link)
		}
	}
	return
}

func FilterTagsByPublishAccess(publishAccess PublishAccess, tags *Tags) (ret *Tags) {
	spans := sql.QueryTagSpans("")
	labelCounts := make(map[string]int)
	for _, span := range spans {
		if CheckPathVisibleByPublishAccess(span.Box, span.Path, publishAccess) {
			label := util.UnescapeHTML(span.Content)
			labelCounts[label] += 1
		}
	}

	ret = &Tags{}
	for _, tag := range *tags {
		tag := reassignTagCounts(tag, labelCounts)
		if tag != nil {
			*ret = append(*ret, tag)
		}
	}
	return
}

func reassignTagCounts(tag *Tag, counts map[string]int) (ret *Tag) {
	var newChildren Tags
	for _, child := range tag.Children {
		child = reassignTagCounts(child, counts)
		if child != nil {
			newChildren = append(newChildren, child)
		}
	}
	tag.Children = newChildren
	tag.Count = counts[tag.Label]
	if tag.Children == nil && tag.Count == 0 {
		return nil
	}
	return tag
}

func FilterLocalStorageByPublishAccess(publishAccess PublishAccess, localStorage map[string]any) (ret map[string]any) {
	ret = localStorage
	// 清空搜索历史记录
	searchKeysItem := ret["local-searchkeys"]
	if searchKeysItem != nil {
		searchKeys := searchKeysItem.(map[string]any)
		if searchKeys != nil {
			searchKeys["keys"] = []string{}
		}
	}
	searchAssetItem := ret["local-searchasset"]
	if searchAssetItem != nil {
		searchAsset := searchAssetItem.(map[string]any)
		if searchAsset != nil {
			searchAsset["k"] = ""
			searchAsset["keys"] = []string{}
		}
	}
	docInfoItem := ret["local-docinfo"]
	if docInfoItem != nil {
		docInfo := docInfoItem.(map[string]any)
		if docInfo != nil {
			idItem := docInfo["id"]
			if idItem != nil {
				id := idItem.(string)
				bt := treenode.GetBlockTree(id)
				if bt != nil {
					if !CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) {
						docInfo["id"] = ""
					}
				}
			}
		}
	}
	return
}

func FilterAssetContentByPublishAccess(c *gin.Context, publishAccess PublishAccess, assetContent []*AssetContent) (ret []*AssetContent) {
	validAssets := []string{}
	bts := treenode.GetBlockTreesByType("d")
	for _, bt := range bts {
		passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
		if CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) && (password == "" || CheckPublishAuthCookie(c, passwordID, password)) {
			assets, err := DocAssets(bt.ID, false)
			if err == nil {
				validAssets = append(validAssets, assets...)
			}
		}
	}

	ret = []*AssetContent{}
	for _, asset := range assetContent {
		if asset == nil {
			continue
		}
		for _, validAsset := range validAssets {
			if validAsset == asset.Path {
				ret = append(ret, asset)
			}
		}
	}
	return
}

func FilterRecentDocsByPublishAccess(c *gin.Context, publishAccess PublishAccess, recentDocs []*RecentDoc) (ret []*RecentDoc) {
	ret = []*RecentDoc{}
	for _, recentDoc := range recentDocs {
		bt := treenode.GetBlockTree(recentDoc.RootID)
		if bt != nil {
			passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
			if CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) && (passwordID == "" || CheckPublishAuthCookie(c, passwordID, password)) {
				ret = append(ret, recentDoc)
			}
		}
	}
	return
}

func FilterCriteriaByPublishAccess(c *gin.Context, publishAccess PublishAccess, criteria []*Criterion) (ret []*Criterion) {
	ret = []*Criterion{}
	// IDPath 元素可能是笔记本 ID、文档 ID，或 "笔记本ID/文档ID[.sy]" 路径串，这里统一解析出文档 ID
	blockIDs := map[string]struct{}{}
	for _, criterion := range criteria {
		for _, p := range criterion.IDPath {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			// 路径形式取末段并去掉 .sy 后缀
			id := strings.TrimSuffix(path.Base(p), ".sy")
			if id != "" && id != "." && id != "/" {
				blockIDs[id] = struct{}{}
			}
		}
	}
	blockIDsSlice := make([]string, 0, len(blockIDs))
	for id := range blockIDs {
		blockIDsSlice = append(blockIDsSlice, id)
	}
	blockTrees := treenode.GetBlockTrees(blockIDsSlice)
	for _, criterion := range criteria {
		accessible := false
		for _, p := range criterion.IDPath {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			id := strings.TrimSuffix(path.Base(p), ".sy")
			if id == "" || id == "." || id == "/" {
				continue
			}
			bt := blockTrees[id]
			if bt == nil {
				// 关联的文档不存在，视为不可访问
				accessible = false
				break
			}
			passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
			if !CheckPathVisibleByPublishAccess(bt.BoxID, bt.Path, publishAccess) || (passwordID != "" && !CheckPublishAuthCookie(c, passwordID, password)) {
				accessible = false
				break
			}
			accessible = true
		}
		if !accessible {
			// 若 IDPath 全部不可访问（或引用了不可见文档），整条丢弃，避免泄露 HPath
			continue
		}

		// 复制一份后再清空搜索/替换关键字，避免污染缓存
		cloned := *criterion
		cloned.K = ""
		cloned.R = ""
		ret = append(ret, &cloned)
	}
	return
}
