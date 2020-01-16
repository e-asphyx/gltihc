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
const filterNames = {
    color: "Color",
    gray: "Gray",
    src: "Source",
    rgba: "Set RGBA component",
    seta: "Set alpha",
    ycc: "Set YCC component",
    prgb: "Permutate RGB",
    prgba: "Permutate RGBA",
    pycc: "Permutate YCC",
    copy: "Copy component",
    ctoa: "Copy component to alpha",
    mix: "Mixer",
    quant: "Quantize",
    qrgba: "Quantize RGBA component",
    qycca: "Quantize YCCA component",
    qy: "Quantize Y",
    inv: "Invert",
    invrgba: "Invert RGBA component",
    inva: "Invert lpha",
    invycc: "Invert YCC component",
    gs: "Gray Scale",
    rasp: "BitRasp",
};
const operatorNames = {
    cmp: "Compose",
    src: "Replace",
    add: "Add",
    addrgbm: "Add RGB modulo 256",
    addyccm: "Add YCC modulo 256",
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
export default class SettingsPage {
    constructor(gltihc) {
        var _a;
        this.gltihc = gltihc;
        const tpl = document.getElementById("settings-tpl");
        this.content = tpl.content.cloneNode(true);
        this.backBtn = this.content.getElementById("settings-back-btn");
        // Cache values as a maps
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
        // Basic settings
        let frag = document.createDocumentFragment();
        optionControls.forEach((o, optIndex) => {
            const id = `input-option-${optIndex}`;
            const el = document.createElement("div");
            el.className = "row responsive";
            el.appendChild((() => {
                const el = document.createElement("label");
                el.setAttribute("for", id);
                el.innerHTML = o.label;
                return el;
            })());
            el.appendChild((() => {
                const el = document.createElement("input");
                el.type = "number";
                el.id = id;
                el.min = String(o.min);
                if (o.max) {
                    el.max = String(o.max);
                }
                if (o.step) {
                    el.step = String(o.step);
                }
                el.value = String(this.gltihc.options[o.prop]);
                el.addEventListener("change", (ev) => {
                    var _a, _b;
                    const val = Number(ev.target.value);
                    this.gltihc.options[o.prop] = val;
                    (_b = (_a = this).onChange) === null || _b === void 0 ? void 0 : _b.call(_a, o.prop, val);
                });
                return el;
            })());
            frag.appendChild(el);
        });
        (_a = this.content.getElementById("settings-basic-form")) === null || _a === void 0 ? void 0 : _a.appendChild(frag);
        const toggleOptions = [
            {
                labels: filterNames,
                values: this.filters,
                opt: "filters",
                parent: this.content.getElementById("settings-filters-form"),
            },
            {
                labels: operatorNames,
                values: this.operators,
                opt: "ops",
                parent: this.content.getElementById("settings-operators-form"),
            },
        ];
        toggleOptions.forEach((to) => {
            var _a;
            frag = document.createDocumentFragment();
            Object.keys(to.labels).forEach((key, optIndex) => {
                const id = `input-toggle-${to.opt}-${optIndex}`;
                const el = document.createElement("div");
                el.className = "row responsive-spread";
                el.appendChild((() => {
                    const el = document.createElement("label");
                    el.setAttribute("for", id);
                    el.innerHTML = to.labels[key];
                    return el;
                })());
                el.appendChild((() => {
                    const el = document.createElement("input");
                    el.type = "checkbox";
                    el.id = id;
                    el.checked = to.values[key] || false;
                    el.addEventListener("change", (ev) => {
                        var _a, _b;
                        to.values[key] = ev.target.checked;
                        const val = Object.keys(to.values).filter((key) => to.values[key]);
                        this.gltihc.options[to.opt] = val;
                        (_b = (_a = this).onChange) === null || _b === void 0 ? void 0 : _b.call(_a, to.opt, val);
                    });
                    return el;
                })());
                frag.appendChild(el);
            });
            (_a = to.parent) === null || _a === void 0 ? void 0 : _a.appendChild(frag);
        });
    }
}
//# sourceMappingURL=settings_page.js.map