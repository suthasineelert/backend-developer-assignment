import http from 'k6/http';
import { sleep, check } from 'k6';

// Define test configurations
export let options = {
    stages: [
        { duration: '30s', target: 50 },   // Ramp up to 50 users in 30s
        { duration: '1m', target: 100 },   // Stay at 100 users for 1 minute
        { duration: '30s', target: 0 },    // Ramp down to 0 users
    ],
};
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM1ODEwNDQsImlkIjoiMDAwMDE4YjBlMWEyMTFlZjk1YTMwMjQyYWMxODAwMDIifQ.ipWJQDKsT5xmdPfwcRgogmvQ1nfbbYeQ8RnOeIFBY08
// Define base URL and authentication token (modify as needed)
const BASE_URL = 'http://localhost:8080'; // Change to your API URL

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

    let params = {
        headers: {
            Authorization: `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
    };

    let res = http.get(`${BASE_URL}/transactions?page=1`, params);

    // Check response status
    check(res, {
        'Response status is 200': (r) => r.status === 200,
        'Response time is < 500ms': (r) => r.timings.duration < 500,
    });

    sleep(1);
}