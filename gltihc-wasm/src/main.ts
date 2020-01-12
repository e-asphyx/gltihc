import { Gltihc } from "./gltihc.js";

class App {
    private gltihc = new Gltihc();
    private imageEl = <HTMLImageElement>document.getElementById("targetImg");
    private msgContainerEl = <HTMLElement>document.getElementById("msgContainer");
    private msgToastEl = <HTMLElement>document.getElementById("msgToast");

    constructor() {
        this.gltihc.initDone.then(() => this.init());
    }

    private async refreshImage() {
        this.showMessage("Processing...");
        try {
            let blob = await this.gltihc.processImage();
            this.imageEl.src = URL.createObjectURL(blob);
            this.imageEl.style.removeProperty("display");
            this.showMessage(null);
        } catch (err) {
            this.showMessage(err);
        }
    }

    private showMessage(msg: string | null) {
        if (msg) {
            this.msgToastEl.innerHTML = msg;
            this.msgContainerEl.style.removeProperty("display");
        } else {
            this.msgContainerEl.style.setProperty("display", "none");
        }
    }

    private init() {
        var uploader = <HTMLInputElement>document.getElementById("uploader");
        uploader?.addEventListener("change", (ev) => {
            this.gltihc.source = (<HTMLInputElement>ev.target)?.files?.[0];
            this.refreshImage();
        });
        document.getElementById("uploadBtn")?.addEventListener("click", () => {
            uploader?.click();
        });
        document.getElementById("refreshBtn")?.addEventListener("click", () => this.refreshImage());
    };
}

// @ts-ignore
var app = new App();



