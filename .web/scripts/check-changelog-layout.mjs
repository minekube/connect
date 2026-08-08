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
    return [...document.querySelectorAll('.changelog-date-group:not([hidden])')].map(group => {
      const heading = group.querySelector('h2')
      const railElement = group.querySelector(':scope > .changelog-date-rail')
      const nodeElement = group.querySelector(':scope > .changelog-date-node')
      const groupRect = group.getBoundingClientRect()
      const headingRect = heading.getBoundingClientRect()
      const railRect = railElement.getBoundingClientRect()
      const nodeRect = nodeElement.getBoundingClientRect()
      const rail = getComputedStyle(railElement)

      return {
        date: heading.textContent.trim(),
        groupTop: groupRect.top,
        groupBottom: groupRect.bottom,
        dateTop: headingRect.top,
        dateBottom: headingRect.bottom,
        datePosition: getComputedStyle(heading).position,
        railTop: railRect.top,
        railBottom: railRect.bottom,
        railX: railRect.left + railRect.width / 2,
        nodeX: nodeRect.left + nodeRect.width / 2,
        nodeY: nodeRect.top + nodeRect.height / 2,
        nodeOffset: nodeRect.top + nodeRect.height / 2 - groupRect.top,
        railBackgroundColor: rail.backgroundColor,
        start: group.classList.contains('is-visible-start'),
        end: group.classList.contains('is-visible-end'),
      }
    })
  })
}

function close(actual, expected, message, tolerance = 1) {
  assert(
    Math.abs(actual - expected) <= tolerance,
    `${message}: expected ${expected.toFixed(2)}, received ${actual.toFixed(2)}`,
  )
}

