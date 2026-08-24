import { normalizePagesResponse } from './worker-response.mjs';

const redirects = new Map([
  [
    '/discord',
    { status: 302, destination: 'https://discord.com/invite/6vMDqWE' },
  ],
  [
    '/discord/',
    { status: 302, destination: 'https://discord.com/invite/6vMDqWE' },
  ],
  ['/guide/changelog', { status: 301, destination: '/changelog/' }],
  ['/guide/changelog/', { status: 301, destination: '/changelog/' }],
  ['/guide/changelog.html', { status: 301, destination: '/changelog/' }],
]);

const createPagesRedirect = (request, redirect, search) => {
  const location = `${redirect.destination}${search}`;
  const body = `Redirecting to ${location}`;
  const headers = new Headers({
    'access-control-allow-origin': '*',
    'content-type': 'text/plain;charset=UTF-8',
    location,
    'referrer-policy': 'strict-origin-when-cross-origin',
    'x-content-type-options': 'nosniff',
  });
  let responseBody = null;

  if (request.method !== 'HEAD') {
    const encodedBody = new TextEncoder().encode(body);
    headers.set('content-length', String(encodedBody.byteLength));
    responseBody = body;

    if (typeof FixedLengthStream !== 'undefined') {
      const { readable, writable } = new FixedLengthStream(
        encodedBody.byteLength
      );
      const writer = writable.getWriter();
      void writer.write(encodedBody);
      void writer.close();
      responseBody = readable;
    }
  }

  return new Response(responseBody, {
    status: redirect.status,
    headers,
  });
};

const createPagesCanonicalRedirect = (request, pathname, search) => {
  const destination = pathname.endsWith('/index.html')
    ? pathname.slice(0, -'index.html'.length)
    : pathname.slice(0, -'.html'.length);
  const headers = new Headers({
    'access-control-allow-origin': '*',
    location: `${destination}${search}`,
    'referrer-policy': 'strict-origin-when-cross-origin',
  });

  if (request.method !== 'HEAD') {
    headers.set('content-length', '0');
  }

  return new Response(null, {
    status: 308,
    headers,
  });
};

const createFixedLengthBody = (response) => {
  const body = response.body;
  const contentLength = response.headers.get('content-length');

  if (
    body === null ||
    contentLength === null ||
    !/^\d+$/.test(contentLength) ||
    !Number.isSafeInteger(Number(contentLength)) ||
    typeof FixedLengthStream === 'undefined'
  ) {
    return body;
  }

  const { readable, writable } = new FixedLengthStream(Number(contentLength));
  void body.pipeTo(writable);
  return readable;
};

const normalizeMarkdownResponse = (response) => {
  const normalized = normalizePagesResponse(response);
  const headers = new Headers(normalized.headers);
  headers.set('content-type', 'text/markdown; charset=utf-8');

  return new Response(normalized.body, {
    status: normalized.status,
    statusText: normalized.statusText,
    headers,
  });
};

export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    const redirect = redirects.get(url.pathname);

    if (redirect) {
      return createPagesRedirect(request, redirect, url.search);
    }

    // Public Bedrock verifier metadata is proxied to the Connect watch so the
    // connect-java plugin's pinned metadata-origin (connect.minekube.com) can
    // resolve both v1 identity keys and the v2 Bedrock principal key. The
    // watch serves the canonical responses with Cache-Control/ETag; pass the
    // origin response through untouched (including 404s when signing is not
    // configured) so plugin capability checks see the same reality as direct
    // watch access.
    if (url.pathname.startsWith('/.well-known/minekube-connect/')) {
      const upstream = 'https://watch-connect.minekube.net' + url.pathname + url.search;
      const upstreamRequest = new Request(upstream, {
        method: request.method,
        headers: request.headers,
      });
      const upstreamResponse = await fetch(upstreamRequest);
      const headers = new Headers(upstreamResponse.headers);
      headers.set('access-control-allow-origin', '*');
      return new Response(upstreamResponse.body, {
        status: upstreamResponse.status,
        statusText: upstreamResponse.statusText,
        headers,
      });
    }

    if (url.pathname.endsWith('.html')) {
      return createPagesCanonicalRedirect(request, url.pathname, url.search);
    }

    const response = await env.ASSETS.fetch(request);
    if (response.status !== 404) {
      if (url.pathname.endsWith('.md')) {
        return normalizeMarkdownResponse(response);
      }

      return normalizePagesResponse(response);
    }

    const notFoundRequest = new Request(new URL('/404', request.url), {
      method: 'GET',
    });
    const notFoundPage = await env.ASSETS.fetch(notFoundRequest);

    return normalizePagesResponse(
      new Response(createFixedLengthBody(notFoundPage), {
        status: 404,
        headers: notFoundPage.headers,
      })
    );
  },
};
