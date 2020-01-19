function setDisabled(n, disabled) {
    if (disabled) {
        n.setAttribute("disabled", "");
    }
    else {
        n.removeAttribute("disabled");
    }
}
const CORS_PROXY = "https://cors-anywhere.herokuapp.com/";
export default class MainPage {
    constructor(gltihc) {
        var _a, _b, _c, _d;
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
            this.message("No image loaded");
            setDisabled(uploadBtn, false);
        });
        const dz = this.content.getElementById("drop-zone");
        (_a = dz) === null || _a === void 0 ? void 0 : _a.addEventListener("dragenter", (ev) => {
            var _a;
            (_a = ev.target) === null || _a === void 0 ? void 0 : _a.classList.add("drag-over");
        });
        (_b = dz) === null || _b === void 0 ? void 0 : _b.addEventListener("dragleave", (ev) => {
            var _a;
            (_a = ev.target) === null || _a === void 0 ? void 0 : _a.classList.remove("drag-over");
        });
        (_c = dz) === null || _c === void 0 ? void 0 : _c.addEventListener("dragover", (ev) => {
            ev.preventDefault();
        });
        (_d = dz) === null || _d === void 0 ? void 0 : _d.addEventListener("drop", (ev) => {
            var _a;
            ev.preventDefault();
            (_a = ev.target) === null || _a === void 0 ? void 0 : _a.classList.remove("drag-over");
            const dt = ev.dataTransfer;
            if (!dt) {
                return;
            }
            const uri = dt.getData("text/uri-list") || dt.getData("text/plain");
            if (uri) {
                fetch(CORS_PROXY + uri).then((resp) => {
                    console.log(resp);
                    if (resp.ok) {
                        return resp.blob();
                    }
                    else {
                        throw new Error(resp.statusText);
                    }
                }).then((blob) => {
                    console.log(blob);
                    if (blob.type.match("^image/")) {
                        this.gltihc.source = blob;
                        this.refreshImage();
                    }
                }).catch((err) => this.message(String(err)));
                return;
            }
            for (const f of dt.files) {
                if (f.type.match("^image/")) {
                    this.gltihc.source = f;
                    this.refreshImage();
                    return;
                }
            }
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