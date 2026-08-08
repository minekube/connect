import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import test from 'node:test';
import worker from '../worker.mjs';
import { normalizePagesResponse } from '../worker-response.mjs';

const redirectPaths = [
  '/*.html',
  '/guide/*',
  '/changelog/*',
  '/discord',
  '/discord/',
];

const readConfig = async (name) =>
  JSON.parse(await readFile(new URL(`../${name}`, import.meta.url), 'utf8'));

const hasHeaderRule = (headers, path, name, value) => {
  const start = headers.indexOf(`${path}\n`);
  const end = headers.indexOf('\n\n', start);
  const rule = headers.slice(start, end === -1 ? undefined : end);
  return rule.includes(
    value === undefined
      ? `  ${name}`
      : `  ${name}: ${value}`
  );
};

test('production and canary Workers preserve the Connect Pages contract', async () => {
  const production = await readConfig('wrangler.jsonc');
  const canary = await readConfig('wrangler.canary.jsonc');
  const headers = await readFile(
    new URL('../docs/public/_headers', import.meta.url),
    'utf8'
  );
  const redirects = await readFile(
    new URL('../docs/public/_redirects', import.meta.url),
    'utf8'
  );
  const packageJson = JSON.parse(
    await readFile(new URL('../package.json', import.meta.url), 'utf8')
  );
  const shared = {
    compatibility_date: '2026-08-07',
    main: 'worker.mjs',
    assets: {
      directory: 'docs/.vitepress/dist',
      binding: 'ASSETS',
      run_worker_first: redirectPaths,
    },
  };

  assert.equal(packageJson.packageManager, 'yarn@4.17.1');
  assert.equal(packageJson.devDependencies.wrangler, '4.115.0');
  assert.equal(
    packageJson.scripts['test:worker'],
    'node --test scripts/workers-migration.test.mjs'
  );
  assert.equal(packageJson.scripts['deploy:worker'], 'wrangler deploy');
  assert.equal(
    packageJson.scripts['deploy:worker:canary'],
    'wrangler deploy --config wrangler.canary.jsonc'
  );

  for (const config of [production, canary]) {
    assert.equal(config.compatibility_date, shared.compatibility_date);
    assert.equal(config.main, shared.main);
    assert.deepEqual(config.assets, shared.assets);
  }

  assert.equal(production.name, 'connect-minekube-worker');
  assert.equal(production.workers_dev, false);
  assert.equal(production.preview_urls, false);
  assert.deepEqual(production.routes, [
    { pattern: 'connect.minekube.com', custom_domain: true },
  ]);

  assert.equal(canary.name, 'connect-minekube-worker-canary');
  assert.equal(canary.workers_dev, false);
  assert.equal(canary.preview_urls, false);
  assert.deepEqual(canary.routes, [
    { pattern: 'connect-docs-canary.minekube.com', custom_domain: true },
  ]);

  assert.match(
    headers,
    /\/\*\n  Access-Control-Allow-Origin: \*\n  X-Content-Type-Options: nosniff\n  Referrer-Policy: strict-origin-when-cross-origin/
  );

  assert.doesNotMatch(redirects, / 308$/m);

  for (const path of [
    '/',
    '/blog/',
    '/blog/gate-api',
    '/blog/gate-lite',
    '/blog/managed-bedrock-support',
    '/blog/minekube-ai-support-loop',
    '/blog/minekube-browser-apis-agents',
    '/blog/minekube-control-plane',
    '/plans',
    '/team',
  ]) {
    assert.equal(
      hasHeaderRule(headers, path, 'Content-Type', 'text/html; charset=utf-8'),
      true,
      `missing HTML charset header rule for ${path}`
    );
  }

  for (const path of ['/guide/*', '/changelog/*']) {
    assert.equal(
      hasHeaderRule(headers, path, 'Content-Type'),
      false,
      `${path} must not overlap generated Markdown route content types`
    );
  }

  for (const [path, name, value] of [
    ['/*.css', 'Content-Type', 'text/css; charset=utf-8'],
    ['/*.js', 'Content-Type', 'application/javascript'],
    ['/*.jpeg', 'cache-control', 'public, max-age=14400, must-revalidate'],
    [
      '/vp-icons.css',
      'cache-control',
      'public, max-age=14400, must-revalidate',
    ],
  ]) {
    assert.equal(
      hasHeaderRule(headers, path, name, value),
      true,
      `missing ${name} header rule for ${path}`
    );
  }

  assert.equal(
    hasHeaderRule(
      headers,
      '/*.md',
      'Content-Type',
      'text/markdown; charset=utf-8'
    ),
    true,
    'Markdown routes must set one exact Markdown content type'
  );
});

