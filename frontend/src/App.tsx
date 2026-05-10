import React, { useEffect, useMemo, useState } from "react";
import { Activity, BarChart3, CheckCircle2, FileUp, KeyRound, Play, Shield } from "lucide-react";
import { createRoot } from "react-dom/client";
import "./styles.css";

type ProcessStatus = {
  id: string;
  status: string;
  publicationId?: string;
  error?: string;
  createdAt: string;
  updatedAt: string;
};

type StatusResponse = {
  status: string;
  uptimeSec: number;
  processes: ProcessStatus[];
};

type MetricsResponse = {
  uptimeSec: number;
  processes: number;
  metrics: {
    requestsTotal: number;
    processesOk: number;
    processesFail: number;
  };
};

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

function App() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState("");
  const [role, setRole] = useState<string>("");
  const [twoFactor, setTwoFactor] = useState("");
  const [title, setTitle] = useState("Example Publication");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [filePreview, setFilePreview] = useState("Choose a publication file to upload.");
  const [status, setStatus] = useState<StatusResponse | null>(null);
  const [metrics, setMetrics] = useState<MetricsResponse | null>(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    const savedToken = window.localStorage.getItem("lcp-token") || "";
    const savedUser = window.localStorage.getItem("lcp-username") || "";
    const savedRole = window.localStorage.getItem("lcp-role") || "";
    const saved2fa = window.localStorage.getItem("lcp-2fa") || "";
    setToken(savedToken);
    setUsername(savedUser);
    setRole(savedRole);
    setTwoFactor(saved2fa);
  }, []);

  const authHeaders = useMemo(
    () => ({
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json"
    }),
    [token]
  );

  async function login() {
    const response = await fetch(`${API_BASE}/api/v1/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username,
        password,
        twoFactor
      })
    });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || "login failed");
    setToken(body.token);
    setRole(body.role || "");
    window.localStorage.setItem("lcp-token", body.token);
    window.localStorage.setItem("lcp-username", username);
    window.localStorage.setItem("lcp-role", body.role || "");
    window.localStorage.setItem("lcp-2fa", twoFactor);
    setMessage(`Signed in as ${body.subject}`);
  }

  async function refreshStatus() {
    const response = await fetch(`${API_BASE}/api/v1/lcp/status`, { headers: authHeaders });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || "status request failed");
    setStatus(body);
  }

  async function refreshMetrics() {
    const response = await fetch(`${API_BASE}/api/v1/admin/metrics`, {
      headers: { ...authHeaders, "X-2FA-Code": twoFactor }
    });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || "metrics request failed");
    setMetrics(body);
  }

  async function processContent() {
    if (!selectedFile) {
      throw new Error("choose a publication file first");
    }

    const fileBase64 = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        const result = reader.result;
        if (typeof result !== "string") {
          reject(new Error("file read failed"));
          return;
        }
        resolve(result.split(",").pop() || "");
      };
      reader.onerror = () => reject(new Error("file read failed"));
      reader.readAsDataURL(selectedFile);
    });

    const response = await fetch(`${API_BASE}/api/v1/lcp/process`, {
      method: "POST",
      headers: {
        ...authHeaders,
        ...(role === "admin" && twoFactor ? { "X-2FA-Code": twoFactor } : {})
      },
      body: JSON.stringify({ title, file: fileBase64 })
    });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || body.error || "process request failed");
    setMessage(`Process ${body.id} completed`);
    await refreshStatus();
  }

  function onFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0] || null;
    setSelectedFile(file);
    setFilePreview(file ? `${file.name} · ${file.type || "unknown type"} · ${Math.ceil(file.size / 1024)} KiB` : "Choose a publication file to upload.");
  }

  async function run(action: () => Promise<void>) {
    setMessage("");
    try {
      await action();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "request failed");
    }
  }

  return (
    <main className="shell">
      <header className="topbar">
        <div>
          <h1>LCP Admin</h1>
          <p>Operations dashboard for publications, processing, and runtime health.</p>
        </div>
        <div className="status-pill">
          <Activity size={18} />
          {status?.status || "not loaded"}
        </div>
      </header>

      <section className="grid">
        <div className="panel auth-panel">
          <h2><Shield size={18} /> Admin Login</h2>
          <label>
            Username
            <input value={username} onChange={(event) => setUsername(event.target.value)} />
          </label>
          <label>
            Password
            <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
          </label>
          <label>
            Admin 2FA
            <input value={twoFactor} onChange={(event) => setTwoFactor(event.target.value)} />
          </label>
          <button onClick={() => run(login)}>
            <Shield size={18} />
            Sign In
          </button>
          <label>
            JWT
            <textarea readOnly value={token} placeholder="JWT appears here after sign in" />
          </label>
          <div className="file-meta">Role: {role || "unset"}</div>
        </div>

        <div className="panel">
          <h2><Play size={18} /> Process</h2>
          <label>
            Title
            <input value={title} onChange={(event) => setTitle(event.target.value)} />
          </label>
          <label>
            Publication File
            <div className="file-picker">
              <label className="file-button">
                <FileUp size={18} />
                <span>Select file</span>
                <input type="file" onChange={onFileChange} />
              </label>
              <div className="file-meta">{filePreview}</div>
            </div>
          </label>
          <button onClick={() => run(processContent)}>
            <CheckCircle2 size={18} />
            Upload and Process
          </button>
        </div>

        <div className="panel">
          <h2><BarChart3 size={18} /> Metrics</h2>
          <button onClick={() => run(refreshMetrics)}>
            <KeyRound size={18} />
            Load Metrics
          </button>
          <dl className="metrics">
            <dt>Uptime</dt>
            <dd>{metrics?.uptimeSec ?? 0}s</dd>
            <dt>Requests</dt>
            <dd>{metrics?.metrics.requestsTotal ?? 0}</dd>
            <dt>OK / Failed</dt>
            <dd>{metrics ? `${metrics.metrics.processesOk} / ${metrics.metrics.processesFail}` : "0 / 0"}</dd>
          </dl>
        </div>
      </section>

      <section className="panel">
        <div className="section-head">
          <h2><Activity size={18} /> Process Status</h2>
          <button onClick={() => run(refreshStatus)}>Refresh</button>
        </div>
        {message && <p className="message">{message}</p>}
        <div className="table">
          <div className="row header">
            <span>ID</span>
            <span>Status</span>
            <span>Publication</span>
            <span>Updated</span>
          </div>
          {(status?.processes || []).map((item) => (
            <div className="row" key={item.id}>
              <span>{item.id}</span>
              <span>{item.status}</span>
              <span>{item.publicationId || "-"}</span>
              <span>{new Date(item.updatedAt).toLocaleString()}</span>
            </div>
          ))}
        </div>
      </section>
    </main>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
