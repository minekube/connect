import assert from 'node:assert/strict'
import test from 'node:test'

import worker, { connectSpigotDownload } from './worker.mjs'

test('redirects the legacy Connect download path to the stable Spigot jar', async () => {
  const response = await worker.fetch(
    new Request('https://developers.minekube.com/connect/download'),
  )

  assert.equal(response.status, 302)
  assert.equal(response.headers.get('location'), connectSpigotDownload)
})

test('redirects a trailing-slash legacy Connect download path too', async () => {
  const response = await worker.fetch(
    new Request('https://developers.minekube.com/connect/download/'),
  )

  assert.equal(response.status, 302)
  assert.equal(response.headers.get('location'), connectSpigotDownload)
})

test('preserves the archived host catch-all for every other path', async () => {
  const response = await worker.fetch(new Request('https://developers.minekube.com/guide/'))

  assert.equal(response.status, 200)
  assert.match(response.headers.get('content-type') ?? '', /^text\/html/)
  assert.equal(
    await response.text(),
    '<!DOCTYPE html><meta http-equiv="Refresh" content="0; url=\'https://minekube.com\'" />',
  )
})
