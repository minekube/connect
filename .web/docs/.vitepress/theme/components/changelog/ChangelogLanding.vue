<script setup lang="ts">
import {nextTick, onMounted, ref, watch} from 'vue'
import {useData} from 'vitepress'

const {frontmatter} = useData()
const content = ref<HTMLElement | null>(null)
const selectedProduct = ref('All products')
const visibleCount = ref(0)
const indexed = ref(false)

const products = ['All products', 'Connect', 'Connect plugin', 'Gate', 'GeyserLite', 'Craftless']

function createDateGroups() {
  if (!content.value) return

  for (const heading of content.value.querySelectorAll<HTMLHeadingElement>(':scope > h2')) {
    const list = heading.nextElementSibling
    if (!(list instanceof HTMLUListElement)) continue

    const group = document.createElement('section')
    group.className = 'changelog-date-group'
    if (heading.id) group.setAttribute('aria-labelledby', heading.id)
    const rail = document.createElement('span')
    rail.className = 'changelog-date-rail'
    rail.setAttribute('aria-hidden', 'true')
    const node = document.createElement('span')
    node.className = 'changelog-date-node'
    node.setAttribute('aria-hidden', 'true')
    heading.before(group)
    group.append(rail, node, heading, list)
  }
}

function applyProductFilter() {
  if (!content.value) return

  createDateGroups()

  const entries = [...content.value.querySelectorAll<HTMLLIElement>('.changelog-date-group > ul > li')]
  for (const entry of entries) {
    const product = entry.querySelector<HTMLElement>('.VPBadge')?.textContent?.trim()
    entry.hidden = selectedProduct.value !== 'All products' && product !== selectedProduct.value
  }

  const groups = [...content.value.querySelectorAll<HTMLElement>('.changelog-date-group')]
  for (const group of groups) {
    const list = group.querySelector<HTMLUListElement>(':scope > ul')
    if (!list) continue
    const hasVisibleEntry = [...list.children].some(entry => !(entry as HTMLElement).hidden)
    group.hidden = !hasVisibleEntry
    group.classList.remove('is-visible-start')
    group.classList.remove('is-visible-end')
  }

  const visibleGroups = groups.filter(group => !group.hidden)
  visibleGroups[0]?.classList.add('is-visible-start')
  visibleGroups.at(-1)?.classList.add('is-visible-end')

  visibleCount.value = entries.filter(entry => !entry.hidden).length
  indexed.value = true
}

onMounted(() => nextTick(applyProductFilter))
watch(selectedProduct, () => nextTick(applyProductFilter))
</script>

<template>
  <div class="changelog-shell">
    <header class="changelog-hero">
      <div class="changelog-hero-copy">
        <p class="changelog-kicker">Minekube release history</p>
        <h1>{{ frontmatter.title }}</h1>
        <p class="changelog-intro">{{ frontmatter.description }}</p>
        <div class="changelog-actions">
          <a class="changelog-primary-action" href="#august-8-2026">Read the latest changes</a>
          <a
            class="changelog-secondary-action"
            href="/changelog.rss"
            target="_blank"
            rel="noreferrer"
          >Subscribe by RSS</a>
        </div>
      </div>

      <aside class="changelog-summary" aria-label="How to read the changelog">
        <p class="changelog-summary-label">At a glance</p>
        <p class="changelog-summary-value">{{ frontmatter.latestBatch }}</p>
        <p class="changelog-summary-copy">
          Product and status markers tell you whether a change is already live, needs an upgrade, or records a past
          incident.
        </p>
        <div class="changelog-summary-rule" />
        <a href="/llms.txt">Machine-readable documentation →</a>
      </aside>
    </header>

    <div class="changelog-feed-shell">
      <section class="changelog-toolbar" aria-label="Filter and subscribe to the changelog">
        <label class="changelog-filter">
          <span>Filter changelog</span>
          <select id="changelog-product" v-model="selectedProduct" name="product">
            <option v-for="product in products" :key="product" :value="product">{{ product }}</option>
          </select>
        </label>

        <div class="changelog-toolbar-links">
          <span class="changelog-result-count" aria-live="polite">
            {{ visibleCount ? `${visibleCount} updates` : 'Updates' }}
          </span>
          <a href="/llms.txt">LLM index</a>
          <a href="/changelog.rss" target="_blank" rel="noreferrer">RSS feed</a>
        </div>
      </section>

      <section class="changelog-status-key" aria-label="Changelog status legend">
        <span><i class="legend-dot legend-live" /><strong>Live</strong> on hosted Connect</span>
        <span><i class="legend-dot legend-upgrade" /><strong>Upgrade</strong> to receive it</span>
        <span><i class="legend-dot legend-history" /><strong>Did not work</strong> records past impact</span>
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
  --changelog-border: color-mix(in srgb, var(--vp-c-divider) 78%, transparent);
  --changelog-rail: color-mix(in srgb, var(--vp-c-brand-1) 24%, var(--vp-c-divider));
  width: min(1180px, 100%);
  margin: 0 auto;
  padding: 64px 24px 112px;
}

