import assert from 'node:assert/strict'
import {spawn} from 'node:child_process'
import {existsSync} from 'node:fs'
import process from 'node:process'
import {chromium} from 'playwright-core'

const host = '127.0.0.1'
const port = Number(process.env.CHANGELOG_TEST_PORT ?? 4174)
const baseURL = `http://${host}:${port}`
const root = new URL('..', import.meta.url)

const chromeCandidates = [
  process.env.CHROME_PATH,
  '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
  '/usr/bin/google-chrome',
  '/usr/bin/google-chrome-stable',
  '/usr/bin/chromium',
  '/usr/bin/chromium-browser',
].filter(Boolean)

const executablePath = chromeCandidates.find(candidate => existsSync(candidate))
assert(executablePath, `Chrome was not found. Checked: ${chromeCandidates.join(', ')}`)

let server
let browser

async function serverReady() {
  try {
    const response = await fetch(`${baseURL}/changelog/`)
    return response.ok
  } catch {
    return false
  }
}

async function startServer() {
  if (await serverReady()) return

  const yarnPath = process.env.npm_execpath
  assert(yarnPath, 'Run this check through Yarn so the VitePress serve command can be resolved.')

  server = spawn(yarnPath, ['serve', '--host', host, '--port', String(port)], {
    cwd: root,
    env: {...process.env, NO_COLOR: '1'},
    stdio: ['ignore', 'pipe', 'pipe'],
  })

  let output = ''
  server.stdout.on('data', chunk => { output += chunk })
  server.stderr.on('data', chunk => { output += chunk })

  for (let attempt = 0; attempt < 80; attempt += 1) {
    if (await serverReady()) return
    if (server.exitCode !== null) throw new Error(`Preview server exited early.\n${output}`)
    await new Promise(resolve => setTimeout(resolve, 100))
  }

  throw new Error(`Preview server did not become ready.\n${output}`)
}

async function readGeometry(page) {
  return page.evaluate(() => {
    const columnRect = document.querySelector('.changelog-column').getBoundingClientRect()
    const feedRect = document.querySelector('.changelog-content').getBoundingClientRect()
    const entries = [...document.querySelectorAll('.changelog-entry:not([hidden])')].map(entry => {
      const dateSlot = entry.querySelector('.changelog-entry-date-slot')
      const date = entry.querySelector('.changelog-entry-date')
      const node = entry.querySelector('.changelog-entry-node')
      const body = entry.querySelector('.changelog-entry-body')
      const entryRect = entry.getBoundingClientRect()
      const dateSlotRect = dateSlot.getBoundingClientRect()
      const dateRect = date.getBoundingClientRect()
      const nodeRect = node.getBoundingClientRect()
      const bodyRect = body.getBoundingClientRect()
      const rail = getComputedStyle(entry, '::before')
      const railTop = Number.parseFloat(rail.top)
      const railBottom = Number.parseFloat(rail.bottom)

      return {
        date: date.textContent.trim(),
        product: entry.querySelector('.VPBadge')?.textContent?.trim(),
        entryTop: entryRect.top,
        entryBottom: entryRect.bottom,
        entryLeft: entryRect.left,
        entryWidth: entryRect.width,
        bodyLeft: bodyRect.left,
        bodyWidth: bodyRect.width,
        dateLeft: dateRect.left,
        dateRight: dateRect.right,
        dateTop: dateRect.top,
        dateBottom: dateRect.bottom,
        datePosition: getComputedStyle(date).position,
        dateSlotTop: dateSlotRect.top,
        dateSlotBottom: dateSlotRect.bottom,
        dateSlotWidth: dateSlotRect.width,
        nodeX: nodeRect.left + nodeRect.width / 2,
        nodeY: nodeRect.top + nodeRect.height / 2,
        nodeWidth: nodeRect.width,
        nodeDisplay: getComputedStyle(node).display,
        railX: entryRect.left + Number.parseFloat(rail.left),
        railStart: entryRect.top + railTop,
        railEnd: entryRect.bottom - railBottom,
        railTop,
        railBottom,
        railDisplay: rail.display,
        railBackground: rail.backgroundImage || rail.backgroundColor,
        end: entry.hasAttribute('data-rail-end'),
      }
    })

    return {
      viewportWidth: innerWidth,
      documentWidth: document.documentElement.scrollWidth,
      columnLeft: columnRect.left,
      columnWidth: columnRect.width,
      feedLeft: feedRect.left,
      entries,
    }
  })
}

function close(actual, expected, message, tolerance = 1) {
  assert(
    Math.abs(actual - expected) <= tolerance,
    `${message}: expected ${expected.toFixed(2)}, received ${actual.toFixed(2)}`,
  )
}

