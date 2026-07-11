import type {SettingTabBuilder} from "../setting/builder";
import {Constants} from "../../constants";
import {fetchPost} from "../../util/fetch";
import {confirmDialog} from "../../dialog/confirmDialog";
import {showMessage} from "../../dialog/message";
import {processSync} from "../../dialog/processSystem";
import {writeText} from "../../protyle/util/compatibility";
import {bindSyncCloudListEvent, renderSyncCloudList, setKey} from "../../sync/syncGuide";
import {Dialog} from "../../dialog";
import {genConfigItemMainHtml, genConfigItemName} from "../render/fragments";
import {getSyncProviderConfigKeywords} from "./syncUi";
import {patchSyncConfig} from "./syncRuntime";
import {openHistory} from "../../history/history";

const registerSyncGroup = (tab: SettingTabBuilder) => {
    const group = tab.group("sync", window.siyuan.languages.configGroupSync);

    group.select("sync.provider", {
        title: window.siyuan.languages.syncProvider,
        desc: window.siyuan.languages.syncProviderTip,
        options: [
            // 官方 SiYuan 云端（value 0，指向 fork 上游服务）暂不提供，仅保留自建 S3/WebDAV/本地
            {value: 2, label: "S3"},
            {value: 3, label: "WebDAV"},
            ...(["std", "docker"].includes(window.siyuan.config.system.container) ? [{value: 4, label: window.siyuan.languages.localFileSystem}] : []),
        ],
        save: (value) => patchSyncConfig("sync.provider", value),
    });
    group.slot({
        key: "syncProviderConfig",
        keywords: getSyncProviderConfigKeywords(),
        html: () => '<div id="syncProviderConfig" class="b3-label config-item"></div>',
    });
    group.slot({
        key: "cloudSpace",
        keywords: [window.siyuan.languages.cloudStorage, window.siyuan.languages.trafficStat, window.siyuan.languages.backup],
        html: () => '<div id="cloudSpace" class="b3-label config-item"></div>',
    });
    group.switch("sync.enabled", {
        title: window.siyuan.languages.openSyncTip1,
        desc: window.siyuan.languages.openSyncTip2,
        save: (value) => patchSyncConfig("sync.enabled", value),
    });
    group.switch("sync.generateConflictDoc", {
        title: window.siyuan.languages.generateConflictDoc,
        desc: window.siyuan.languages.generateConflictDocTip,
        save: (value) => patchSyncConfig("sync.generateConflictDoc", value),
    });
    group.select("sync.mode", {
        title: window.siyuan.languages.syncMode,
        desc: window.siyuan.languages.syncModeTip,
        options: [
            {value: 1, label: window.siyuan.languages.syncMode1},
            {value: 2, label: window.siyuan.languages.syncMode2},
            {value: 3, label: window.siyuan.languages.syncMode3},
        ],
        save: (value) => patchSyncConfig("sync.mode", value),
    });
    group.number("sync.interval", {
        title: window.siyuan.languages.syncInterval,
        desc: window.siyuan.languages.syncIntervalTip,
        min: 30,
        max: 43200,
        unit: window.siyuan.languages.second,
        save: (value) => patchSyncConfig("sync.interval", value),
    });
    group.switch("sync.perception", {
        title: window.siyuan.languages.syncPerception,
        desc: window.siyuan.languages.syncPerceptionTip,
        save: (value) => patchSyncConfig("sync.perception", value),
    });
    group.slot({
        key: "syncCloudDir",
        keywords: [window.siyuan.languages.cloudSyncDir, window.siyuan.languages.cloudSyncDirTip, window.siyuan.languages.config],
        html: () => `<div class="b3-label config-item" id="syncCloudDirBlock">
    <div class="fn__flex config-wrap">
        ${genConfigItemMainHtml(window.siyuan.languages.cloudSyncDir, window.siyuan.languages.cloudSyncDirTip)}
        <div class="fn__space"></div>
        <button class="b3-button b3-button--outline fn__flex-center fn__size200" data-action="config">
            <svg><use xlink:href="#iconSettings"></use></svg>${window.siyuan.languages.config}
        </button>
    </div>
    <div id="syncCloudList" class="fn__none"></div>
</div>`,
        afterMount: mountSyncCloudDir,
    });
    group.slot({
        key: "syncCloudBackup",
        keywords: [
            window.siyuan.languages.cloudBackup,
            window.siyuan.languages.cloudBackupTip,
            window.siyuan.languages.dataSnapshot,
        ],
        html: () => `<div class="b3-label config-item" id="syncCloudBackupBlock">
    <div class="fn__flex config-wrap">
        ${genConfigItemMainHtml(window.siyuan.languages.cloudBackup, window.siyuan.languages.cloudBackupTip)}
        <div class="fn__space"></div>
        <button class="b3-button b3-button--outline fn__flex-center fn__size200" id="openCloudBackup">
            <svg><use xlink:href="#iconHistory"></use></svg>${window.siyuan.languages.dataSnapshot}
        </button>
    </div>
</div>`,
        afterMount: (root) => {
            root.querySelector("#openCloudBackup")?.addEventListener("click", () => {
                openHistory(window.siyuan.ws.app, "repo");
            });
        },
    });
};

