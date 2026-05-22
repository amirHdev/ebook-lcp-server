import { LcpClient } from "../src/index.js";
const baseUrl = process.env.LCP_BASE_URL ?? "http://127.0.0.1:8080";
const username = process.env.LCP_USERNAME ?? "admin";
const password = process.env.LCP_PASSWORD ?? "admin";
const twoFactor = process.env.LCP_2FA_CODE ?? "123456";
const bookPath = process.env.BOOK_PATH ?? "examples/pride-and-prejudice/pride-and-prejudice.epub";
const client = new LcpClient({
    baseUrl,
    twoFactorCode: twoFactor
});
const login = await client.login(username, password, twoFactor);
const health = await client.health();
const ready = await client.ready();
const publication = await client.uploadPublicationFromFile(bookPath, "SDK Smoke Book");
const license = await client.createLicense({
    publicationID: publication.id,
    userID: "sdk-reader-01",
    passphrase: "open-sesame",
    hint: "demo"
});
const lcpl = await client.downloadLcpl(license.id);
const statusDocument = await client.getLicenseStatusDocument(license.id);
const licenses = await client.listAdminLicenses();
console.log(JSON.stringify({
    loginRole: login.role,
    health,
    ready,
    publicationID: publication.id,
    licenseID: license.id,
    lcplBytes: lcpl.length,
    status: statusDocument.status,
    adminLicenseCount: licenses.licenses.length
}, null, 2));
