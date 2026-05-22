import { readFile } from "node:fs/promises";
export class LcpClient {
    baseUrl;
    token;
    twoFactorCode;
    fetchImpl;
    constructor(options = {}) {
        this.baseUrl = (options.baseUrl ?? "http://127.0.0.1:8080").replace(/\/+$/, "");
        this.token = options.token;
        this.twoFactorCode = options.twoFactorCode;
        this.fetchImpl = options.fetchImpl ?? fetch;
    }
    setToken(token) {
        this.token = token;
    }
    async login(username, password, twoFactor) {
        const response = await this.request("/api/v1/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username,
                password,
                ...(twoFactor ? { twoFactor } : {})
            })
        });
        this.token = response.token;
        return response;
    }
    health() {
        return this.request("/healthz");
    }
    ready() {
        return this.request("/readyz");
    }
    status() {
        return this.request("/api/v1/lcp/status", {
            headers: this.authHeaders()
        });
    }
    async uploadPublicationFromFile(filePath, title) {
        const buffer = await readFile(filePath);
        return this.uploadPublication(buffer, title ?? fileStem(filePath));
    }
    async uploadPublication(file, title) {
        const payload = {
            query: "mutation UploadPublication($title: String!, $file: Upload!) { uploadPublication(title: $title, file: $file) { id title downloadURL } }",
            variables: {
                title,
                file: Buffer.from(file).toString("base64")
            }
        };
        return this.graphql("uploadPublication", payload);
    }
    async createLicense(input) {
        const payload = {
            query: "mutation CreateLicense($publicationID: ID!, $userID: ID!, $passphrase: String!, $hint: String!) { createLicense(publicationID: $publicationID, userID: $userID, passphrase: $passphrase, hint: $hint) { id publicationID userID publicationURL passphrase hint } }",
            variables: input
        };
        return this.graphql("createLicense", payload);
    }
    async revokeLicense(licenseID) {
        await this.request(`/api/v1/admin/licenses/${licenseID}/revoke`, {
            method: "POST",
            headers: this.adminHeaders()
        });
        return { status: "revoked", licenseID };
    }
    listAdminLicenses() {
        return this.request("/api/v1/admin/licenses", {
            headers: this.adminHeaders()
        });
    }
    listAudit(limit) {
        const query = typeof limit === "number" ? `?limit=${limit}` : "";
        return this.request(`/api/v1/admin/audit${query}`, {
            headers: this.adminHeaders()
        });
    }
    async downloadLcpl(licenseID) {
        return this.requestText(`/api/v1/licenses/${licenseID}/lcpl`, {
            headers: this.authHeaders()
        });
    }
    async getLicenseStatusDocument(licenseID) {
        return this.request(`/licenses/${licenseID}/status`);
    }
    async graphql(field, payload) {
        const response = await this.request("/graphql", {
            method: "POST",
            headers: {
                ...this.authHeaders(),
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payload)
        });
        if (response.errors?.length) {
            throw new Error(response.errors[0].message);
        }
        if (!response.data || !(field in response.data)) {
            throw new Error(`GraphQL response missing ${field}`);
        }
        return response.data[field];
    }
    authHeaders() {
        if (!this.token) {
            throw new Error("Client token is not set. Call login() first or provide a token.");
        }
        return {
            Authorization: `Bearer ${this.token}`
        };
    }
    adminHeaders() {
        return {
            ...this.authHeaders(),
            ...(this.twoFactorCode ? { "X-2FA-Code": this.twoFactorCode } : {})
        };
    }
    async request(path, init = {}) {
        const response = await this.fetchImpl(`${this.baseUrl}${path}`, init);
        if (!response.ok) {
            throw new Error(`${path} returned ${response.status} ${response.statusText}: ${await response.text()}`);
        }
        return (await response.json());
    }
    async requestText(path, init = {}) {
        const response = await this.fetchImpl(`${this.baseUrl}${path}`, init);
        if (!response.ok) {
            throw new Error(`${path} returned ${response.status} ${response.statusText}: ${await response.text()}`);
        }
        return response.text();
    }
}
function fileStem(path) {
    const normalized = path.replace(/\\/g, "/");
    const lastSegment = normalized.split("/").pop() ?? normalized;
    const suffix = lastSegment.lastIndexOf(".");
    return suffix > 0 ? lastSegment.slice(0, suffix) : lastSegment;
}
