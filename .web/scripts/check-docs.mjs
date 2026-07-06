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
  'use the normal Bedrock port `19132`',
  'It does not mean you need to open UDP `19132`',
  'official Microsoft/Xbox Bedrock authentication',
  'Bedrock Identity Enforcement',
  'metadata-url',
  'endpoint and organization',
  'Discord support response draft',
])

assertAll('docs/guide/offline-mode.md', [
  'Connect-managed Bedrock identity',
  'official Microsoft/Xbox Bedrock auth',
  'Do not enable `allow-offline-mode-players` just because a Bedrock player is joining through Connect',
])

assertAll('docs/guide/connectors/plugin.md', [
  'Bedrock Identity',
  'metadata-url',
  'enforcement: warn',
])

assertAll('docs/guide/compatibility.md', [
  'Velocity snapshots',
  'NLogin',
  'AuthMe',
  'FastLogin',
  'Gate Lite with an online-mode backend',
  'AuthSession/passthrough for Lite backend routes is not available today',
  'Arclight/Ketting/Forge hybrids',
  'FabricProxy-Lite',
  'CrossStitch',
  'Polymer',
  'NeoForge 1.21.x / Proxy-Compatible-Forge through Connect',
  'Connect compatibility investigation',
])

assertAll('docs/guide/connectors/gate.md', [
  'Current behavior',
  'Not supported today',
  'Connect -> Gate Lite -> Online Mode Backend',
  'Connect passthrough/AuthSession support',
])

assertAll('docs/guide/auth-api.md', [
  'The AuthSession API is not available in production yet',
  'Gate Lite behind Connect is different',
  'Use standard Gate with Connect enabled or the Connect Java Plugin',
])

console.log('Docs content assertions passed.')