function assertRailGeometry(groups, viewport) {
  assert(groups.length > 1, `${viewport}: expected multiple visible date groups`)
  assert.equal(groups.filter(group => group.start).length, 1, `${viewport}: one visible rail start`)
  assert.equal(groups.filter(group => group.end).length, 1, `${viewport}: one visible rail end`)
  assert(groups[0].start, `${viewport}: first visible group owns the rail start`)
  assert(groups.at(-1).end, `${viewport}: last visible group owns the rail end`)

  for (const group of groups) {
    close(group.railX, group.nodeX, `${viewport}: rail/node x alignment for ${group.date}`)
    assert(group.railTop >= group.groupTop, `${viewport}: ${group.date} rail may not escape above its group`)
    assert(group.railBottom <= group.groupBottom, `${viewport}: ${group.date} rail may not escape below its group`)
    assert.notEqual(group.railBackgroundColor, 'rgba(0, 0, 0, 0)', `${viewport}: rail stays visible`)
  }

  close(groups[0].railTop, groups[0].nodeY, `${viewport}: rail starts at the first node`)
  for (let index = 0; index < groups.length - 1; index += 1) {
    close(
      groups[index].railBottom,
      groups[index + 1].railTop,
      `${viewport}: connected rail between ${groups[index].date} and ${groups[index + 1].date}`,
    )
  }
  close(groups.at(-1).railBottom, groups.at(-1).nodeY, `${viewport}: rail stops at the final node`)
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
  await page.waitForSelector('.changelog-date-group')

  let groups = await readGeometry(page)
  assert.equal(groups.length, 26, 'desktop: every date is grouped exactly once')
  assertRailGeometry(groups, 'desktop')

  const stickyTarget = await page.evaluate(() => {
    const visibleGroups = [...document.querySelectorAll('.changelog-date-group:not([hidden])')]
    const index = visibleGroups.findIndex(candidate => candidate.getBoundingClientRect().height > 600)
    const group = visibleGroups[index]
    const nextGroup = visibleGroups[index + 1]
    const rect = group.getBoundingClientRect()
    const nodeRect = group.querySelector(':scope > .changelog-date-node').getBoundingClientRect()
    const nextHeadingRect = nextGroup.querySelector('h2').getBoundingClientRect()
    return {
      top: scrollY + rect.top,
      bottom: scrollY + rect.bottom,
      nodeOffset: nodeRect.top + nodeRect.height / 2 - rect.top,
      id: group.getAttribute('aria-labelledby'),
      nextId: nextGroup.getAttribute('aria-labelledby'),
      nextHeadingTop: scrollY + nextHeadingRect.top,
    }
  })
  await page.evaluate(top => scrollTo(0, top + 260), stickyTarget.top)
  await page.waitForTimeout(100)

  const sticky = await page.evaluate(id => {
    const heading = document.getElementById(id)
    const group = heading.closest('.changelog-date-group')
    const headingRect = heading.getBoundingClientRect()
    const groupRect = group.getBoundingClientRect()
    return {
      top: headingRect.top,
      bottom: headingRect.bottom,
      groupTop: groupRect.top,
      groupBottom: groupRect.bottom,
      configuredTop: Number.parseFloat(getComputedStyle(heading).top),
      position: getComputedStyle(heading).position,
    }
  }, stickyTarget.id)
  assert.equal(sticky.position, 'sticky', 'desktop: release date uses sticky positioning')
  close(sticky.top, sticky.configuredTop, 'desktop: active date sticks at configured top')
  assert(sticky.top >= sticky.groupTop && sticky.bottom <= sticky.groupBottom, 'desktop: sticky date stays inside its release block')

  const staticNode = await page.evaluate(id => {
    const heading = document.getElementById(id)
    const group = heading.closest('.changelog-date-group')
    const node = group.querySelector(':scope > .changelog-date-node')
    const groupRect = group.getBoundingClientRect()
    const nodeRect = node.getBoundingClientRect()
    return nodeRect.top + nodeRect.height / 2 - groupRect.top
  }, stickyTarget.id)
  close(staticNode, stickyTarget.nodeOffset, 'desktop: sticky date does not drag its timeline node')

  await page.evaluate(({bottom, top}) => scrollTo(0, bottom - top + 2), {
    bottom: stickyTarget.bottom,
    top: sticky.configuredTop,
  })
  await page.waitForTimeout(100)
  const outgoingBoundary = await page.evaluate(id => {
    const heading = document.getElementById(id)
    const group = heading.closest('.changelog-date-group')
    return {
      dateBottom: heading.getBoundingClientRect().bottom,
      groupBottom: group.getBoundingClientRect().bottom,
    }
  }, stickyTarget.id)
  assert(
    outgoingBoundary.dateBottom <= outgoingBoundary.groupBottom + 0.5,
    'desktop: outgoing sticky date stops at its own release boundary',
  )

  await page.evaluate(({headingTop, stickyTop}) => scrollTo(0, headingTop - stickyTop + 2), {
    headingTop: stickyTarget.nextHeadingTop,
    stickyTop: sticky.configuredTop,
  })
  await page.waitForTimeout(100)
  const incomingTop = await page.locator(`#${stickyTarget.nextId}`).evaluate(heading => heading.getBoundingClientRect().top)
  close(incomingTop, sticky.configuredTop, 'desktop: incoming date takes over at the sticky boundary')

  await page.evaluate(() => scrollTo(0, 0))
  await page.locator('.changelog-filter select').selectOption('Gate')
  await page.waitForFunction(() => document.querySelector('.changelog-result-count')?.textContent.includes('15 updates'))
  groups = await readGeometry(page)
  assert.equal(groups.length, 12, 'filtered desktop: only Gate date groups remain')
  assertRailGeometry(groups, 'filtered desktop')
  assert.equal(
    await page.locator('.changelog-date-group li:not([hidden]) .VPBadge:first-child').allTextContents()
      .then(labels => labels.every(label => label.trim() === 'Gate')),
    true,
    'filtered desktop: every visible entry belongs to Gate',
  )

  await page.setViewportSize({width: 390, height: 1000})
  await page.reload({waitUntil: 'networkidle'})
  await page.evaluate(() => scrollTo(0, 0))
  await page.waitForSelector('.changelog-date-group')
  groups = await readGeometry(page)
  assertRailGeometry(groups, 'mobile')
  assert(groups.every(group => group.datePosition === 'relative'), 'mobile: dates do not stick over narrow content')

  console.log('Changelog browser geometry passed for desktop, filtering, sticky boundaries, and mobile.')
} finally {
  await browser?.close()
  if (server && server.exitCode === null) {
    server.kill('SIGTERM')
  }
}
