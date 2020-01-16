import { Gltihc, Option } from "./gltihc.js";

const optionControls: Array<{
    label: string;
    min: number;
    max?: number;
    step?: number;
    prop: Option;
}> = [
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

interface ToggleOptions {
    [name: string]: string;
}

const filterNames: ToggleOptions = {
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

const operatorNames: ToggleOptions = {
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

interface ToggleValues {
    [n: string]: boolean;
}

function defaultToggleValues(opt: ToggleOptions): ToggleValues {
    return Object.assign({}, ...Object.keys(opt).map<ToggleValues>((n) => ({ [n]: true })));
}

function toggleValues(opt: string[]): ToggleValues {
    return Object.assign({}, ...opt.map<ToggleValues>((n) => ({ [n]: true })));
}

export default class SettingsPage {
    public onChange?: (o?: Option, v?: any) => void;

    public readonly content: DocumentFragment;
    public readonly backBtn: HTMLInputElement;

    private filters: ToggleValues;
    private operators: ToggleValues;

    constructor(private gltihc: Gltihc) {
        const tpl = <HTMLTemplateElement>document.getElementById("settings-tpl");

        this.content = <DocumentFragment>tpl.content.cloneNode(true);
        this.backBtn = <HTMLInputElement>this.content.getElementById("settings-back-btn");

        // Cache values as a maps
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
                    const val = Number((<HTMLInputElement>ev.target).value);
                    this.gltihc.options[o.prop] = val;
                    this.onChange?.(o.prop, val);
                });
                return el;
            })());
            frag.appendChild(el);
        });
        this.content.getElementById("settings-basic-form")?.appendChild(frag);

        const toggleOptions: Array<{
            labels: ToggleOptions;
            values: ToggleValues;
            opt: Option;
            parent: Node | null;
        }> = [
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
                        to.values[key] = (<HTMLInputElement>ev.target).checked;
                        const val = Object.keys(to.values).filter((key) => to.values[key]);
                        this.gltihc.options[to.opt] = val;
                        this.onChange?.(to.opt, val);
                    });
                    return el;
                })());
                frag.appendChild(el);
            });
            to.parent?.appendChild(frag);
        });
    }
}
