# TypeScript SDK

Small dependency-light TypeScript client for the local LCP platform API.

## Build

```bash
./frontend/node_modules/.bin/tsc -p sdk/typescript/tsconfig.json
```

## Smoke test

```bash
node sdk/typescript/dist/examples/smoke.js
```

## Example

```ts
import { LcpClient } from "./src/index.js";

const client = new LcpClient({
  baseUrl: "http://127.0.0.1:8080",
  twoFactorCode: "123456"
});

await client.login("admin", "admin", "123456");
const publication = await client.uploadPublicationFromFile(
  "examples/pride-and-prejudice/pride-and-prejudice.epub",
  "Pride and Prejudice"
);
const license = await client.createLicense({
  publicationID: publication.id,
  userID: "reader-01",
  passphrase: "open-sesame",
  hint: "demo"
});
```
