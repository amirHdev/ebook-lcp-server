# Roadmap

## Active roadmap

| Priority | Item                                                        | Status                 | Why it matters                                                                                                                                                  |
| -------- | ----------------------------------------------------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| P4       | Run official EDRLab validation with production certificates | Ready for external run | The repo-owned packet and evidence workflow are complete; the remaining work is the real external certification run.                                            |
| P7       | Maintain public project polish and community growth         | In progress            | Screenshots, demo media, docs polish, and SDK examples are now in repo; the remaining work is release cadence, fast issue response, and steady community posts. |

## Notes

- Deployment guides are complete under `docs/deploy.md`, `docs/deploy-flyio.md`, and `docs/deploy-railway.md`.
- Kubernetes deployment, HPA, NetworkPolicy, metrics, backup/restore, license lifecycle, and 100-VU load behavior were validated on local Minikube on 2026-05-27; production sign-off remains external.
- The immediate product focus is adoption first, then operator ergonomics, then certification, broader client support, and visible community momentum.
