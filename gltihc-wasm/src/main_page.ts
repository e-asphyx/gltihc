import { Gltihc } from "./gltihc.js";

function setDisabled(n: Element, disabled: boolean) {
    if (disabled) {
        n.setAttribute("disabled", "");
    } else {
        n.removeAttribute("disabled");
    }
}

export default class MainPage {
    public readonly content: DocumentFragment;
    public readonly settingsBtn: HTMLInputElement;

    private imageEl: HTMLImageElement;
    private messageElem: HTMLElement;
    private messageContainer: HTMLElement;
    private refreshBtn: HTMLInputElement;

    constructor(private gltihc: Gltihc) {
        const tpl = <HTMLTemplateElement>document.getElementById("main-tpl");

        this.content = <DocumentFragment>tpl.content.cloneNode(true);
        this.imageEl = <HTMLImageElement>this.content.getElementById("target-img");
        this.settingsBtn = <HTMLInputElement>this.content.getElementById("settings-btn");

        this.messageContainer = <HTMLElement>this.content.getElementById("message-box");
        this.messageElem = <HTMLElement>this.content.getElementById("message");

        this.refreshBtn = <HTMLInputElement>this.content.getElementById("refresh-btn");
        this.refreshBtn.addEventListener("click", () => this.refreshImage());

        const uploader = <HTMLInputElement>this.content.getElementById("uploader");
        uploader.addEventListener("change", (ev) => {
            const file = (<HTMLInputElement>ev.target)?.files?.[0];
            if (!file) {
                setDisabled(this.refreshBtn, true);
            } else {
                this.gltihc.source = file;
                this.refreshImage();
            }
        });

        const uploadBtn = <HTMLInputElement>this.content.getElementById("upload-btn");
        uploadBtn.addEventListener("click", () => uploader.click());

        this.gltihc.initDone.then(() => {
            setDisabled(uploadBtn, false);
        });
    }

    private refreshImage() {
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
            } catch (err) {
                this.message(err);
            }
            setDisabled(this.refreshBtn, false);
        });
    }

    private message(msg: string | null) {
        if (msg) {
            this.messageElem.innerHTML = msg;
            this.messageContainer.style.removeProperty("display");
        } else {
            this.messageContainer.style.setProperty("display", "none");
        }
    }
}