const mountSyncCloudDir = (root: HTMLElement) => {
    const cloudListElement = root.querySelector("#syncCloudList");
    if (cloudListElement) {
        bindSyncCloudListEvent(cloudListElement);
        root.querySelector('#syncCloudDirBlock [data-action="config"]')?.addEventListener("click", () => {
            const hidden = cloudListElement.classList.toggle("fn__none");
            if (!hidden) {
                renderSyncCloudList(cloudListElement, true);
            }
        });
    }
};

const registerRepoGroup = (tab: SettingTabBuilder) => {
    const group = tab.group("repo", window.siyuan.languages.configGroupLocalDataRepo);

    group.slot({
        key: "repoKey",
        keywords: [
            window.siyuan.languages.dataRepoKey,
            window.siyuan.languages.dataRepoKeyTip1,
            window.siyuan.languages.dataRepoKeyTip2,
            window.siyuan.languages.importKey,
            window.siyuan.languages.genKey,
            window.siyuan.languages.genKeyByPW,
            window.siyuan.languages.copyKey,
            window.siyuan.languages.resetRepo,
        ],
        html: () => `<div class="fn__flex b3-label config-item config-wrap">
    <div class="fn__flex-1 fn__flex-center">
        ${genConfigItemName(window.siyuan.languages.dataRepoKey)}
        <div class="fn__hr--small"></div>
        <div class="b3-label__text">
            ${window.siyuan.languages.dataRepoKeyTip1}
            <div class="fn__hr--small"></div>
            <span class="ft__error">${window.siyuan.languages.dataRepoKeyTip2}</span>
        </div>
    </div>
    <div class="fn__space"></div>
    <div class="fn__size200 fn__flex-center fn__none" id="repoKeyActionsEmpty">
        <button class="b3-button b3-button--outline fn__block" id="importKey"><svg><use xlink:href="#iconDownload"></use></svg>${window.siyuan.languages.importKey}</button>
        <div class="fn__hr"></div>
        <button class="b3-button b3-button--outline fn__block" id="initKey"><svg><use xlink:href="#iconLock"></use></svg>${window.siyuan.languages.genKey}</button>
        <div class="fn__hr"></div>
        <button class="b3-button b3-button--outline fn__block" id="initKeyByPW"><svg><use xlink:href="#iconHand"></use></svg>${window.siyuan.languages.genKeyByPW}</button>
    </div>
    <div class="fn__size200 fn__flex-center fn__none" id="repoKeyActionsSet">
        <button class="b3-button b3-button--outline fn__block" id="copyKey"><svg><use xlink:href="#iconCopy"></use></svg>${window.siyuan.languages.copyKey}</button>
        <div class="fn__hr"></div>
        <button class="b3-button b3-button--outline fn__block" id="resetRepo"><svg><use xlink:href="#iconUndo"></use></svg>${window.siyuan.languages.resetRepo}</button>
    </div>
</div>`,
        afterMount: mountRepoKey,
    });
    group.stack({
        key: "repoPurge",
        keywords: [
            window.siyuan.languages.dataRepoPurge,
            window.siyuan.languages.dataRepoPurgeTip,
            window.siyuan.languages.dataRepoAutoPurgeIndexRetentionDays,
            window.siyuan.languages.dataRepoAutoPurgeRetentionIndexesDaily,
        ],
        afterMount: (root) => {
            root.querySelector("#purgeRepo")?.addEventListener("click", () => {
                confirmDialog("♻️ " + window.siyuan.languages.dataRepoPurge, window.siyuan.languages.dataRepoPurgeConfirm, () => {
                    fetchPost("/api/repo/purgeRepo");
                });
            });
        },
    }, (stack) => {
        stack.title(window.siyuan.languages.dataRepoPurge);
        stack.desc(window.siyuan.languages.dataRepoPurgeTip);
        stack.button({
            id: "purgeRepo",
            label: window.siyuan.languages.purge,
            icon: "iconTrashcan",
        });
        stack.number("repo.indexRetentionDays", {
            desc: window.siyuan.languages.dataRepoAutoPurgeIndexRetentionDays,
            min: 1,
        });
        stack.number("repo.retentionIndexesDaily", {
            desc: window.siyuan.languages.dataRepoAutoPurgeRetentionIndexesDaily,
            min: 1,
        });
    });
};

