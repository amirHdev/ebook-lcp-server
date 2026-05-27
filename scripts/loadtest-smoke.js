import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: Number(__ENV.VUS || 100),
  duration: __ENV.DURATION || "30s",
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<200"],
  },
};

const baseURL = __ENV.BASE_URL || "http://localhost:8080";
const token = __ENV.TOKEN || "";
const requestPause = Number(__ENV.REQUEST_PAUSE_SECONDS || 0.1);

export default function () {
  const res = http.get(`${baseURL}/api/v1/lcp/status`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  check(res, {
    "status is 200": (r) => r.status === 200,
  });
  sleep(requestPause);
}
