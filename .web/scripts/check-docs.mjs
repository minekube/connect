import {readFileSync} from 'node:fs'
import {resolve} from 'node:path'

const root = new URL('..', import.meta.url)

function readDoc(path) {
  return readFileSync(resolve(root.pathname, path), 'utf8')
}

function assertIncludes(content, expected, file) {
  if (!content.includes(expected)) {
    throw new Error(`${file} is missing required docs coverage: ${expected}`)
  }
}

function assertAll(file, required) {
  const content = readDoc(file)

  for (const expected of required) {
    assertIncludes(content, expected, file)
  }
}

assertAll('docs/guide/bedrock.md', [
  'Connect-managed Bedrock',
  'self-hosted Gate direct Bedrock',
  'backend Floodgate API',
  'online-mode',
  'custom domain',
  'Gate Lite',
  'Paper/Velocity/Bungee connector',
  'Do not enable backend Gate `bedrock: true` for Connect-managed Bedrock',
  'Discord support response draft',
])

assertAll('docs/guide/compatibility.md', [
  'Velocity snapshots',
  'NLogin',
  'AuthMe',
  'FastLogin',
  'Arclight/Ketting/Forge hybrids',
  'FabricProxy-Lite',
  'CrossStitch',
  'Polymer',
  'NeoForge 1.21.x / Proxy-Compatible-Forge through Connect',
  'Connect compatibility investigation',
])

console.log('Docs content assertions passed.')
