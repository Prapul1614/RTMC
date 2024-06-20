import grpc from 'k6/net/grpc';
import { Client, Stream } from 'k6/experimental/grpc'
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

// Load tokens from tokens.json file
const tokens = new SharedArray('tokens', function () {
    return JSON.parse(open('./tokens.json'));
});

let tokenIndex = 0;

const texts = new SharedArray('texts', function () {
    return [
      "transaction transaction location",
      "urgent urgent money transfer",
      "stock price now decrease",
      "three failes login attempts",
      "delay traffic delay delay",
      "you are mentions in and good for mentions",
      "energy high very very",
      "your request for product recall successfull",
      "Yoo its going to be an earthquake",
      "hospital patients waiting patients no rooms for patients",
    ];
});

const client = new Client();
client.load([], './rule.proto');

// https://www.youtube.com/watch?v=ghuo8m7AXEM
// https://github.com/grafana/k6/blob/master/examples/grpc_client_streaming.js
// https://github.com/grafana/xk6
// https://github.com/grafana/xk6-grpc
export let options = {
    stages: [
        { duration: '30s', target: 10 }, // ramp-up to 10 users
        { duration: '1m', target: 10 }, // stay at 10 users
        { duration: '30s', target: 0 },  // ramp-down to 0 users
    ]
};


export default () => {
    const token = tokens[tokenIndex];
    tokenIndex = (tokenIndex + 1) % tokens.length;

    client.connect('localhost:3000', { plaintext: true });

    const metadata = { authorization: 'Bearer ' + token };
    const callOptions = {
        metadata,
    };
    
    const stream = new Stream(client, 'rulepb.RuleService/StreamData', callOptions)

    console.log(tokenIndex)

    stream.on('data', (message) => {
        //console.log('Received message: ', message);
    });

    stream.on('end', () => {
        console.log('Stream ended');
    });

    stream.on('error', (err) => {
        //console.error('Stream error: ', err);
    });

    for (let i = 0; i < 100; i++) {
        const data = { text: texts[randomIntBetween(0, texts.length - 1)] };
        stream.write(data);
        sleep(0.5); 
    }

    // Close the stream
    stream.end();

    sleep(2)

    client.close();
};
