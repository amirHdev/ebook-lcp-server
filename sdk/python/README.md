# Python SDK

Small dependency-light Python client for the local LCP platform API.

## Smoke test

```bash
PYTHONPATH=sdk/python python3 sdk/python/examples/smoke.py
```

## Example

```python
from lcp_sdk import LcpClient

client = LcpClient(base_url="http://127.0.0.1:8080", two_factor_code="123456")
client.login("admin", "admin", "123456")
publication = client.upload_publication_from_file(
    "examples/pride-and-prejudice/pride-and-prejudice.epub",
    "Pride and Prejudice",
)
license_doc = client.create_license(
    publication["id"],
    "reader-01",
    "open-sesame",
    "demo",
)
```