function assertCloudflareRail(geometry, viewport) {
  const {entries} = geometry
  assert(entries.length > 1, `${viewport}: expected multiple visible changelog entries`)
  assert.equal(entries.filter(entry => entry.end).length, 1, `${viewport}: exactly one final rail segment`)
  assert(entries.at(-1).end, `${viewport}: final visible entry owns the fading rail segment`)

  for (const entry of entries) {
    close(entry.railX, entry.nodeX, `${viewport}: rail/node x alignment for ${entry.date}`)
    close(entry.railStart, entry.nodeY, `${viewport}: rail starts at its node for ${entry.date}`)
    close(entry.nodeWidth, 11, `${viewport}: Cloudflare node size for ${entry.date}`)
    close(entry.dateSlotWidth, 88, `${viewport}: Cloudflare date gutter width for ${entry.date}`)
    close(entry.entryLeft - entry.dateRight, 28, `${viewport}: Cloudflare date-to-rail gap for ${entry.date}`)
    close(entry.bodyLeft - entry.entryLeft, 64, `${viewport}: Cloudflare body indent for ${entry.date}`)
    close(entry.bodyWidth, 640, `${viewport}: Cloudflare reading width for ${entry.date}`)
    close((entry.dateTop + entry.dateBottom) / 2, entry.nodeY, `${viewport}: date/node vertical alignment for ${entry.date}`)
    assert.equal(entry.railDisplay, 'block', `${viewport}: entry rail is visible`)
    assert.equal(entry.nodeDisplay, 'block', `${viewport}: entry node is visible`)
  }

  for (let index = 0; index < entries.length - 1; index += 1) {
    close(
      entries[index].railEnd,
      entries[index + 1].railStart,
      `${viewport}: seamless entry-to-entry rail at index ${index}`,
    )
  }

  close(entries.at(-1).railBottom, 0, `${viewport}: final segment ends inside its entry`)
  assert.match(entries.at(-1).railBackground, /linear-gradient/, `${viewport}: final rail fades like Cloudflare's`)
}

try {
  await startServer()
  browser = await chromium.launch({
    executablePath,
    headless: true,
    args: process.getuid?.() === 0 ? ['--no-sandbox'] : [],
  })

  const page = await browser.newPage({viewport: {width: 1440, height: 1200}})
  await page.goto(`${baseURL}/changelog/`, {waitUntil: 'networkidle'})
  await page.waitForSelector('.changelog-entry')

  let geometry = await readGeometry(page)
  assert.equal(geometry.columnWidth, 768, 'desktop: Cloudflare reading column is 768px')
  close(geometry.columnLeft, geometry.feedLeft, 'desktop: feed shares the column left edge')
  assert.equal(geometry.entries.length, 48, 'desktop: every update becomes exactly one timeline entry')
  assertCloudflareRail(geometry, 'desktop')
  assert.equal(
    geometry.entries.slice(0, 3).every(entry => entry.date === 'August 8, 2026'),
    true,
    'desktop: dates repeat for each same-day entry like Cloudflare',
  )
  assert.equal(
    await page.locator('.changelog-date-rail, .changelog-date-node').count(),
    0,
    'desktop: legacy grouped rail elements are gone',
  )

  const stickyTarget = await page.evaluate(() => {
    const entry = [...document.querySelectorAll('.changelog-entry:not([hidden])')]
      .find(candidate => candidate.getBoundingClientRect().height > 240)
    const rect = entry.getBoundingClientRect()
    return {top: scrollY + rect.top, selector: [...document.querySelectorAll('.changelog-entry')].indexOf(entry)}
  })
  await page.evaluate(top => scrollTo(0, top + 80), stickyTarget.top)
  await page.waitForTimeout(100)
  const sticky = await page.locator('.changelog-entry').nth(stickyTarget.selector).evaluate(entry => {
    const date = entry.querySelector('.changelog-entry-date')
    const slot = entry.querySelector('.changelog-entry-date-slot')
    const dateRect = date.getBoundingClientRect()
    const slotRect = slot.getBoundingClientRect()
    return {
      top: dateRect.top,
      bottom: dateRect.bottom,
      slotTop: slotRect.top,
      slotBottom: slotRect.bottom,
      position: getComputedStyle(date).position,
      configuredTop: Number.parseFloat(getComputedStyle(date).top),
    }
  })
  assert.equal(sticky.position, 'sticky', 'desktop: date remains sticky within its individual update')
  close(sticky.top, sticky.configuredTop, 'desktop: active date sticks below navigation')
  assert(sticky.top >= sticky.slotTop && sticky.bottom <= sticky.slotBottom, 'desktop: sticky date stays inside its update')

  await page.evaluate(() => scrollTo(0, 0))
  await page.locator('.changelog-filter select').selectOption('Gate')
  await page.waitForFunction(() => document.querySelector('[aria-live="polite"]')?.textContent.includes('15 updates'))
  geometry = await readGeometry(page)
  assert.equal(geometry.entries.length, 15, 'filtered desktop: only Gate updates remain')
  assert.equal(geometry.entries.every(entry => entry.product === 'Gate'), true, 'filtered desktop: every update is Gate')
  assertCloudflareRail(geometry, 'filtered desktop')

  await page.setViewportSize({width: 390, height: 1000})
  await page.reload({waitUntil: 'networkidle'})
  await page.waitForSelector('.changelog-entry')
  geometry = await readGeometry(page)
  assert.equal(geometry.viewportWidth, 390, 'mobile: browser uses the requested viewport')
  assert(geometry.documentWidth <= geometry.viewportWidth, 'mobile: changelog does not overflow horizontally')
  assert.equal(
    geometry.entries.every(entry => entry.railDisplay === 'none' && entry.nodeDisplay === 'none'),
    true,
    'mobile: Cloudflare rail and nodes collapse away',
  )
  assert.equal(
    geometry.entries.every(entry => entry.datePosition === 'static'),
    true,
    'mobile: dates stack above entry content',
  )
  assert.equal(
    geometry.entries.every(entry => Math.abs(entry.bodyLeft - entry.entryLeft) <= 1),
    true,
    'mobile: entry content uses the full narrow column',
  )

  console.log('Changelog matches Cloudflare rail geometry on desktop, filtering, sticky dates, and mobile.')
} finally {
  await browser?.close()
  if (server && server.exitCode === null) {
    server.kill('SIGTERM')
  }
}
