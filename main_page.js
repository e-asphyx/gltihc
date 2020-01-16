function setDisabled(n, disabled) {
    if (disabled) {
        n.setAttribute("disabled", "");
    }
    else {
        n.removeAttribute("disabled");
    }
}
export default class MainPage {
    constructor(gltihc) {
        this.gltihc = gltihc;
        const tpl = document.getElementById("main-tpl");
        this.content = tpl.content.cloneNode(true);
        this.imageEl = this.content.getElementById("target-img");
        this.settingsBtn = this.content.getElementById("settings-btn");
        this.messageContainer = this.content.getElementById("message-box");
        this.messageElem = this.content.getElementById("message");
        this.refreshBtn = this.content.getElementById("refresh-btn");
        this.refreshBtn.addEventListener("click", () => this.refreshImage());
        const uploader = this.content.getElementById("uploader");
        uploader.addEventListener("change", (ev) => {
            var _a, _b;
            const file = (_b = (_a = ev.target) === null || _a === void 0 ? void 0 : _a.files) === null || _b === void 0 ? void 0 : _b[0];
            if (!file) {
                setDisabled(this.refreshBtn, true);
            }
            else {
                this.gltihc.source = file;
                this.refreshImage();
            }
        });
        const uploadBtn = this.content.getElementById("upload-btn");
        uploadBtn.addEventListener("click", () => uploader.click());
        this.gltihc.initDone.then(() => {
            setDisabled(uploadBtn, false);
        });
    }
    refreshImage() {
        if (!this.gltihc.source) {
            this.message("No image loaded");
            return;
        }
        setDisabled(this.refreshBtn, true);
        this.message("Processing...");
        // Redraw
        setTimeout(async () => {
            try {
                const blob = await this.gltihc.processImage();
                this.imageEl.src = URL.createObjectURL(blob);
                this.imageEl.style.removeProperty("display");
                this.message(null);
            }
            catch (err) {
                this.message(err);
            }
            setDisabled(this.refreshBtn, false);
        });
    }
    message(msg) {
        if (msg) {
            this.messageElem.innerHTML = msg;
            this.messageContainer.style.removeProperty("display");
        }
        else {
            this.messageContainer.style.setProperty("display", "none");
        }
    }
}
//# sourceMappingURL=main_page.js.map