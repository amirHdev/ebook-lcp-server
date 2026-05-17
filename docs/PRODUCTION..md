# Production Deployment Guide

This document outlines best practices for deploying **ebook-lcp-server** in a production environment.

---

## 1. Security Hardening

### Transport Security

- Always place a reverse proxy (**Caddy**, **Traefik**, or **Nginx**) in front of the application.
- Enable automatic TLS using Let's Encrypt.
- Force HTTPS and redirect all HTTP traffic.
- Use TLS 1.3 only and modern cipher suites.

### Authentication

- Enable **2FA** for all admin accounts.
- Use strong, randomly generated JWT secret keys.
- Consider IP whitelisting for the admin interface when possible.
- Rotate credentials and secrets periodically.

### Secrets Management

- Never store sensitive keys in code or Docker images.
- Recommended solutions:
  - Docker Secrets / Kubernetes Secrets
  - External Secrets Operator
  - HashiCorp Vault

### Runtime Security

- Run containers as non-root user (already done with distroless images).
- Use read-only filesystem where possible.
- Drop all unnecessary capabilities.
- Set strict resource limits (CPU, Memory, PIDs).

---

## 2. Recommended Production Architecture

````text
Internet
   │
   ▼
Caddy/Traefik (TLS + Rate Limiting)
   │
   ▼
ebook-lcp-server (2-3 replicas)
   │
   ▼
PostgreSQL + Redis (optional)
   │
   ▼
S3 / MinIO Storage

## 3. Database & Storage Configuration

- **Database**: Use PostgreSQL instead of JSON/in-memory storage.
- **Storage**: Use S3-compatible storage (MinIO, AWS S3, Backblaze B2, etc.).
- Enable SSL/TLS on database connections.
- Set up regular automated backups.

---

## 4. Monitoring & Observability

**Included in the project:**
- Prometheus metrics
- Grafana dashboards

**Recommended additions:**
- Loki + Promtail for centralized logging
- Alertmanager for notifications
- Uptime Kuma or similar for external monitoring

---

## 5. High Availability & Scaling

- Run multiple replicas of the LCP server.
- Use managed PostgreSQL or replication.
- Configure proper health checks and readiness probes.
- Use Redis for session/distributed locking if scaling horizontally.

---

## 6. Backup Strategy

- Daily PostgreSQL backups
- Regular snapshots of storage bucket
- Secure backup of encryption keys and certificates
- Periodically test restore procedures

---

## 7. Logging

- Use structured (JSON) logging in production.
- Set log level to `info` or higher.
- Ship logs to Loki or an external service.

---

## 8. Maintenance & Updates

- Keep dependencies updated (Dependabot is already configured).
- Regularly scan for vulnerabilities using `trivy`.
- Have a tested rollback plan.

---

## 9. Pre-Production Checklist

- [ ] TLS/HTTPS configured
- [ ] 2FA enabled on admin accounts
- [ ] PostgreSQL + S3 in use
- [ ] Secrets managed externally
- [ ] Resource limits defined
- [ ] Monitoring & alerting active
- [ ] Backups configured and tested
- [ ] Rate limiting enabled
- [ ] No test/demo data present

---

## 10. Useful Commands

```bash
# Vulnerability scan
trivy fs .

# Build production image
docker build -t ebook-lcp-server:prod -f Dockerfile .

# Run production compose
docker compose -f docker-compose.prod.yml up -d


✅ Features fixed in this version:
1. All headings now consistent (`##`).
2. Lists properly formatted with `-`.
3. Checklists correctly use `[ ]`.
4. Code blocks are marked with language hints (`bash`).
5. Sections separated with `---` for readability.

If you want, I can **merge this with your earlier sections** and produce a **complete, fully Markdown-ready Production Deployment Guide** in one file—it will be ready to paste into GitHub or internal docs.

Do you want me to do that?
````