const mountRepoKey = (root: HTMLElement) => {
    const emptyElement = root.querySelector("#repoKeyActionsEmpty");
    const setElement = root.querySelector("#repoKeyActionsSet");
    const toggleRepoKeyActions = () => {
        const hasKey = Boolean(window.siyuan.config.repo.key);
        emptyElement?.classList.toggle("fn__none", hasKey);
        setElement?.classList.toggle("fn__none", !hasKey);
    };
    toggleRepoKeyActions();
    root.querySelector("#importKey")?.addEventListener("click", () => {
        const passwordDialog = new Dialog({
            title: "🔑 " + window.siyuan.languages.key,
            content: `<div class="b3-dialog__content">
    <textarea spellcheck="false" style="resize: vertical;" class="b3-text-field fn__block" placeholder="${window.siyuan.languages.keyPlaceholder}"></textarea>
</div>
<div class="b3-dialog__action">
    <button class="b3-button b3-button--cancel">${window.siyuan.languages.cancel}</button><div class="fn__space"></div>
    <button class="b3-button b3-button--text">${window.siyuan.languages.confirm}</button>
</div>`,
            width: "520px",
        });
        passwordDialog.element.setAttribute("data-key", Constants.DIALOG_PASSWORD);
        const textAreaElement = passwordDialog.element.querySelector("textarea");
        textAreaElement.focus();
        const btnsElement = passwordDialog.element.querySelectorAll(".b3-button");
        btnsElement[0].addEventListener("click", () => {
            passwordDialog.destroy();
        });
        btnsElement[1].addEventListener("click", () => {
            fetchPost("/api/repo/importRepoKey", {key: textAreaElement.value}, (response) => {
                window.siyuan.config.repo.key = response.data.key;
                toggleRepoKeyActions();
                passwordDialog.destroy();
            });
        });
    });
    root.querySelector("#initKey")?.addEventListener("click", () => {
        confirmDialog("🔑 " + window.siyuan.languages.genKey, window.siyuan.languages.initRepoKeyTip, () => {
            fetchPost("/api/repo/initRepoKey", {}, (response) => {
                window.siyuan.config.repo.key = response.data.key;
                toggleRepoKeyActions();
            });
        });
    });
    root.querySelector("#initKeyByPW")?.addEventListener("click", () => {
        setKey(false, () => {
            toggleRepoKeyActions();
        });
    });
    root.querySelector("#copyKey")?.addEventListener("click", () => {
        writeText(window.siyuan.config.repo.key);
        showMessage(window.siyuan.languages.copied);
    });
    root.querySelector("#resetRepo")?.addEventListener("click", () => {
        confirmDialog("⚠️ " + window.siyuan.languages.resetRepo, window.siyuan.languages.resetRepoTip, () => {
            fetchPost("/api/repo/resetRepo", {}, () => {
                window.siyuan.config.repo.key = "";
                window.siyuan.config.sync.enabled = false;
                processSync();
                toggleRepoKeyActions();
            });
        });
    });
};