test('the Worker preserves Pages response headers and redirect metadata', async () => {
  const html = normalizePagesResponse(
    new Response('<!doctype html>', {
      headers: { 'content-type': 'text/html' },
    })
  );
  assert.equal(html.headers.get('access-control-allow-origin'), '*');
  assert.equal(html.headers.get('content-type'), 'text/html; charset=utf-8');
  assert.equal(html.headers.get('x-content-type-options'), 'nosniff');
  assert.equal(
    html.headers.get('referrer-policy'),
    'strict-origin-when-cross-origin'
  );

  const encoded = normalizePagesResponse(
    new Response('already encoded', {
      headers: { 'content-type': 'Text/Plain; Charset=UTF-8' },
    })
  );
  assert.equal(
    encoded.headers.get('content-type'),
    'Text/Plain; Charset=UTF-8'
  );

  const missing = normalizePagesResponse(
    new Response('missing', {
      status: 404,
      headers: {
        'cache-control': 'public, max-age=0, must-revalidate',
        'content-type': 'text/html',
      },
    })
  );
  assert.equal(missing.headers.get('cache-control'), 'no-store');

  const redirect = normalizePagesResponse(
    new Response('redirecting', {
      status: 301,
      headers: {
        location: '/changelog/',
        'content-type': 'text/plain',
      },
    })
  );
  assert.equal(redirect.status, 301);
  assert.equal(redirect.headers.get('location'), '/changelog/');
  assert.equal(redirect.headers.get('content-type'), 'text/plain; charset=utf-8');

  const assetRedirect = normalizePagesResponse(
    new Response(null, {
      status: 302,
      headers: { location: 'https://discord.com/invite/6vMDqWE' },
    })
  );
  assert.equal(
    assetRedirect.headers.get('location'),
    'https://discord.com/invite/6vMDqWE'
  );
  assert.equal(
    assetRedirect.headers.get('content-type'),
    'text/plain;charset=UTF-8'
  );
});

test('the Worker recreates Pages redirects before Static Assets routing', async () => {
  const assets = {
    fetch() {
      throw new Error('redirect requests must not reach the asset binding');
    },
  };
  const redirects = [
    ['/discord', 302, 'https://discord.com/invite/6vMDqWE'],
    ['/discord/', 302, 'https://discord.com/invite/6vMDqWE'],
    ['/guide/changelog', 301, '/changelog/'],
    ['/guide/changelog/', 301, '/changelog/'],
    ['/guide/changelog.html', 301, '/changelog/'],
  ];

  for (const [path, status, destination] of redirects) {
    const response = await worker.fetch(
      new Request(`https://connect.minekube.com${path}?source=test`),
      { ASSETS: assets }
    );
    const location = `${destination}?source=test`;
    const body = `Redirecting to ${location}`;

    assert.equal(response.status, status);
    assert.equal(response.headers.get('location'), location);
    assert.equal(
      response.headers.get('content-type'),
      'text/plain;charset=UTF-8'
    );
    assert.equal(
      response.headers.get('content-length'),
      String(new TextEncoder().encode(body).byteLength)
    );
    assert.equal(response.headers.get('access-control-allow-origin'), '*');
    assert.equal(response.headers.get('x-content-type-options'), 'nosniff');
    assert.equal(
      response.headers.get('referrer-policy'),
      'strict-origin-when-cross-origin'
    );
    assert.equal(await response.text(), body);

    const head = await worker.fetch(
      new Request(`https://connect.minekube.com${path}`, { method: 'HEAD' }),
      { ASSETS: assets }
    );
    assert.equal(head.status, status);
    assert.equal(head.headers.get('location'), destination);
    assert.equal(head.headers.get('content-length'), null);
    assert.equal(await head.text(), '');
  }
});

test('the Worker recreates Pages canonical HTML redirects', async () => {
  const assets = {
    fetch() {
      throw new Error('canonical redirects must not reach the asset binding');
    },
  };
  const redirects = [
    ['/index.html', '/'],
    ['/guide/topologies.html', '/guide/topologies'],
    ['/guide/connectors/index.html', '/guide/connectors/'],
    ['/blog/index.html', '/blog/'],
  ];

  for (const [path, destination] of redirects) {
    const response = await worker.fetch(
      new Request(`https://connect.minekube.com${path}?source=test`),
      { ASSETS: assets }
    );

    assert.equal(response.status, 308);
    assert.equal(response.headers.get('location'), `${destination}?source=test`);
    assert.equal(response.headers.get('content-length'), '0');
    assert.equal(response.headers.get('content-type'), null);
    assert.equal(response.headers.get('x-content-type-options'), null);
    assert.equal(response.headers.get('access-control-allow-origin'), '*');
    assert.equal(
      response.headers.get('referrer-policy'),
      'strict-origin-when-cross-origin'
    );
    assert.equal(await response.text(), '');

    const head = await worker.fetch(
      new Request(`https://connect.minekube.com${path}`, { method: 'HEAD' }),
      { ASSETS: assets }
    );
    assert.equal(head.status, 308);
    assert.equal(head.headers.get('location'), destination);
    assert.equal(head.headers.get('content-length'), null);
    assert.equal(await head.text(), '');
  }
});

