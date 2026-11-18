<template>
  <section class="mk-compare" aria-labelledby="compare-title">
    <div class="head">
      <h2 id="compare-title" class="h2">Compare the plans</h2>
      <p class="sub">Every organization starts on Free, forever. Plus is for networks that want their own branding and fewer ecosystem defaults.</p>
    </div>
    <div class="scroll">
      <table>
        <thead>
          <tr>
            <th scope="col" class="feat"><span class="sr">Feature</span></th>
            <th v-for="tier in tiers" :key="tier.id" scope="col" :class="tier.mostPopular ? 'col-plus' : ''">
              <span class="t-name">{{ tier.name }}</span>
              <span class="t-price">{{ tier.priceMonthly }}<small>/mo</small></span>
              <a :href="tier.href" class="t-cta">{{ tier.ctaText }}</a>
            </th>
          </tr>
        </thead>
        <tbody v-for="section in sections" :key="section.name">
          <tr class="group">
            <th scope="colgroup">{{ section.name }}</th>
            <td v-for="tier in tiers" :key="tier.id" :class="tier.mostPopular ? 'col-plus' : ''" aria-hidden="true"></td>
          </tr>
          <tr v-for="feature in section.features" :key="feature.name">
            <th scope="row">{{ feature.name }}</th>
            <td v-for="tier in tiers" :key="tier.id" :class="tier.mostPopular ? 'col-plus' : ''">
              <span v-if="typeof feature.tiers[tier.name] === 'string'">{{ feature.tiers[tier.name] }}</span>
              <template v-else>
                <span v-if="feature.tiers[tier.name] === true" class="yes"><svg class="ck" viewBox="0 0 20 20" aria-hidden="true"><path fill="currentColor" d="M16.7 4.2a.75.75 0 0 1 .1 1l-8 10.5a.75.75 0 0 1-1.1.1l-4.5-4.5a.75.75 0 1 1 1-1l3.9 3.8 7.5-9.8a.75.75 0 0 1 1.1-.1Z"/></svg><span class="sr">Included in {{ tier.name }}</span></span>
                <span v-else class="no" aria-hidden="false"><span aria-hidden="true">–</span><span class="sr">Not included in {{ tier.name }}</span></span>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup>
import { discount, plans } from './settings';

const tiers = [
  {
    name: 'Free',
    id: 'tier-free-table',
    ctaText: 'Quick start',
    href: '/guide/quick-start',
    priceMonthly: plans.free.price,
    description: 'Quis suspendisse ut fermentum neque vivamus non tellus.',
    mostPopular: false,
  },
  {
    name: 'Plus',
    id: 'tier-plus-table',
    ctaText: plans.plus.ctaText,
    href: plans.plus.href,
    priceMonthly: plans.plus.price,
    description: 'Quis eleifend a tincidunt pellentesque. A tempor in sed.',
    mostPopular: true,
  },
  // {
  //   name: 'Enterprise',
  //   id: 'tier-enterprise',
  //   href: '#',
  //   priceMonthly: '$109',
  //   description: 'Orci volutpat ut sed sed neque, dui eget. Quis tristique non.',
  //   mostPopular: false,
  // },
];
const sections = [
  {
    name: 'Platform Features',
    features: [
      { name: 'DDoS protected global edge', tiers: { Free: true, Plus: true } },
      { name: 'Custom domains', tiers: { Free: true, Plus: true } },
      { name: 'Custom wildcard domains', tiers: { Free: false, Plus: true } },
      {
        name: 'Unlimited endpoints & connectors',
        tiers: { Free: true, Plus: true },
      },
      { name: 'Unlimited players', tiers: { Free: true, Plus: true } },
      { name: 'Unlimited org members', tiers: { Free: true, Plus: true } },
      { name: 'Activity Tracking', tiers: { Free: true, Plus: true } },
      { name: 'Custom fallback motd', tiers: { Free: false, Plus: true } },
      // { name: 'Team members', tiers: { Free: 'Up to 20 users', Plus: 'Up to 50 users' } },
    ],
  },
  {
    name: 'Opt-out ecosystem default options',
    features: [
      {
        name: 'Prevent Fallback to Browser Hub',
        tiers: { Free: false, Plus: true },
      },
      {
        name: 'Prevent Browser Endpoint Listing',
        tiers: { Free: false, Plus: true },
      },
      { name: 'Prevent /browser Command', tiers: { Free: false, Plus: true } },
      { name: 'Prevent Fallback Motd', tiers: { Free: false, Plus: true } },
      { name: 'Prevent Default Tablist', tiers: { Free: false, Plus: true } },
    ],
  },
  {
    name: 'Support',
    features: [
      {
        name: '24/7 online community support',
        tiers: { Free: true, Plus: true },
      },
      {
        name: 'Staff response within 24 hours',
        tiers: { Free: false, Plus: true },
      },
      { name: '1:1 onboarding tour', tiers: { Free: false, Plus: true } },
      {
        name: 'Early access to new features',
        tiers: { Free: false, Plus: true },
      },
    ],
  },
  // {
  //   name: 'Reporting',
  //   features: [
  //     { name: 'Advanced analytics', tiers: { Basic: true, Essential: true, Premium: true } },
  //     { name: 'Basic reports', tiers: { Essential: true, Premium: true } },
  //     { name: 'Professional reports', tiers: { Premium: true } },
  //     { name: 'Custom report builder', tiers: { Premium: true } },
  //   ],
  // },
  // {
  //   name: 'Support',
  //   features: [
  //     { name: '24/7 online support', tiers: { Basic: true, Essential: true, Premium: true } },
  //     { name: 'Quarterly product workshops', tiers: { Essential: true, Premium: true } },
  //     { name: 'Priority phone support', tiers: { Essential: true, Premium: true } },
  //     { name: '1:1 onboarding tour', tiers: { Premium: true } },
  //   ],
  // },
];
</script>

