<!--
  Timeline geometry adapted from Cloudflare's ChangelogLayout, ChangelogFeed,
  ChangelogEntry, DateRail, and Header components at cloudflare/cloudflare-docs.
  Copyright (c) 2021 Cloudflare, Inc.; MIT license in .web/THIRD_PARTY_NOTICES.md.
-->
<script setup lang="ts">
import {nextTick, onMounted, ref, watch} from 'vue'
import {useData} from 'vitepress'

const {frontmatter} = useData()
const content = ref<HTMLElement | null>(null)
const selectedProduct = ref('All products')
const visibleCount = ref(0)
const indexed = ref(false)

const products = ['All products', 'Connect', 'Connect plugin', 'Gate', 'GeyserLite', 'Craftless']

function dateTimeValue(label: string) {
  const parsed = new Date(`${label} 00:00:00 UTC`)
  return Number.isNaN(parsed.getTime()) ? '' : parsed.toISOString().slice(0, 10)
}

function createEntries() {
  if (!content.value) return

  for (const heading of content.value.querySelectorAll<HTMLHeadingElement>(':scope > h2')) {
    const list = heading.nextElementSibling
    if (!(list instanceof HTMLUListElement)) continue

    const dateLabel = heading.textContent?.replace(/[\u200B#]/g, '').trim() ?? ''
    const group = document.createElement('section')
    group.className = 'changelog-date-group'
    if (heading.id) group.setAttribute('aria-labelledby', heading.id)

    heading.before(group)
    group.append(heading, list)
    heading.classList.add('changelog-source-date')
    list.classList.add('changelog-entry-list')

    for (const entry of list.querySelectorAll<HTMLLIElement>(':scope > li')) {
      entry.classList.add('changelog-entry')

      const dateSlot = document.createElement('span')
      dateSlot.className = 'changelog-entry-date-slot'
      const time = document.createElement('time')
      time.className = 'changelog-entry-date'
      time.textContent = dateLabel
      const dateTime = dateTimeValue(dateLabel)
      if (dateTime) time.dateTime = dateTime
      dateSlot.append(time)

      const node = document.createElement('span')
      node.className = 'changelog-entry-node'
      node.setAttribute('aria-hidden', 'true')

      const title = entry.querySelector<HTMLElement>(':scope > strong')
      const badges = [...entry.querySelectorAll<HTMLElement>(':scope > .VPBadge')]
      if (title && badges.length) {
        const pills = document.createElement('span')
        pills.className = 'changelog-entry-pills'
        pills.append(...badges)
        title.after(pills)
      }

      const body = document.createElement('div')
      body.className = 'changelog-entry-body'
      body.append(...entry.childNodes)
      entry.append(dateSlot, node, body)
    }
  }
}

function applyProductFilter() {
  if (!content.value) return

  createEntries()

  const entries = [...content.value.querySelectorAll<HTMLLIElement>('.changelog-entry')]
  for (const entry of entries) {
    const product = entry.querySelector<HTMLElement>('.VPBadge')?.textContent?.trim()
    entry.hidden = selectedProduct.value !== 'All products' && product !== selectedProduct.value
    entry.removeAttribute('data-rail-end')
  }

  for (const group of content.value.querySelectorAll<HTMLElement>('.changelog-date-group')) {
    const entriesInGroup = [...group.querySelectorAll<HTMLLIElement>('.changelog-entry')]
    group.hidden = !entriesInGroup.some(entry => !entry.hidden)
  }

  const visibleEntries = entries.filter(entry => !entry.hidden)
  visibleEntries.at(-1)?.setAttribute('data-rail-end', '')
  visibleCount.value = visibleEntries.length
  indexed.value = true
}

onMounted(() => nextTick(applyProductFilter))
watch(selectedProduct, () => nextTick(applyProductFilter))
</script>

<template>
  <div class="changelog-shell">
    <div class="changelog-backdrop" aria-hidden="true">
      <span class="changelog-backdrop-rule changelog-backdrop-rule-left" />
      <span class="changelog-backdrop-rule changelog-backdrop-rule-right" />
      <span class="changelog-backdrop-dots" />
      <span class="changelog-backdrop-rule changelog-backdrop-rule-outer" />
    </div>

    <div class="changelog-column">
      <header class="changelog-header">
        <h1>{{ frontmatter.title }}</h1>
        <p>{{ frontmatter.description }}</p>
      </header>

      <section class="changelog-toolbar" aria-label="Filter and subscribe to the changelog">
        <label class="changelog-filter">
          <span class="visually-hidden">Filter changelog by product</span>
          <select id="changelog-product" v-model="selectedProduct" name="product" aria-label="Filter changelog by product">
            <option v-for="product in products" :key="product" :value="product">{{ product }}</option>
          </select>
        </label>

        <div class="changelog-toolbar-links">
          <a class="changelog-minimal-link" href="/llms.txt">View LLM docs</a>
          <a class="changelog-rss-link" href="/changelog.rss" target="_blank" rel="noreferrer" aria-label="Subscribe to the changelog RSS feed">
            <svg aria-hidden="true" viewBox="0 0 24 24"><path d="M5 4a15 15 0 0 1 15 15h-3A12 12 0 0 0 5 7V4Zm0 6a9 9 0 0 1 9 9h-3a6 6 0 0 0-6-6v-3Zm2 7a2 2 0 1 1-4 0 2 2 0 0 1 4 0Z" /></svg>
            Subscribe to RSS
          </a>
        </div>
        <span class="visually-hidden" aria-live="polite">{{ visibleCount }} updates shown</span>
      </section>

      <article ref="content" class="changelog-content">
        <slot />
        <p v-if="indexed && visibleCount === 0" class="changelog-empty">No updates match this product.</p>
      </article>
    </div>
  </div>
</template>

<style scoped>
.changelog-shell {
  --changelog-column-width: min(768px, calc(100vw - 3rem));
  --changelog-outer-width: 680px;
  --changelog-border: color-mix(in srgb, var(--vp-c-divider) 82%, transparent);
  --changelog-border-strong: color-mix(in srgb, var(--vp-c-text-3) 72%, var(--vp-c-divider));
  position: relative;
  width: 100%;
  min-height: 100vh;
}

.changelog-column {
  position: relative;
  z-index: 1;
  width: var(--changelog-column-width);
  margin-inline: auto;
  padding-bottom: 48px;
}

.changelog-backdrop {
  position: absolute;
  z-index: 0;
  inset: 0;
  display: none;
  overflow: hidden;
  pointer-events: none;
}

.changelog-backdrop-rule {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 1px;
  transform: translateX(-50%);
  background-image: linear-gradient(to bottom, var(--changelog-border) 50%, transparent 50%);
  background-size: 1px 32px;
  background-repeat: repeat-y;
}

.changelog-backdrop-rule-left {
  left: calc(50% - var(--changelog-column-width) / 2);
}

.changelog-backdrop-rule-right {
  left: calc(50% + var(--changelog-column-width) / 2);
}

.changelog-backdrop-rule-outer {
  left: calc(50% + var(--changelog-outer-width));
}

.changelog-backdrop-dots {
  position: absolute;
  top: 0;
  bottom: 0;
  left: calc(50% + var(--changelog-column-width) / 2);
  width: calc(var(--changelog-outer-width) - var(--changelog-column-width) / 2);
  background-image: radial-gradient(var(--changelog-border) 0.75px, transparent 0.75px);
  background-position: left top;
  background-size: 12px 12px;
}

.changelog-header {
  max-width: 672px;
  margin-inline: auto;
  padding-top: 56px;
  text-align: center;
}

.changelog-header h1 {
  margin: 0;
  border: 0;
  color: var(--vp-c-text-1);
  font-size: clamp(44px, 5vw, 56px);
  font-weight: 500;
  letter-spacing: -0.035em;
  line-height: 1.05;
  text-wrap: balance;
}

.changelog-header p {
  max-width: 448px;
  margin: 16px auto 0;
  color: var(--vp-c-text-2);
  font-size: 18px;
  line-height: 1.625;
  text-wrap: pretty;
}

.changelog-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin: 48px 0;
  padding: 40px 64px 0;
}

.changelog-filter {
  width: 100%;
  max-width: 259px;
}

.changelog-filter select {
  width: 100%;
  height: 40px;
  padding: 0 40px 0 14px;
  border: 1px solid var(--changelog-border);
  border-radius: 8px;
  color: var(--vp-c-text-1);
  background-color: var(--vp-c-bg);
  font: inherit;
  font-size: 16px;
  font-weight: 400;
  transition: border-color 160ms ease;
}

.changelog-filter select:hover {
  border-color: var(--changelog-border-strong);
}

.changelog-filter select:focus-visible {
  outline: 2px solid var(--vp-c-brand-1);
  outline-offset: -2px;
}

.changelog-toolbar-links {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

.changelog-toolbar-links a {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 40px;
  padding: 0 14px;
  border-radius: 8px;
  color: var(--vp-c-text-1);
  font-size: 14px;
  font-weight: 600;
  line-height: 1;
  text-decoration: none;
}

.changelog-toolbar-links a:hover {
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg-soft);
}

.changelog-rss-link {
  gap: 8px;
  border: 1px solid var(--changelog-border);
  background: var(--vp-c-bg-soft);
}

.changelog-rss-link svg {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  clip-path: inset(50%);
  white-space: nowrap;
}

.changelog-content {
  --body-x: 56px;
  --date-w: 88px;
  --node-y: 14.4px;
  position: relative;
}

.changelog-content :deep(> h1) {
  display: none;
}

.changelog-content :deep(.changelog-policy) {
  margin: 0 0 52px;
  padding: 0 var(--body-x) 24px;
  border-bottom: 1px solid var(--changelog-border);
  color: var(--vp-c-text-2);
  font-size: 14px;
  line-height: 1.7;
}

.changelog-content :deep(.changelog-policy summary) {
  cursor: pointer;
  color: var(--vp-c-text-1);
  font-weight: 600;
}

.changelog-content :deep(.changelog-policy p:last-child) {
  margin-bottom: 0;
}

.changelog-content :deep(.changelog-date-group),
.changelog-content :deep(.changelog-entry-list) {
  display: contents;
}

.changelog-content :deep(.changelog-source-date) {
  display: none;
}

.changelog-content :deep(.changelog-entry) {
  position: relative;
  margin: 0;
  padding: 0 0 52px var(--body-x);
  color: var(--vp-c-text-2);
  line-height: 1.72;
  list-style: none;
}

.changelog-content :deep(.changelog-entry::before) {
  position: absolute;
  top: var(--node-y);
  bottom: calc(-1 * var(--node-y));
  left: 0;
  display: none;
  width: 1px;
  transform: translateX(-50%);
  background: var(--changelog-border-strong);
  content: '';
}

.changelog-content :deep(.changelog-entry[data-rail-end]::before) {
  bottom: 0;
  background: linear-gradient(
    to bottom,
    var(--changelog-border-strong) 0,
    var(--changelog-border-strong) calc(100% - 48px),
    transparent 100%
  );
}

.changelog-content :deep(.changelog-entry-node) {
  position: absolute;
  z-index: 2;
  top: var(--node-y);
  left: 0;
  display: none;
  width: 11px;
  height: 11px;
  box-sizing: border-box;
  transform: translate(-50%, -50%);
  border: 1.5px solid var(--changelog-border-strong);
  border-radius: 3px;
  background: var(--vp-c-bg);
  transition: border-color 160ms cubic-bezier(0.23, 1, 0.32, 1), background-color 160ms cubic-bezier(0.23, 1, 0.32, 1);
}

.changelog-content :deep(.changelog-entry:hover > .changelog-entry-node) {
  border-color: var(--vp-c-text-1);
  background-color: var(--vp-c-text-1);
}

.changelog-content :deep(.changelog-entry-date-slot) {
  position: absolute;
  top: 0;
  right: calc(100% + 28px);
  bottom: 52px;
  width: var(--date-w);
  padding-top: 7.25px;
  color: var(--vp-c-text-3);
  font-family: var(--vp-font-family-mono);
  font-size: 13px;
  font-weight: 500;
  line-height: 1.1;
  text-align: right;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.changelog-content :deep(.changelog-entry-date) {
  position: sticky;
  top: 92px;
  display: block;
}

.changelog-content :deep(.changelog-entry-body) {
  max-width: 640px;
  font-size: 16px;
  text-wrap: pretty;
}

.changelog-content :deep(.changelog-entry-body > strong) {
  display: block;
  margin: 0;
  color: var(--vp-c-text-1);
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.02em;
  line-height: 1.25;
  text-wrap: balance;
}

.changelog-content :deep(.changelog-entry-pills) {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 11px 0 18px;
}

.changelog-content :deep(.changelog-entry-pills .VPBadge) {
  margin: 0;
}

.changelog-empty {
  margin: 0;
  padding: 0 var(--body-x) 48px;
  color: var(--vp-c-text-2);
  font-size: 15px;
}

@media (min-width: 768px) {
  .changelog-backdrop,
  .changelog-content :deep(.changelog-entry::before),
  .changelog-content :deep(.changelog-entry-node) {
    display: block;
  }
}

@media (min-width: 1024px) {
  .changelog-content {
    --body-x: 64px;
  }
}

@media (max-width: 767px) {
  .changelog-shell {
    --changelog-column-width: calc(100vw - 2rem);
  }

  .changelog-header {
    padding-top: 36px;
  }

  .changelog-header h1 {
    font-size: 42px;
  }

  .changelog-header p {
    font-size: 16px;
  }

  .changelog-toolbar {
    align-items: stretch;
    flex-direction: column;
    margin: 36px 0 44px;
    padding: 28px 0 0;
  }

  .changelog-filter {
    max-width: none;
  }

  .changelog-toolbar-links {
    justify-content: space-between;
    margin-left: 0;
  }

  .changelog-content {
    --body-x: 0px;
  }

  .changelog-content :deep(.changelog-policy) {
    margin-bottom: 44px;
  }

  .changelog-content :deep(.changelog-entry) {
    padding-bottom: 44px;
  }

  .changelog-content :deep(.changelog-entry-date-slot) {
    position: static;
    display: block;
    width: auto;
    margin-bottom: 8px;
    padding: 0;
    text-align: left;
  }

  .changelog-content :deep(.changelog-entry-date) {
    position: static;
  }

  .changelog-content :deep(.changelog-entry-body > strong) {
    font-size: 20px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .changelog-content :deep(.changelog-entry-node) {
    transition: none;
  }
}
</style>
