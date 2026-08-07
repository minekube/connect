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

export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    const redirect = redirects.get(url.pathname);

    if (redirect) {
      return createPagesRedirect(request, redirect, url.search);
    }

    if (url.pathname.endsWith('.html')) {
      return createPagesCanonicalRedirect(request, url.pathname, url.search);
    }

    const response = await env.ASSETS.fetch(request);
    if (response.status !== 404) {
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
