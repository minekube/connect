#!/usr/bin/env node

const url = process.argv[2] ?? 'https://developers.minekube.com/connect/download'
const minimumJarBytes = 1_000_000
const allowedContentTypes = new Set([
  'application/java-archive',
  'application/octet-stream',
  'application/zip',
])

const response = await fetch(url, {
  headers: {
    range: 'bytes=0-3',
  },
  redirect: 'follow',
})

if (response.status !== 206) {
  throw new Error(`expected final HTTP 206 range response, got ${response.status}`)
}

const contentType = response.headers.get('content-type')?.split(';', 1)[0].trim()
if (!contentType || !allowedContentTypes.has(contentType)) {
  throw new Error(`expected jar content type, got ${contentType ?? 'missing'}`)
}

const contentRange = response.headers.get('content-range')
const totalBytes = Number.parseInt(contentRange?.match(/^bytes 0-3\/(\d+)$/)?.[1] ?? '', 10)
if (!Number.isSafeInteger(totalBytes) || totalBytes < minimumJarBytes) {
  throw new Error(
    `expected jar size of at least ${minimumJarBytes} bytes, got ${contentRange ?? 'missing content-range'}`,
  )
}

const bytes = new Uint8Array(await response.arrayBuffer())
if (bytes.length !== 4 || bytes[0] !== 0x50 || bytes[1] !== 0x4b) {
  throw new Error(`expected PK jar magic, got ${Buffer.from(bytes).toString('hex') || 'no bytes'}`)
}

if (!response.redirected) {
  throw new Error('expected at least one HTTP redirect before the final jar response')
}

console.log(
  JSON.stringify({
    source: url,
    finalUrl: response.url,
    status: response.status,
    contentType,
    totalBytes,
    magic: Buffer.from(bytes).toString('hex'),
  }),
)
