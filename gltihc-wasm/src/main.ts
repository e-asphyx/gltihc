import { Gltihc, Option } from "./gltihc.js";

const optionControls: Array<{
    label: string;
    min: number;
    max?: number;
    step?: number;
    prop: Option;
}> = [
        {
            label: "Block size",
            min: 1,
            prop: "blockSize",
        },
        {
            label: "Minimum iterations",
            min: 1,
            prop: "minIterations",
        },
        {
            label: "Maximum iterations",
            min: 1,
            prop: "maxIterations",
        },
        {
            label: "Minimum filter chain length",
            min: 1,
            prop: "minFilters",
        },
        {
            label: "Maximum filter chain length",
            min: 1,
            prop: "maxFilters",
        },
        {
            label: "Minimum segment size",
            max: 100,
            min: 0.01,
            prop: "minSegmentSize",
            step: 0.01,
        },
        {
            label: "Maximum segment size",
            max: 100,
            min: 0.01,
            prop: "maxSegmentSize",
            step: 0.01,
        },
    ];

interface ToggleOptions {
    [name: string]: string;
}

const filterNames: ToggleOptions = {
    color: "Color",
    gray: "Gray",
    src: "Source",
    rgba: "Set RGBA Component",
    seta: "Set Alpha",
    ycc: "Set YCC Component",
    prgb: "Permutate RGB",
    prgba: "Permutate RGBA",
    pycc: "Permutate YCC",
    copy: "Copy Component",
    ctoa: "Copy to Alpha",
    mix: "Mixer",
    quant: "Quantize",
    qrgba: "Quantize RGBA",
    qycca: "Quantize YCCA",
    qy: "Quantize Y",
    inv: "Invert",
    invrgba: "Invert RGBA Component",
    invycc: "Invert YCC Component",
    gs: "Gray Scale",
    rasp: "BitRasp",
};

const operatorNames: ToggleOptions = {
    cmp: "Compose",
    src: "Replace",
    add: "Add",
    addrgbm: "Add RGB Modulo 256",
    addyccm: "Add YCC Modulo 256",
    mulrgb: "Multiply RGB",
    mulycc: "Multiply YCC",
    xorrgb: "Xor RGB",
    xorycc: "Xor YCC",
};

interface ToggleValues {
    [n: string]: boolean;
}

function defaultToggleValues(opt: ToggleOptions): ToggleValues {
    return Object.assign({}, ...Object.keys(opt).map<ToggleValues>((n) => ({ [n]: true })));
}

function toggleValues(opt: string[]): ToggleValues {
    return Object.assign({}, ...opt.map<ToggleValues>((n) => ({ [n]: true })));
}

class App {
    private gltihc = new Gltihc();
    private imageEl = <HTMLImageElement>document.getElementById("target-img");
    private msgContainerEl = <HTMLElement>document.getElementById("msg-container");
    private msgToastEl = <HTMLElement>document.getElementById("msg-toast");
    private refreshBtn = <HTMLInputElement>document.getElementById("refresh-btn");
    private filters: ToggleValues;
    private operators: ToggleValues;
    private storage = window.localStorage;

    constructor() {
        this.gltihc.initDone.then(() => this.postWASMInit());
        const val = this.storage.getItem("gltihc");
        if (val) {
            this.gltihc.options = JSON.parse(val);
        }

        if (this.gltihc.options.filters instanceof Array) {
            this.filters = toggleValues(this.gltihc.options.filters);
        } else {
            this.filters = defaultToggleValues(filterNames);
        }
        if (this.gltihc.options.ops instanceof Array) {
            this.operators = toggleValues(this.gltihc.options.ops);
        } else {
            this.operators = defaultToggleValues(operatorNames);
        }

        this.initSettings();
    }

    private async refreshImage() {
        this.showMessage("Processing...");
        try {
            const blob = await this.gltihc.processImage();
            this.imageEl.src = URL.createObjectURL(blob);
            this.imageEl.style.removeProperty("display");
            this.showMessage(null);
            this.refreshBtn.disabled = false;
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

    private saveConfig() {
        this.storage.setItem("gltihc", JSON.stringify(this.gltihc.options));
    }

    private formElement(label: string, id: string): Node {
        const el = document.createElement("div");
        el.className = "row responsive-label";
        el.appendChild((() => {
            const el = document.createElement("div");
            el.className = "col-sm-12 col-md-3";
            el.appendChild((() => {
                const el = document.createElement("label");
                el.setAttribute("for", id);
                el.innerHTML = label;
                return el;
            })());
            return el;
        })());
        return el;
    }

    private toggleElements(labels: ToggleOptions, values: ToggleValues, opt: Option, prefix: string): Node[] {
        return Object.keys(labels).map<Node>((key, i) => {
            const optId = `${prefix}-${i}`;
            const el = this.formElement(labels[key], optId);
            el.appendChild((() => {
                const el = document.createElement("div");
                el.className = "col-sm-12 col-md";
                el.appendChild((() => {
                    const el = document.createElement("input");
                    el.type = "checkbox";
                    el.id = optId;
                    el.checked = values[key] || false;
                    el.addEventListener("change", (ev) => {
                        values[key] = (<HTMLInputElement>ev.target).checked;
                        this.gltihc.options[opt] = Object.keys(values).filter((key) => values[key]);
                        this.saveConfig();
                    });
                    return el;
                })());
                return el;
            })());
            return el;
        });
    }

    private initSettings() {
        // Numeric options
        const optForm = document.getElementById("options-form");
        optionControls.forEach((o, i) => {
            const optId = `option-input-${i}`;
            const el = this.formElement(o.label, optId);
            el.appendChild((() => {
                const el = document.createElement("div");
                el.className = "col-sm-12 col-md";
                el.appendChild((() => {
                    const el = document.createElement("input");
                    el.type = "number";
                    el.id = optId;
                    el.min = String(o.min);
                    if (o.max) {
                        el.max = String(o.max);
                    }
                    if (o.step) {
                        el.step = String(o.step);
                    }
                    el.value = String(this.gltihc.options[o.prop]);
                    el.addEventListener("change", (ev) => {
                        this.gltihc.options[o.prop] = Number((<HTMLInputElement>ev.target).value);
                        this.saveConfig();
                    });
                    return el;
                })());
                return el;
            })());
            optForm?.appendChild(el);
        });

        // Filters
        const filtersForm = document.getElementById("filters-form");
        for (const child of this.toggleElements(filterNames, this.filters, "filters", "filter-input")) {
            filtersForm?.appendChild(child);
        }

        // Operators
        const opsForm = document.getElementById("operators-form");
        for (const child of this.toggleElements(operatorNames, this.operators, "ops", "operator-input")) {
            opsForm?.appendChild(child);
        }
    }

    private postWASMInit() {
        const uploader = <HTMLInputElement>document.getElementById("uploader");
        uploader?.addEventListener("change", (ev) => {
            this.gltihc.source = (<HTMLInputElement>ev.target)?.files?.[0];
            this.refreshImage();
        });
        document.getElementById("upload-btn")?.addEventListener("click", () => {
            uploader?.click();
        });
        this.refreshBtn.addEventListener("click", () => this.refreshImage());

        const modalCtl = <HTMLInputElement>document.getElementById("modal-control");
        document.getElementById("settings-btn")?.addEventListener("click", () => {
            modalCtl.checked = true;
        });
    }
}

// @ts-ignore
const app = new App();
