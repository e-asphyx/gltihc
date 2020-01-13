import "./wasm_exec/wasm_exec.js";
const assemblyPath = "./gltihc.wasm";
export class Gltihc {
    constructor() {
        this.options = {
            blockSize: 16,
            filters: null,
            maxFilters: 4,
            maxIterations: 10,
            maxSegmentSize: 0.2,
            minFilters: 1,
            minIterations: 10,
            minSegmentSize: 0.01,
            ops: null,
        };
        this.initDone = new Promise((resolve) => {
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
    processImage() {
        return new Promise((resolve, reject) => {
            if (!this.source) {
                reject("No input data");
                return;
            }
            const reader = new FileReader();
            reader.onload = (ev) => {
                var _a, _b, _c;
                if (!(((_a = ev.target) === null || _a === void 0 ? void 0 : _a.result) instanceof ArrayBuffer)) {
                    reject();
                    return;
                }
                const bytes = new Uint8Array(ev.target.result);
                const result = (_c = (_b = this).gltihcProcessImage) === null || _c === void 0 ? void 0 : _c.call(_b, bytes, this.options);
                if (result instanceof Uint8Array) {
                    resolve(new Blob([result], { type: "image/jpeg" }));
                }
                else {
                    reject(result);
                }
            };
            reader.readAsArrayBuffer(this.source);
        });
    }
}
//# sourceMappingURL=gltihc.js.map