test('the Worker normalizes HTML and Markdown below clean-URL trees', async () => {
  const assets = {
    fetch(request) {
      const pathname = new URL(request.url).pathname;
      if (pathname.endsWith('.md')) {
        return new Response('# Bedrock', {
          headers: {
            'content-type':
              'text/markdown; charset=utf-8, text/html; charset=utf-8',
          },
        });
      }

      return new Response('<!doctype html>', {
        headers: { 'content-type': 'text/html' },
      });
    },
  };

  for (const [path, type] of [
    ['/guide/bedrock', 'text/html; charset=utf-8'],
    ['/changelog/', 'text/html; charset=utf-8'],
    ['/guide/bedrock.md', 'text/markdown; charset=utf-8'],
  ]) {
    const response = await worker.fetch(
      new Request(`https://connect.minekube.com${path}`),
      { ASSETS: assets }
    );

    assert.equal(response.status, 200);
    assert.equal(response.headers.get('content-type'), type);
  }
});

test('the Worker serves the branded 404 with the Pages cache contract', async () => {
  const requests = [];
  const assets = {
    fetch(request) {
      const pathname = new URL(request.url).pathname;
      requests.push(pathname);

      if (pathname === '/404') {
        return new Response('<!doctype html><h1>Page not found</h1>', {
          headers: { 'content-type': 'text/html; charset=utf-8' },
        });
      }

      return new Response('generic asset miss', { status: 404 });
    },
  };
  const response = await worker.fetch(
    new Request('https://connect.minekube.com/missing-page'),
    { ASSETS: assets }
  );

  assert.deepEqual(requests, ['/missing-page', '/404']);
  assert.equal(response.status, 404);
  assert.equal(response.headers.get('cache-control'), 'no-store');
  assert.equal(response.headers.get('content-type'), 'text/html; charset=utf-8');
  assert.equal(
    await response.text(),
    '<!doctype html><h1>Page not found</h1>'
  );
});

test('the branded 404 fetch ignores conditional and range request headers', async () => {
  const brandedBody = '<!doctype html><h1>Page not found</h1>';
  const requests = [];
  const assets = {
    fetch(request) {
      requests.push(request);
      if (new URL(request.url).pathname !== '/404') {
        return new Response('generic asset miss', { status: 404 });
      }

      if (request.headers.has('range')) {
        return new Response(brandedBody.slice(0, 12), {
          status: 206,
          headers: {
            'content-length': '12',
            'content-range': 'bytes 0-11/39',
          },
        });
      }

      if (request.headers.has('if-none-match')) {
        return new Response(null, { status: 304 });
      }

      return new Response(brandedBody, {
        headers: { 'content-type': 'text/html; charset=utf-8' },
      });
    },
  };

  const response = await worker.fetch(
    new Request('https://connect.minekube.com/missing-page', {
      headers: {
        'if-modified-since': 'Wed, 21 Oct 2015 07:28:00 GMT',
        'if-none-match': '"cached-404"',
        range: 'bytes=0-11',
      },
    }),
    { ASSETS: assets }
  );

  assert.equal(response.status, 404);
  assert.equal(await response.text(), brandedBody);
  assert.equal(requests[1].method, 'GET');
  assert.equal(requests[1].headers.get('if-modified-since'), null);
  assert.equal(requests[1].headers.get('if-none-match'), null);
  assert.equal(requests[1].headers.get('range'), null);
});

test('the branded 404 response uses a fixed-length body for GET and HEAD', async () => {
  const brandedBody = '<!doctype html><h1>Page not found</h1>';
  const contentLength = new TextEncoder().encode(brandedBody).byteLength;
  const fixedLengths = [];
  const previousFixedLengthStream = globalThis.FixedLengthStream;

  globalThis.FixedLengthStream = class extends TransformStream {
    constructor(length) {
      super();
      fixedLengths.push(length);
    }
  };

  try {
    const assets = {
      fetch(request) {
        if (new URL(request.url).pathname !== '/404') {
          return new Response('generic asset miss', { status: 404 });
        }

        return new Response(brandedBody, {
          headers: {
            'content-length': String(contentLength),
            'content-type': 'text/html; charset=utf-8',
          },
        });
      },
    };

    const response = await worker.fetch(
      new Request('https://connect.minekube.com/missing-page'),
      { ASSETS: assets }
    );
    const head = await worker.fetch(
      new Request('https://connect.minekube.com/missing-page', {
        method: 'HEAD',
      }),
      { ASSETS: assets }
    );

    assert.deepEqual(fixedLengths, [contentLength, contentLength]);
    assert.equal(response.headers.get('content-length'), String(contentLength));
    assert.equal(head.headers.get('content-length'), String(contentLength));
    assert.equal(await response.text(), brandedBody);
  } finally {
    if (previousFixedLengthStream === undefined) {
      delete globalThis.FixedLengthStream;
    } else {
      globalThis.FixedLengthStream = previousFixedLengthStream;
    }
  }
});