.changelog-hero {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(280px, 0.55fr);
  gap: 28px;
  overflow: hidden;
  padding: 56px;
  border: 1px solid var(--changelog-border);
  border-radius: 30px;
  background:
    radial-gradient(circle at 92% 8%, color-mix(in srgb, var(--my-purple) 30%, transparent), transparent 34%),
    radial-gradient(circle at 3% 100%, color-mix(in srgb, var(--vp-c-brand-2) 25%, transparent), transparent 42%),
    color-mix(in srgb, var(--vp-c-bg-soft) 88%, transparent);
  box-shadow: 0 26px 70px color-mix(in srgb, var(--vp-c-brand-1) 9%, transparent);
}

.changelog-hero::after {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image: linear-gradient(color-mix(in srgb, var(--vp-c-divider) 23%, transparent) 1px, transparent 1px),
    linear-gradient(90deg, color-mix(in srgb, var(--vp-c-divider) 23%, transparent) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: linear-gradient(to right, transparent 20%, black 100%);
  content: '';
}

.changelog-hero-copy,
.changelog-summary {
  position: relative;
  z-index: 1;
}

.changelog-kicker,
.changelog-summary-label {
  margin: 0 0 14px;
  color: var(--vp-c-brand-1);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.15em;
  text-transform: uppercase;
}

.changelog-hero h1 {
  max-width: 720px;
  margin: 0;
  background: linear-gradient(110deg, var(--vp-c-text-1) 20%, var(--vp-c-brand-1) 72%, var(--my-purple));
  background-clip: text;
  color: transparent;
  font-size: clamp(48px, 8vw, 86px);
  font-weight: 800;
  letter-spacing: -0.055em;
  line-height: 0.98;
}

.changelog-intro {
  max-width: 680px;
  margin: 24px 0 0;
  color: var(--vp-c-text-2);
  font-size: clamp(17px, 2vw, 21px);
  line-height: 1.65;
}

.changelog-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 32px;
}

