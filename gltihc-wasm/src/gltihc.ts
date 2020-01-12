import "./wasm_exec/wasm_exec.js";

interface Options {
    minIterations: number;
    maxIterations: number;
    blockSize: number;
    minSegmentSize: number;
    maxSegmentSize: number;
    minFilters: number;
    maxFilters: number;
    filters?: string[];
    ops?: string[];
}

type ProcessImageFunc = (src: Uint8Array, opt: Options) => Uint8Array | string;

declare global {
    interface Window {
        _gltihcInitDone: () => void;
        _gltihcProcessImage: ProcessImageFunc;
    }
}

const assemblyPath = "./gltihc.wasm";

export class Gltihc {
    options: Options = {
        minIterations: 10,
        maxIterations: 10,
        blockSize: 16,
        minSegmentSize: 0.01,
        maxSegmentSize: 0.2,
        minFilters: 1,
        maxFilters: 4,
    };
    readonly initDone: Promise<any>;
    source?: Blob;

    private gltihcProcessImage?: ProcessImageFunc;

    constructor() {
        this.initDone = new Promise<any>((resolve) => {
            window._gltihcInitDone = () => {
                this.gltihcProcessImage = window._gltihcProcessImage;
                resolve();
            }
            const go = new Go();
            WebAssembly.instantiateStreaming(fetch(assemblyPath), go.importObject).then((result) => {
                go.run(result.instance);
            });
        })
    }

    processImage(): Promise<Blob> {
        return new Promise<Blob>((resolve, reject) => {
            if (!this.source) {
                reject("No input data");
                return;
            }
            let reader = new FileReader();
            reader.onload = (ev) => {
                if (!(ev.target?.result instanceof ArrayBuffer)) {
                    reject();
                    return;
                }
                let bytes = new Uint8Array(ev.target.result);
                let result = this.gltihcProcessImage?.(bytes, this.options)
                if (result instanceof Uint8Array) {
                    resolve(new Blob([result], { type: "image/png" }));
                } else {
                    reject(result);
                }
            };
            reader.readAsArrayBuffer(this.source);
        })
    }
}