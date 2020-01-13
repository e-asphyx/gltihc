import "./wasm_exec/wasm_exec.js";

export type Option = "minIterations" | "maxIterations" | "blockSize" | "minSegmentSize" |
    "maxSegmentSize" | "minFilters" | "maxFilters" | "filters" | "ops" | "maxWidth" | "maxHeight";
;

export type Options = {
    [prop in Option]: number | string[] | null;
};

type ProcessImageFunc = (src: Uint8Array, opt: Options) => Uint8Array | string;

declare global {
    // tslint:disable-next-line: interface-name
    interface Window {
        _gltihcInitDone: () => void;
        _gltihcProcessImage: ProcessImageFunc;
    }
}

const assemblyPath = "./gltihc.wasm";

export class Gltihc {
    public options: Options = {
        blockSize: 16,
        filters: null,
        maxFilters: 4,
        maxIterations: 10,
        maxSegmentSize: 0.2,
        minFilters: 1,
        minIterations: 10,
        minSegmentSize: 0.01,
        ops: null,
        maxHeight: 1024,
        maxWidth: 1024,
    };

    public readonly initDone: Promise<any>;
    public source?: Blob;

    private gltihcProcessImage?: ProcessImageFunc;

    constructor() {
        this.initDone = new Promise<any>((resolve) => {
            window._gltihcInitDone = () => {
                this.gltihcProcessImage = window._gltihcProcessImage;
                resolve();
            };
            const go = new Go();
            WebAssembly.instantiateStreaming(fetch(assemblyPath), go.importObject).then((result) => {
                go.run(result.instance);
            });
        });
    }

    public processImage(): Promise<Blob> {
        return new Promise<Blob>((resolve, reject) => {
            if (!this.source) {
                reject("No input data");
                return;
            }
            const reader = new FileReader();
            reader.onload = (ev) => {
                if (!(ev.target?.result instanceof ArrayBuffer)) {
                    reject();
                    return;
                }
                const bytes = new Uint8Array(ev.target.result);
                const result = this.gltihcProcessImage?.(bytes, this.options);
                if (result instanceof Uint8Array) {
                    resolve(new Blob([result], { type: "image/jpeg" }));
                } else {
                    reject(result);
                }
            };
            reader.readAsArrayBuffer(this.source);
        });
    }
}
