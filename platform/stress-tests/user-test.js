import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 10,  // Number of virtual users
  duration: '30s',  // Test duration
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.01'],    // Less than 1% of requests should fail
  },
};

export default function () {
  // Login to get token
  const loginPayload = JSON.stringify({
    user_id: '000018b0e1a211ef95a30242ac180002',
    pin: '123456',
  });

  const loginHeaders = {
    'Content-Type': 'application/json',
  };

  const loginRes = http.post('http://localhost:8080/api/v1/auth/verify-pin', loginPayload, { headers: loginHeaders });

  check(loginRes, {
    'login successful': (r) => r.status === 200,
  });

  // Extract token from login response
  const authToken = JSON.parse(loginRes.body).tokens.access;

  // Define headers with authorization
  const headers = {
    'Authorization': `Bearer ${authToken}`,
    'Content-Type': 'application/json',
  };

  // Test user profile endpoint with authorization
  const profileResponse = http.get('http://localhost:8080/api/v1/user/profile', { headers });
  check(profileResponse, {
    'profile status is 200': (r) => r.status === 200,
    'profile response time < 200ms': (r) => r.timings.duration < 200,
  });

  const greetingResponse = http.get('http://localhost:8080/api/v1/user/greeting', { headers });
  check(greetingResponse, {
    'greeting status is 200': (r) => r.status === 200,
    'greeting response time < 200ms': (r) => r.timings.duration < 200,
  });

  const greetingPayload = JSON.stringify({
    message: 'Hello World',
  })
  const updateGreetingResponse = http.put('http://localhost:8080/api/v1/user/greeting', greetingPayload, { headers });
  check(updateGreetingResponse, {
    'update greeting status is 200': (r) => r.status === 200,
    'update greeting response time < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(1);
}