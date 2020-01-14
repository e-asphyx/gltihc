import { Gltihc } from "./gltihc.js";
const optionControls = [
    {
        label: "Maximum width",
        min: 0,
        prop: "maxWidth",
    },
    {
        label: "Maximum height",
        min: 0,
        prop: "maxHeight",
    },
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
const optionsVer = "0";
const filterNames = {
    color: "Color",
    gray: "Gray",
    src: "Source",
    rgba: "Set RGBA Component",
    seta: "Set Alpha",
    ycc: "Set YCC Component",
    prgb: "Permutate RGB",
    prgba: "Permutate RGBA",
    pycc: "Permutate YCC",
    copy: "Copy Component to Component",
    ctoa: "Copy Component to Alpha",
    mix: "Mixer",
    quant: "Quantize",
    qrgba: "Quantize RGBA Component",
    qycca: "Quantize YCCA Component",
    qy: "Quantize Y",
    inv: "Invert",
    invrgba: "Invert RGBA Component",
    inva: "Invert Alpha",
    invycc: "Invert YCC Component",
    gs: "Gray Scale",
    rasp: "BitRasp",
};
const operatorNames = {
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
function defaultToggleValues(opt) {
    return Object.assign({}, ...Object.keys(opt).map((n) => ({ [n]: true })));
}
function toggleValues(opt) {
    return Object.assign({}, ...opt.map((n) => ({ [n]: true })));
}
class App {
    constructor() {
        this.gltihc = new Gltihc();
        this.imageEl = document.getElementById("target-img");
        this.msgContainerEl = document.getElementById("msg-container");
        this.msgToastEl = document.getElementById("msg-toast");
        this.refreshBtn = document.getElementById("refresh-btn");
        this.storage = window.localStorage;
        this.gltihc.initDone.then(() => this.postWASMInit());
        let val = this.storage.getItem("gltihcVer");
        if (val === optionsVer) {
            val = this.storage.getItem("gltihc");
            if (val) {
                this.gltihc.options = JSON.parse(val);
            }
        }
        if (this.gltihc.options.filters instanceof Array) {
            this.filters = toggleValues(this.gltihc.options.filters);
        }
        else {
            this.filters = defaultToggleValues(filterNames);
        }
        if (this.gltihc.options.ops instanceof Array) {
            this.operators = toggleValues(this.gltihc.options.ops);
        }
        else {
            this.operators = defaultToggleValues(operatorNames);
        }
        this.initSettings();
    }
    refreshImage() {
        if (!this.gltihc.source) {
            this.showMessage("No image loaded");
            return;
        }
        this.showMessage("Processing...");
        // Redraw
        setTimeout(async () => {
            try {
                const blob = await this.gltihc.processImage();
                this.imageEl.src = URL.createObjectURL(blob);
                this.imageEl.style.removeProperty("display");
                this.showMessage(null);
                this.refreshBtn.disabled = false;
            }
            catch (err) {
                this.showMessage(err);
            }
        });
    }
    showMessage(msg) {
        if (msg) {
            this.msgToastEl.innerHTML = msg;
            this.msgContainerEl.style.removeProperty("display");
        }
        else {
            this.msgContainerEl.style.setProperty("display", "none");
        }
    }
    saveConfig() {
        this.storage.setItem("gltihc", JSON.stringify(this.gltihc.options));
        this.storage.setItem("gltihcVer", optionsVer);
    }
    formElement(label, id) {
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
    toggleElements(labels, values, opt, prefix) {
        return Object.keys(labels).map((key, i) => {
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
                        values[key] = ev.target.checked;
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
    initSettings() {
        var _a, _b;
        // Numeric options
        const optForm = document.getElementById("options-form");
        optionControls.forEach((o, i) => {
            var _a;
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
                        this.gltihc.options[o.prop] = Number(ev.target.value);
                        this.saveConfig();
                    });
                    return el;
                })());
                return el;
            })());
            (_a = optForm) === null || _a === void 0 ? void 0 : _a.appendChild(el);
        });
        // Filters
        const filtersForm = document.getElementById("filters-form");
        for (const child of this.toggleElements(filterNames, this.filters, "filters", "filter-input")) {
            (_a = filtersForm) === null || _a === void 0 ? void 0 : _a.appendChild(child);
        }
        // Operators
        const opsForm = document.getElementById("operators-form");
        for (const child of this.toggleElements(operatorNames, this.operators, "ops", "operator-input")) {
            (_b = opsForm) === null || _b === void 0 ? void 0 : _b.appendChild(child);
        }
    }
    postWASMInit() {
        var _a, _b, _c;
        const uploader = document.getElementById("uploader");
        (_a = uploader) === null || _a === void 0 ? void 0 : _a.addEventListener("change", (ev) => {
            var _a, _b;
            const file = (_b = (_a = ev.target) === null || _a === void 0 ? void 0 : _a.files) === null || _b === void 0 ? void 0 : _b[0];
            this.gltihc.source = file;
            this.refreshImage();
        });
        (_b = document.getElementById("upload-btn")) === null || _b === void 0 ? void 0 : _b.addEventListener("click", () => {
            var _a;
            (_a = uploader) === null || _a === void 0 ? void 0 : _a.click();
        });
        this.refreshBtn.addEventListener("click", () => this.refreshImage());
        const modalCtl = document.getElementById("modal-control");
        (_c = document.getElementById("settings-btn")) === null || _c === void 0 ? void 0 : _c.addEventListener("click", () => {
            modalCtl.checked = true;
        });
    }
}
// @ts-ignore
const app = new App();
//# sourceMappingURL=main.js.map