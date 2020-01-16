import { Gltihc } from "./gltihc.js";
import MainPage from "./main_page.js";
import SettingsPage from "./settings_page.js";

const optionsVer = "0";

function createPage(frag: Node): Node {
    const el = document.createElement("div");
    el.className = "app-container";
    el.appendChild(frag);
    return el;
}

type Page = "main" | "settings";
type Pages = {
    [page in Page]: {
        node: Node,
        title: string,
    };
};

class App {
    private gltihc = new Gltihc();
    private mainPage: MainPage;
    private settingsPage: SettingsPage;

    private storage = window.localStorage;
    private root = document.getElementById("root");

    private pages: Pages;
    private currentPage: Page = "main";

    constructor() {
        // Load config
        let val = this.storage.getItem("gltihcVer");
        if (val === optionsVer) {
            val = this.storage.getItem("gltihc");
            if (val) {
                this.gltihc.options = JSON.parse(val);
            }
        }

        this.mainPage = new MainPage(this.gltihc);
        this.settingsPage = new SettingsPage(this.gltihc);
        this.pages = {
            main: {
                node: createPage(this.mainPage.content),
                title: "Gltihc",
            },
            settings: {
                node: createPage(this.settingsPage.content),
                title: "Gltihc | Settings",
            },
        };

        this.settingsPage.onChange = () => {
            this.storage.setItem("gltihc", JSON.stringify(this.gltihc.options));
            this.storage.setItem("gltihcVer", optionsVer);
        };

        window.addEventListener("popstate", (ev) => {
            if (ev.state) {
                this.showPage(ev.state);
            }
        });

        this.mainPage.settingsBtn.addEventListener("click", () => this.push("settings"));
        this.settingsPage.backBtn.addEventListener("click", () => history.back());

        // Display main page
        this.root?.appendChild(this.pages[this.currentPage].node);
        this.push(this.currentPage);
    }

    private push(page: Page) {
        history.pushState(page, this.pages[page].title);
        this.showPage(page);
    }

    private showPage(page: Page) {
        if (this.root && page !== this.currentPage) {
            this.root.replaceChild(this.pages[page].node, this.pages[this.currentPage].node);
        }
        this.currentPage = page;
    }
}

// @ts-ignore
const app = new App();
