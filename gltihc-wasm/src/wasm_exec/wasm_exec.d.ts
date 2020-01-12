export interface Imports {
    [key: string]: (...args: any[]) => any;
}

declare global {
    declare class Go {
        constructor();

        argv: string[];
        env: { [key: string]: string; };
        exit: (code: number) => void;
        importObject: { [namespace: string]: Imports; };

        run(instance: WebAssembly.Instance): Promise<any>;
    }
}