<style scoped>
.mk-compare { padding-top: 1rem; }
.h2 { margin: 0; font-family: var(--mk-display); font-weight: 750; font-stretch: 82%; font-size: clamp(2rem, 1.3rem + 2.2vw, 3.25rem); line-height: 1.02; letter-spacing: -0.02em; color: var(--sl-color-white); }
.head { display: grid; grid-template-columns: auto minmax(16rem, 30rem); justify-content: space-between; align-items: center; gap: 1.25rem clamp(2.5rem, 6vw, 6rem); }
.head h2 { white-space: nowrap; }
@media (max-width: 72rem) { .head h2 { white-space: normal; } }
@media (max-width: 56rem) { .head { grid-template-columns: 1fr; } }

.sub { margin: 0; font-size: 1.0625rem; line-height: 1.6; color: var(--sl-color-gray-2); }
.mk-compare { min-width: 0; max-width: 100%; }
.scroll { position: relative; max-width: 100%; margin-top: 3rem; overflow-x: auto; border-radius: 1.1rem; background: var(--sl-color-gray-6); }
table { width: 100%; min-width: 28rem; border-collapse: collapse; display: table; margin: 0; font-size: 1rem; }
th, td { padding: 0.85rem 1.25rem; text-align: left; border: 0; border-bottom: 1px solid var(--sl-color-hairline-light); background: transparent; color: var(--sl-color-gray-1); }
thead th { vertical-align: bottom; padding-top: 1.5rem; padding-bottom: 1.25rem; }
td, thead th:not(.feat) { text-align: center; width: 22%; }
.col-plus { background: color-mix(in srgb, var(--mk-orange) 9%, transparent); }
.t-name { display: block; font-weight: 700; color: var(--sl-color-white); }
.col-plus .t-name { color: var(--mk-orange); }
.t-price { display: block; margin-top: 0.25rem; font-family: var(--mk-display); font-weight: 800; font-stretch: 80%; font-size: 2rem; color: var(--sl-color-white); }
.t-price small { font-family: var(--sl-font); font-size: 0.9rem; font-weight: 600; color: var(--sl-color-gray-3); margin-left: 0.15rem; }
.t-cta { display: inline-flex; margin-top: 0.75rem; padding: 0.5rem 1rem; border-radius: 999px; font-size: 0.9375rem; font-weight: 700; text-decoration: none; background: var(--sl-color-gray-5); color: var(--sl-color-white); }
.col-plus .t-cta { background: var(--mk-orange); color: #160a01; }
.group > * { padding-top: 1.75rem; border-bottom-color: var(--sl-color-gray-5); }
.group th { font-size: 0.9375rem; font-weight: 700; color: var(--sl-color-white); border-bottom-color: var(--sl-color-gray-5); }
tbody th[scope='row'] { font-weight: 500; }
tbody tr:last-child > * { border-bottom: 0; }
.yes { display: inline-flex; color: var(--mk-orange); }
.yes .ck { width: 1.2rem; height: 1.2rem; }
.no { color: var(--sl-color-gray-3); font-weight: 600; }
.sr { position: absolute; width: 1px; height: 1px; overflow: hidden; clip: rect(0 0 0 0); white-space: nowrap; }
</style>
