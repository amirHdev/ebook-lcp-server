declare const process: {
  env: Record<string, string | undefined>;
};

declare const Buffer: {
  from(data: Uint8Array | ArrayBuffer | string, encoding?: string): {
    toString(encoding?: string): string;
  };
};

declare module "node:fs/promises" {
  export function readFile(path: string): Promise<Uint8Array>;
}