.changelog-actions a {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 44px;
  padding: 0 18px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 700;
  text-decoration: none;
  transition: transform 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.changelog-actions a:hover {
  transform: translateY(-2px);
}

.changelog-primary-action {
  color: var(--vp-button-brand-text) !important;
  background: var(--vp-button-brand-bg);
}

.changelog-primary-action:hover {
  background: var(--vp-button-brand-hover-bg);
}

.changelog-secondary-action {
  border: 1px solid var(--changelog-border);
  color: var(--vp-c-text-1) !important;
  background: color-mix(in srgb, var(--vp-c-bg) 68%, transparent);
}

.changelog-summary {
  align-self: stretch;
  padding: 26px;
  border: 1px solid color-mix(in srgb, var(--vp-c-brand-1) 22%, var(--vp-c-divider));
  border-radius: 22px;
  background: color-mix(in srgb, var(--vp-c-bg) 82%, transparent);
  backdrop-filter: blur(12px);
}

.changelog-summary-value {
  margin: 0;
  color: var(--vp-c-text-1);
  font-size: 26px;
  font-weight: 750;
  letter-spacing: -0.025em;
}

.changelog-summary-copy {
  margin: 14px 0 0;
  color: var(--vp-c-text-2);
  font-size: 14px;
  line-height: 1.7;
}

.changelog-summary-rule {
  height: 1px;
  margin: 24px 0 18px;
  background: var(--changelog-border);
}

.changelog-summary a {
  color: var(--vp-c-brand-1);
  font-size: 13px;
  font-weight: 700;
  text-decoration: none;
}

.changelog-feed-shell {
  position: relative;
  width: min(980px, 100%);
  margin: 42px auto 0;
}

.changelog-feed-shell::after {
  position: absolute;
  z-index: -1;
  top: 0;
  right: -72px;
  bottom: 0;
  width: 92px;
  border-right: 1px dashed var(--changelog-border);
  border-left: 1px dashed var(--changelog-border);
  background-image: radial-gradient(color-mix(in srgb, var(--vp-c-divider) 65%, transparent) 0.8px, transparent 0.8px);
  background-size: 12px 12px;
  content: '';
}

.changelog-toolbar {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 24px;
  padding: 26px 0;
  border-top: 1px dashed var(--changelog-border);
  border-bottom: 1px solid var(--changelog-border);
}

.changelog-filter {
  display: grid;
  gap: 8px;
  color: var(--vp-c-text-2);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.changelog-filter select {
  width: 270px;
  min-height: 46px;
  padding: 0 42px 0 15px;
  border: 1px solid var(--changelog-border);
  border-radius: 10px;
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg-soft);
  font: inherit;
  font-size: 14px;
  font-weight: 650;
  letter-spacing: normal;
  text-transform: none;
}

.changelog-filter select:focus-visible {
  outline: 2px solid var(--vp-c-brand-1);
  outline-offset: 2px;
}

.changelog-toolbar-links {
  display: flex;
  align-items: center;
  gap: 18px;
  font-size: 13px;
}

.changelog-toolbar-links a {
  color: var(--vp-c-text-1);
  font-weight: 650;
  text-decoration: none;
}

.changelog-toolbar-links a:hover {
  color: var(--vp-c-brand-1);
}

.changelog-result-count {
  color: var(--vp-c-text-3);
  font-variant-numeric: tabular-nums;
}

.changelog-status-key {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 28px;
  padding: 18px 0;
  border-bottom: 1px dashed var(--changelog-border);
  color: var(--vp-c-text-2);
  font-size: 12px;
}

.changelog-status-key span {
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.changelog-status-key strong {
  color: var(--vp-c-text-1);
}

.legend-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 999px;
  box-shadow: 0 0 0 4px color-mix(in srgb, currentColor 13%, transparent);
}

.legend-live { color: var(--vp-c-brand-1); background: currentColor; }
.legend-upgrade { color: var(--vp-c-yellow-1); background: currentColor; }
.legend-history { color: var(--vp-c-red-1); background: currentColor; }

.changelog-content {
  width: min(700px, calc(100% - 220px));
  margin: 36px auto 0;
}

.changelog-content :deep(> h1) {
  display: none;
}

.changelog-content :deep(.changelog-policy) {
  margin: 0 0 56px;
  padding: 0 0 20px;
  border-bottom: 1px solid var(--changelog-border);
  color: var(--vp-c-text-2);
  font-size: 14px;
  line-height: 1.7;
}

.changelog-content :deep(.changelog-policy summary) {
  cursor: pointer;
  color: var(--vp-c-text-1);
  font-weight: 700;
}

.changelog-content :deep(.changelog-policy p:last-child) {
  margin-bottom: 0;
}

.changelog-content :deep(.changelog-date-group) {
  --date-column: 150px;
  --date-gap: 34px;
  --group-top-space: 56px;
  --group-bottom-space: 56px;
  --rail-x: calc(var(--date-column) + var(--date-gap));
  --node-y: calc(var(--group-top-space) + 6px);
  position: relative;
  display: grid;
  grid-template-columns: var(--date-column) var(--date-gap) minmax(0, 1fr);
  align-items: start;
  width: calc(100% + var(--date-column) + var(--date-gap));
  margin: 0 0 0 calc(-1 * (var(--date-column) + var(--date-gap)));
  padding: var(--group-top-space) 0 var(--group-bottom-space);
}

.changelog-content :deep(.changelog-date-group::after) {
  position: absolute;
  right: 0;
  bottom: 0;
  left: calc(var(--rail-x) + 44px);
  height: 1px;
  background: var(--changelog-border);
  content: '';
}

.changelog-content :deep(.changelog-date-group.is-visible-start) {
  --group-top-space: 0px;
}

.changelog-content :deep(.changelog-date-rail) {
  position: absolute;
  z-index: 1;
  top: 0;
  bottom: 0;
  left: var(--rail-x);
  width: 1px;
  background: var(--changelog-rail);
}

.changelog-content :deep(.changelog-date-group.is-visible-start > .changelog-date-rail) {
  top: var(--node-y);
}

.changelog-content :deep(.changelog-date-group.is-visible-end > .changelog-date-rail) {
  bottom: calc(100% - var(--node-y));
}

.changelog-content :deep(.changelog-date-node) {
  position: absolute;
  z-index: 2;
  top: var(--node-y);
  left: calc(var(--rail-x) - 5px);
  width: 11px;
  height: 11px;
  transform: translateY(-50%);
  border: 2px solid color-mix(in srgb, var(--vp-c-brand-1) 72%, var(--vp-c-divider));
  border-radius: 3px;
  background: var(--vp-c-bg);
}

.changelog-content :deep(h2) {
  position: sticky;
  top: 92px;
  grid-column: 1;
  width: 100%;
  margin: 0;
  border-top: 0;
  color: var(--vp-c-text-3);
  font-family: var(--vp-font-family-mono);
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0;
  line-height: 1;
  text-align: right;
  white-space: nowrap;
}

.changelog-content :deep(h2 .header-anchor) {
  display: none;
}

.changelog-content :deep(ul) {
  grid-column: 3;
  margin: 0;
  padding: 0 0 0 44px;
  list-style: none;
}

.changelog-content :deep(li) {
  margin: 0;
  padding: 0 0 32px;
  color: var(--vp-c-text-2);
  line-height: 1.72;
}

.changelog-content :deep(li + li) {
  padding-top: 32px;
  border-top: 1px solid var(--changelog-border);
}

.changelog-content :deep(li strong) {
  display: block;
  margin: 12px 0 7px;
  color: var(--vp-c-text-1);
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.018em;
  line-height: 1.38;
}

.changelog-content :deep(li .VPBadge) {
  margin-right: 6px;
  vertical-align: 0;
}

.changelog-empty {
  margin: 48px 0;
  color: var(--vp-c-text-2);
  font-size: 15px;
}

@media (max-width: 820px) {
  .changelog-shell {
    padding: 32px 16px 80px;
  }

  .changelog-hero {
    grid-template-columns: 1fr;
    padding: 34px 26px;
    border-radius: 24px;
  }

  .changelog-feed-shell {
    margin-top: 28px;
  }

  .changelog-feed-shell::after {
    display: none;
  }

  .changelog-toolbar {
    align-items: stretch;
    flex-direction: column;
  }

  .changelog-filter select {
    width: 100%;
  }

  .changelog-toolbar-links {
    justify-content: space-between;
  }

  .changelog-status-key {
    align-items: flex-start;
    flex-direction: column;
  }

  .changelog-content {
    width: 100%;
    margin-top: 28px;
  }

  .changelog-content :deep(.changelog-policy) {
    margin-bottom: 38px;
  }

  .changelog-content :deep(.changelog-date-group) {
    --group-top-space: 36px;
    --group-bottom-space: 36px;
    --rail-x: 5px;
    display: block;
    width: 100%;
    margin: 0;
    padding: var(--group-top-space) 0 var(--group-bottom-space);
  }

  .changelog-content :deep(.changelog-date-group::after) {
    left: 29px;
  }

  .changelog-content :deep(h2) {
    position: relative;
    top: auto;
    width: auto;
    height: auto;
    margin: 0 0 18px;
    padding-left: 24px;
    text-align: left;
  }

  .changelog-content :deep(ul) {
    margin-left: 5px;
    padding: 0 0 0 24px;
  }

  .changelog-content :deep(li strong) {
    font-size: 18px;
  }
}
</style>
