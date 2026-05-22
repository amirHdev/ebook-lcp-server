import { readFile } from "node:fs/promises";

export interface LoginResponse {
  token: string;
  role?: string;
  subject?: string;
  expiresAt?: string;
}

export interface HealthResponse {
  status: string;
}

export interface StatusEnvelope {
  status: string;
  uptimeSec?: number;
  processes?: unknown[];
}

export interface PublicationResult {
  id: string;
  title?: string;
  downloadURL?: string;
}

export interface LicenseResult {
  id: string;
  publicationID: string;
  userID: string;
  publicationURL?: string;
  passphrase?: string;
  hint?: string;
}

export interface AuditEntry {
  id: string;
  action: string;
  actor?: string;
  resource?: string;
  resourceId?: string;
  createdAt?: string;
}

export interface AdminLicense {
  id: string;
  publicationID: string;
  userID: string;
  status?: string;
  publicationURL?: string;
  endDate?: string;
  createdAt?: string;
}

export interface ClientOptions {
  baseUrl?: string;
  token?: string;
  twoFactorCode?: string;
  fetchImpl?: typeof fetch;
}

export class LcpClient {
  private readonly baseUrl: string;
  private token?: string;
  private readonly twoFactorCode?: string;
  private readonly fetchImpl: typeof fetch;

  constructor(options: ClientOptions = {}) {
    this.baseUrl = (options.baseUrl ?? "http://127.0.0.1:8080").replace(/\/+$/, "");
    this.token = options.token;
    this.twoFactorCode = options.twoFactorCode;
    this.fetchImpl = options.fetchImpl ?? fetch;
  }

  setToken(token: string): void {
    this.token = token;
  }

  async login(username: string, password: string, twoFactor?: string): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>("/api/v1/auth/login", {
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

  health(): Promise<HealthResponse> {
    return this.request<HealthResponse>("/healthz");
  }

  ready(): Promise<HealthResponse> {
    return this.request<HealthResponse>("/readyz");
  }

  status(): Promise<StatusEnvelope> {
    return this.request<StatusEnvelope>("/api/v1/lcp/status", {
      headers: this.authHeaders()
    });
  }

  async uploadPublicationFromFile(filePath: string, title?: string): Promise<PublicationResult> {
    const buffer = await readFile(filePath);
    return this.uploadPublication(buffer, title ?? fileStem(filePath));
  }

  async uploadPublication(file: Uint8Array, title: string): Promise<PublicationResult> {
    const payload = {
      query:
        "mutation UploadPublication($title: String!, $file: Upload!) { uploadPublication(title: $title, file: $file) { id title downloadURL } }",
      variables: {
        title,
        file: Buffer.from(file).toString("base64")
      }
    };

    return this.graphql<PublicationResult>("uploadPublication", payload);
  }

  async createLicense(input: {
    publicationID: string;
    userID: string;
    passphrase: string;
    hint: string;
  }): Promise<LicenseResult> {
    const payload = {
      query:
        "mutation CreateLicense($publicationID: ID!, $userID: ID!, $passphrase: String!, $hint: String!) { createLicense(publicationID: $publicationID, userID: $userID, passphrase: $passphrase, hint: $hint) { id publicationID userID publicationURL passphrase hint } }",
      variables: input
    };

    return this.graphql<LicenseResult>("createLicense", payload);
  }

  async revokeLicense(licenseID: string): Promise<{ status: string; licenseID: string }> {
    await this.request<Record<string, unknown>>(`/api/v1/admin/licenses/${licenseID}/revoke`, {
      method: "POST",
      headers: this.adminHeaders()
    });
    return { status: "revoked", licenseID };
  }

  listAdminLicenses(): Promise<{ licenses: AdminLicense[] }> {
    return this.request<{ licenses: AdminLicense[] }>("/api/v1/admin/licenses", {
      headers: this.adminHeaders()
    });
  }

  listAudit(limit?: number): Promise<{ entries: AuditEntry[] }> {
    const query = typeof limit === "number" ? `?limit=${limit}` : "";
    return this.request<{ entries: AuditEntry[] }>(`/api/v1/admin/audit${query}`, {
      headers: this.adminHeaders()
    });
  }

  async downloadLcpl(licenseID: string): Promise<string> {
    return this.requestText(`/api/v1/licenses/${licenseID}/lcpl`, {
      headers: this.authHeaders()
    });
  }

  async getLicenseStatusDocument(licenseID: string): Promise<Record<string, unknown>> {
    return this.request<Record<string, unknown>>(`/licenses/${licenseID}/status`);
  }

  private async graphql<T>(field: string, payload: Record<string, unknown>): Promise<T> {
    const response = await this.request<{ data?: Record<string, T>; errors?: Array<{ message: string }> }>(
      "/graphql",
      {
        method: "POST",
        headers: {
          ...this.authHeaders(),
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      }
    );

    if (response.errors?.length) {
      throw new Error(response.errors[0].message);
    }

    if (!response.data || !(field in response.data)) {
      throw new Error(`GraphQL response missing ${field}`);
    }

    return response.data[field];
  }

  private authHeaders(): HeadersInit {
    if (!this.token) {
      throw new Error("Client token is not set. Call login() first or provide a token.");
    }
    return {
      Authorization: `Bearer ${this.token}`
    };
  }

  private adminHeaders(): HeadersInit {
    return {
      ...this.authHeaders(),
      ...(this.twoFactorCode ? { "X-2FA-Code": this.twoFactorCode } : {})
    };
  }

  private async request<T>(path: string, init: RequestInit = {}): Promise<T> {
    const response = await this.fetchImpl(`${this.baseUrl}${path}`, init);
    if (!response.ok) {
      throw new Error(`${path} returned ${response.status} ${response.statusText}: ${await response.text()}`);
    }
    return (await response.json()) as T;
  }

  private async requestText(path: string, init: RequestInit = {}): Promise<string> {
    const response = await this.fetchImpl(`${this.baseUrl}${path}`, init);
    if (!response.ok) {
      throw new Error(`${path} returned ${response.status} ${response.statusText}: ${await response.text()}`);
    }
    return response.text();
  }
}

function fileStem(path: string): string {
  const normalized = path.replace(/\\/g, "/");
  const lastSegment = normalized.split("/").pop() ?? normalized;
  const suffix = lastSegment.lastIndexOf(".");
  return suffix > 0 ? lastSegment.slice(0, suffix) : lastSegment;
}