// registerGitBackupGroup 一键将 data/ 目录单向备份（提交并推送）到自建 Git 仓库。
// 说明文本暂以俄语硬编码（本 fork 默认俄语），后续可迁移到 i18n。
const registerGitBackupGroup = (tab: SettingTabBuilder) => {
    const group = tab.group("gitBackup", "Git-бэкап");
    group.slot({
        key: "gitBackupMain",
        keywords: ["git", "backup", "бэкап", "репозиторий", "токен"],
        html: () => `<div class="b3-label config-item">
    <div class="b3-label__text">Односторонний бэкап заметок в ваш Git-репозиторий (push). Данные шифрует Git-хостинг; используйте приватный репозиторий.</div>
    <div class="fn__hr"></div>
    <label class="fn__flex" style="align-items:center"><input class="b3-switch" type="checkbox" data-gb="enabled"><span class="fn__space"></span><span class="ft__on-surface">Включить Git-бэкап</span></label>
    <div class="fn__hr"></div>
    <input class="b3-text-field fn__block" data-gb="repoURL" placeholder="URL репозитория, напр. https://github.com/user/notes.git">
    <div class="fn__hr"></div>
    <input class="b3-text-field fn__block" data-gb="token" type="password" placeholder="Токен доступа (PAT). Оставьте пустым, чтобы не менять">
    <div class="fn__hr"></div>
    <div class="fn__flex">
        <input class="b3-text-field fn__flex-1" data-gb="branch" placeholder="Ветка (main)">
        <span class="fn__space"></span>
        <input class="b3-text-field fn__flex-1" data-gb="username" placeholder="Имя автора / логин (необяз.)">
    </div>
    <div class="fn__hr"></div>
    <input class="b3-text-field fn__block" data-gb="email" placeholder="E-mail автора коммита (необяз.)">
    <div class="fn__hr"></div>
    <label class="fn__flex" style="align-items:center"><input class="b3-switch" type="checkbox" data-gb="autoEnabled"><span class="fn__space"></span><span class="ft__on-surface">Авто-отправка каждые</span><span class="fn__space"></span><input class="b3-text-field" style="width:64px" type="number" min="5" data-gb="autoMinutes"><span class="fn__space"></span><span class="ft__on-surface">мин</span></label>
    <div class="fn__hr"></div>
    <div class="fn__flex">
        <button class="b3-button b3-button--outline fn__flex-center" id="gitBackupSave"><svg><use xlink:href="#iconDownload"></use></svg>Сохранить</button>
        <span class="fn__space"></span>
        <button class="b3-button b3-button--outline fn__flex-center" id="gitBackupNow"><svg><use xlink:href="#iconUpload"></use></svg>Отправить в Git</button>
        <span class="fn__space"></span>
        <button class="b3-button b3-button--outline fn__flex-center" id="gitBackupRestore"><svg><use xlink:href="#iconDownload"></use></svg>Восстановить из Git</button>
    </div>
    <div class="b3-label__text" style="margin-top:6px">Восстановление объединяет данные репозитория с текущими: отсутствующие документы добавляются, ваши локальные не перезаписываются.</div>
</div>`,
        afterMount: (root) => {
            const val = (name: string) => root.querySelector(`[data-gb="${name}"]`) as HTMLInputElement;
            fetchPost("/api/gitbackup/getConf", {}, (response) => {
                const c = response.data.gitBackup;
                val("enabled").checked = c.enabled;
                val("repoURL").value = c.repoURL || "";
                val("branch").value = c.branch || "main";
                val("username").value = c.username || "";
                val("email").value = c.email || "";
                val("autoEnabled").checked = c.autoEnabled;
                val("autoMinutes").value = String(c.autoMinutes || 30);
                if (response.data.tokenSet) {
                    val("token").placeholder = "Токен сохранён — оставьте пустым, чтобы не менять";
                }
            });
            const save = (cb?: () => void) => {
                fetchPost("/api/gitbackup/setConf", {
                    enabled: val("enabled").checked,
                    repoURL: val("repoURL").value.trim(),
                    token: val("token").value,
                    branch: val("branch").value.trim() || "main",
                    username: val("username").value.trim(),
                    email: val("email").value.trim(),
                    autoEnabled: val("autoEnabled").checked,
                    autoMinutes: parseInt(val("autoMinutes").value) || 30,
                }, () => {
                    val("token").value = "";
                    if (cb) {
                        cb();
                    } else {
                        showMessage("Сохранено");
                    }
                });
            };
            const doRestore = () => {
                fetchPost("/api/gitbackup/restore", {}, (response) => {
                    showMessage(response.code === 0 ? `Восстановлено файлов: ${response.data.restored}. Данные объединены.` : response.msg, response.code === 0 ? 5000 : 7000, response.code === 0 ? "info" : "error");
                });
            };
            // После сохранения проверяем, есть ли в репозитории бэкап, и предлагаем восстановить (с объединением)
            const promptRestoreIfRemoteHasBackup = () => {
                fetchPost("/api/gitbackup/checkRemote", {}, (response) => {
                    if (response.code === 0 && response.data.hasBackup) {
                        confirmDialog("♻️ Восстановление из Git", "В репозитории найден бэкап SiYuan. Восстановить его и объединить с текущими заметками? Отсутствующие документы будут добавлены, ваши локальные — сохранены (без перезаписи).", () => {
                            doRestore();
                        });
                    }
                });
            };
            root.querySelector("#gitBackupSave")?.addEventListener("click", () => save(() => {
                showMessage("Сохранено");
                promptRestoreIfRemoteHasBackup();
            }));
            root.querySelector("#gitBackupNow")?.addEventListener("click", () => {
                // сначала сохраняем настройки, затем запускаем отправку
                save(() => {
                    fetchPost("/api/gitbackup/backup", {}, (response) => {
                        showMessage(response.code === 0 ? "Отправлено в Git" : response.msg, response.code === 0 ? 3000 : 7000, response.code === 0 ? "info" : "error");
                    });
                });
            });
            root.querySelector("#gitBackupRestore")?.addEventListener("click", () => {
                // сохраняем настройки, затем восстанавливаем с подтверждением
                save(() => {
                    confirmDialog("♻️ Восстановление из Git", "Восстановить данные из репозитория и объединить с текущими заметками? Отсутствующие документы будут добавлены, ваши локальные — сохранены (без перезаписи).", () => {
                        doRestore();
                    });
                });
            });
        },
    });
};

export const registerSyncTab = (tab: SettingTabBuilder) => {
    // 官方账号（ld246）登录/注册区块暂不启用，后续接入自建服务时再恢复
    // registerAccountGroup(tab);
    registerSyncGroup(tab);
    registerGitBackupGroup(tab);
    registerRepoGroup(tab);
};
