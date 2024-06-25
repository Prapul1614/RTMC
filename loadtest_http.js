import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  vus: 1,
  duration: '10s',
};

export default function () {
  let url = 'http://localhost:8000/classify';
  let authToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE4OTU1MTksInN1YiI6IjY2NzAxOGFhMDc4YzIwODc0ODQwNWRlYSJ9.HSnkKn--yv4FiIjjWQw6zajpU5tcGfuLOsHjdqvwtE4'; // Replace with your actual auth token
  let payload = randomString(25);

  let params = {
    headers: {
      'Content-Type': 'text/plain', 
      'Authorization': `Bearer ${authToken}`, 
    },
  };

  let res = http.post(url, payload, params);

  check(res, {
    'is status 200': (r) => r.status === 200,
  });

}
