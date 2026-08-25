import assert from 'node:assert/strict';
import test from 'node:test';
import { proxyWellKnownMetadata } from '../../functions/.well-known/minekube-connect/[[path]].js';

test('proxies the Bedrock principal v2 metadata from the watch', async () => {
	const upstreamBody = JSON.stringify({
		issuer: 'minekube-connect',
		trust_domain: 'urn:minekube:connect:production',
		keys: [{ kid: 'v2-2026-08' }],
	});
	const originalFetch = globalThis.fetch;
	let called;
	globalThis.fetch = async (url, init) => {
		called = { url: String(url), method: init?.method ?? 'GET' };
		return new Response(upstreamBody, {
			status: 200,
			headers: {
				'content-type': 'application/json',
				'cache-control': 'public, max-age=300',
			},
		});
	};
	try {
		const res = await proxyWellKnownMetadata(
			new Request(
				'https://connect.minekube.com/.well-known/minekube-connect/bedrock-principal-v2.json',
				{ method: 'GET' },
			),
			{ path: ['bedrock-principal-v2.json'] },
		);
		assert.equal(res.status, 200);
		assert.equal(
			called.url,
			'https://watch-connect.minekube.net/.well-known/minekube-connect/bedrock-principal-v2.json',
		);
		assert.equal(called.method, 'GET');
		assert.equal(res.headers.get('access-control-allow-origin'), '*');
		assert.equal(res.headers.get('cache-control'), 'public, max-age=300');
		assert.equal(await res.text(), upstreamBody);
	} finally {
		globalThis.fetch = originalFetch;
	}
});

test('passes through watch 404s unchanged (metadata not configured)', async () => {
	const originalFetch = globalThis.fetch;
	globalThis.fetch = async () => new Response('not found', { status: 404 });
	try {
		const res = await proxyWellKnownMetadata(
			new Request(
				'https://connect.minekube.com/.well-known/minekube-connect/bedrock-principal-v2.json',
			),
			{ path: ['bedrock-principal-v2.json'] },
		);
		assert.equal(res.status, 404);
		assert.equal(res.headers.get('access-control-allow-origin'), '*');
	} finally {
		globalThis.fetch = originalFetch;
	}
});
