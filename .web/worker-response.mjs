const shouldAddUtf8Charset = (contentType) => {
  if (!contentType || contentType.toLowerCase().includes('charset=')) {
    return false;
  }

  const mediaType = contentType.split(';', 1)[0].trim().toLowerCase();
  return (
    mediaType.startsWith('text/') ||
    mediaType === 'application/javascript' ||
    mediaType === 'application/x-javascript'
  );
};

export const normalizePagesResponse = (response) => {
  const headers = new Headers(response.headers);

  headers.set('access-control-allow-origin', '*');
  headers.set('x-content-type-options', 'nosniff');
  headers.set('referrer-policy', 'strict-origin-when-cross-origin');

  if (response.status === 404) {
    headers.set('cache-control', 'no-store');
  }

  const contentType = headers.get('content-type');
  if (
    !contentType &&
    response.status >= 300 &&
    response.status < 400 &&
    headers.has('location')
  ) {
    headers.set('content-type', 'text/plain;charset=UTF-8');
  } else if (shouldAddUtf8Charset(contentType)) {
    headers.set('content-type', `${contentType}; charset=utf-8`);
  }

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
};
