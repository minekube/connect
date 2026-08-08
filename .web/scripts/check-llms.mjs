import {existsSync, readFileSync} from 'node:fs'
import {resolve} from 'node:path'

const root = new URL('..', import.meta.url)
const dist = resolve(root.pathname, 'docs/.vitepress/dist')

function read(path) {
  const file = resolve(dist, path)
  if (!existsSync(file)) throw new Error(`Missing generated LLM route: ${path}`)
  return readFileSync(file, 'utf8')
}

function assertIncludes(content, expected, file) {
  if (!content.includes(expected)) {
    throw new Error(`${file} is missing generated LLM content: ${expected}`)
  }
}

function assertNotIncludes(content, unexpected, file) {
  if (content.includes(unexpected)) {
    throw new Error(`${file} contains unexpected generated LLM content: ${unexpected}`)
  }
}

const index = read('llms.txt')
const full = read('llms-full.txt')
const joining = read('guide/joining.md')
const bedrock = read('guide/bedrock.md')
const offline = read('guide/offline-mode.md')

for (const [file, content, required] of [
  ['llms.txt', index, ['# Minekube Connect', 'https://connect.minekube.com/guide/joining.md']],
  ['llms-full.txt', full, ['## Who Can Join', '## Allowing Offline Java on an Endpoint']],
  ['guide/joining.md', joining, ['Stable native Bedrock/XUID-derived profile', 'Bedrock client without valid Microsoft/Xbox authentication']],
  ['guide/bedrock.md', bedrock, ['without owning or linking Java Edition', 'stable profile derived from the verified Bedrock XUID']],
  ['guide/offline-mode.md', offline, ['cracked.minekube.com', 'allowOfflineModePlayers: true']],
]) {
  for (const expected of required) assertIncludes(content, expected, file)
  assertNotIncludes(content, '<!--@include:', file)
}

for (const excluded of [
  '/changelog/2026-',
  '/guide/downloads.md',
  '/guide/games.md',
  '/guide/includes/',
]) {
  assertNotIncludes(index, excluded, 'llms.txt')
}

for (const match of index.matchAll(/\(https:\/\/connect\.minekube\.com(\/[^)#?]+\.md)/g)) {
  const output = match[1].slice(1)
  if (!existsSync(resolve(dist, output))) {
    throw new Error(`llms.txt points to a missing generated Markdown route: ${match[1]}`)
  }
}

console.log('Generated LLM routes passed.')
