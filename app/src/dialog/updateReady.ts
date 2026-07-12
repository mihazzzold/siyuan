import {Constants} from "../constants";
import {fetchPost} from "../util/fetch";
/// #if !BROWSER
import {ipcRenderer} from "electron";
/// #endif
import {showMessage} from "./message";

// 退出并安装已下载的更新（electron-updater）。
// 本地内核：先优雅关闭内核（execInstallPkg:1 = 跳过内核自带安装器），再退出安装。
// 瘦客户端（远程服务器）：不动远程内核，直接退出安装。
const restartForUpdate = () => {
    /// #if !BROWSER
    const isLocalKernel = ["127.0.0.1", "localhost"].includes(location.hostname);
    if (isLocalKernel) {
        fetchPost("/api/system/exit", {force: true, setCurrentWorkspace: false, execInstallPkg: 1}, () => {
            // дать ядру полностью завершиться, иначе на Windows файлы ядра могут быть заняты установщиком
            setTimeout(() => {
                ipcRenderer.send(Constants.SIYUAN_QUIT_UPDATE, location.port);
            }, 3000);
        });
    } else {
        ipcRenderer.send(Constants.SIYUAN_QUIT_UPDATE, location.port);
    }
    /// #endif
};

// 主进程通知「更新已下载」后，弹出可手动关闭的提示条，带「重启以更新」按钮。
export const showUpdateReady = (version: string) => {
    const lang = window.siyuan.languages;
    const title = (lang.updateReady || "🎉 Update downloaded, restart to apply") +
        (version ? " (v" + version + ")" : "");
    const btnLabel = lang.updateRestartNow || "Restart";
    const btnId = "siyuanUpdateRestartBtn";
    const html = `<div class="fn__flex" style="align-items: center;">
<span style="margin-right: 12px;">${title}</span>
<button class="b3-button b3-button--small" id="${btnId}" style="white-space: nowrap;">${btnLabel}</button>
</div>`;
    // timeout 0：常驻，带关闭按钮；固定 id 避免重复弹出多条
    showMessage(html, 0, "info", "siyuan-update-ready");
    const btn = document.getElementById(btnId);
    if (btn) {
        btn.addEventListener("click", () => {
            btn.setAttribute("disabled", "disabled");
            restartForUpdate();
        });
    }
};
