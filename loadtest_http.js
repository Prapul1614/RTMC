import http from 'k6/http';
import { sleep, check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
  vus: 1,
  duration: '1s',
};

export default function () {
  let url = 'http://localhost:8000/classify';
  let authToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg5ODg3OTYsInN1YiI6IjY2NzAxOGFhMDc4YzIwODc0ODQwNWRlYSJ9.-BYVQ6KMNOAlEtWk9K5vPTDiSfUkNY6kSSCkBhCre8c'; // Replace with your actual auth token
  let payload = randomString(25); // Generates a random string of 10 characters

  let params = {
    headers: {
      'Content-Type': 'text/plain', // Set the Content-Type to text/plain
      'Authorization': `Bearer ${authToken}`, // Add the auth token in the header
    },
  };

  let res = http.post(url, payload, params);

  check(res, {
    'is status 200': (r) => r.status === 200,
  });

}
