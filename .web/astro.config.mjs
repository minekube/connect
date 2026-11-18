import {defineConfig} from 'astro/config'
import starlight from '@astrojs/starlight'
import vue from '@astrojs/vue'
import tailwindcss from '@tailwindcss/vite'
import {remarkCompat} from './src/plugins/markdown.mjs'
import {indexPages} from './src/plugins/index-pages.mjs'
import sitemap from '@astrojs/sitemap'
import {readdirSync} from 'node:fs'

// Directories that hold an index page are served with a trailing slash (/guide/, /blog/).
const indexDirs = new Set(
  readdirSync('./src/content/docs', {recursive: true})
    .map(String).map((f) => f.split('\\').join('/'))
    .filter((f) => /(^|\/)index\.mdx?$/.test(f) && f.includes('/'))
    .map((f) => '/' + f.replace(/\/index\.mdx?$/, '')),
)

export const site = 'https://connect.minekube.com'
const discordLink = 'https://minekube.com/discord'

export default defineConfig({
  site,
  trailingSlash: 'ignore',
  // Load pages on hover/focus so moving through the docs feels instant.
  prefetch: {prefetchAll: true, defaultStrategy: 'hover'},
  // guide/quick-start.md -> guide/quick-start.html, guide/index.md -> guide/index.html,
  // the same file layout VitePress produced, so every existing URL keeps resolving.
  build: {format: 'preserve'},
  markdown: {
    remarkPlugins: [remarkCompat],
  },
  vite: {plugins: [tailwindcss()]},
  integrations: [
    // must run before Starlight so Pagefind indexes the final file layout
    indexPages(),
    // Own sitemap (Starlight skips its built-in one when this is present): no include
    // fragments or dated changelog partials, and index pages keep their trailing slash.
    sitemap({
      filter: (page) => !/\/(design\/|guide\/includes|changelog\/20\d\d-)/.test(page) && !page.endsWith('/404'),
      serialize: (item) => {
        const path = new URL(item.url).pathname.replace(/\/$/, '')
        if (indexDirs.has(path)) item.url = site + path + '/'
        return item
      },
    }),
    starlight({
      title: 'Minekube Connect',
      description: 'The Ingress Tunnel for Minecraft Servers',
      logo: {src: './src/assets/minekube-logo.png', replacesTitle: false},
      favicon: '/favicon.png',
      editLink: {baseUrl: 'https://github.com/minekube/connect/edit/main/.web/'},
      lastUpdated: false,
      social: [
        {icon: 'discord', label: 'Discord', href: discordLink},
        {icon: 'github', label: 'GitHub', href: 'https://github.com/minekube/connect-java'},
        {icon: 'x.com', label: 'X', href: 'https://x.com/minekube'},
      ],
      customCss: ['./src/styles/global.css'],
      routeMiddleware: './src/routeData.ts',
      components: {
        Head: './src/components/Head.astro',
        Header: './src/components/Header.astro',
        PageTitle: './src/components/PageTitle.astro',
        MarkdownContent: './src/components/MarkdownContent.astro',
        Footer: './src/components/Footer.astro',
        ThemeSelect: './src/components/ThemeSwitch.astro',
        Search: './src/components/Search.astro',
      },
      expressiveCode: {
        themes: ['github-dark-default', 'github-light-default'],
        styleOverrides: {borderRadius: '0.5rem', frames: {shadowColor: 'transparent'}},
      },
      head: [
        {tag: 'link', attrs: {rel: 'alternate', type: 'application/rss+xml', title: 'The Minekube Blog RSS Feed', href: `${site}/feed.rss`}},
        {tag: 'link', attrs: {rel: 'alternate', type: 'application/rss+xml', title: 'Minekube Changelog RSS Feed', href: `${site}/changelog.rss`}},
        {tag: 'meta', attrs: {name: 'theme-color', content: '#f38020'}},
        {
          tag: 'script',
          content: `!function(t,e){var o,n,p,r;e.__SV||(window.posthog=e,e._i=[],e.init=function(i,s,a){function g(t,e){var o=e.split(".");2==o.length&&(t=t[o[0]],e=o[1]),t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}}(p=t.createElement("script")).type="text/javascript",p.async=!0,p.src=s.api_host+"/static/array.js",(r=t.getElementsByTagName("script")[0]).parentNode.insertBefore(p,r);var u=e;for(void 0!==a?u=e[a]=[]:a="posthog",u.people=u.people||[],u.toString=function(t){var e="posthog";return"posthog"!==a&&(e+="."+a),t||(e+=" (stub)"),e},u.people.toString=function(){return u.toString(1)+".people (stub)"},o="capture identify alias people.set people.set_once set_config register register_once unregister opt_out_capturing has_opted_out_capturing opt_in_capturing reset isFeatureEnabled onFeatureFlags getFeatureFlag getFeatureFlagPayload reloadFeatureFlags group updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures getActiveMatchingSurveys getSurveys onSessionId".split(" "),n=0;n<o.length;n++)g(u,o[n]);e._i.push([i,s,a])},e.__SV=1)}(document,window.posthog||[]);
posthog.init('phc_h17apkvCV2aUlSQA4BB7WP7AmZHaU14AKnAe9f3ij5S',{api_host:'https://ph.minekube.com'})`,
        },
        {
          // images are not draggable and have no "save image" context menu
          tag: 'script',
          content: `document.addEventListener('contextmenu',e=>{if(e.target instanceof HTMLImageElement)e.preventDefault()});document.addEventListener('dragstart',e=>{if(e.target instanceof HTMLImageElement)e.preventDefault()})`,
        },
        {
          // code-group tabs (rendered by remarkCompat)
          tag: 'script',
          content: `document.addEventListener('click',e=>{const b=e.target.closest('.code-group-tabs button');if(!b)return;const g=b.closest('.code-group');g.querySelectorAll('.code-group-tabs button').forEach(x=>x.setAttribute('aria-selected',x===b));g.querySelectorAll(':scope>.code-group-panel').forEach(p=>p.hidden=p.dataset.panel!==b.dataset.tab)})`,
        },
      ],
      sidebar: [
        {
          label: 'Getting Started',
          items: [
            {label: 'Introduction', link: '/guide/'},
            {label: 'Quick Start', link: '/guide/quick-start'},
          ],
        },
        {
          label: 'Connectors',
          items: [
            {label: 'Overview', link: '/guide/connectors/'},
            {label: 'Gate Connector', link: '/guide/connectors/gate'},
            {label: 'Plugin Connector', link: '/guide/connectors/plugin'},
          ],
        },
        {label: 'AuthSession API', link: '/guide/auth-api', badge: {text: 'Soon', variant: 'danger'}},
        {
          label: 'Guide',
          items: [
            {label: 'Joining Servers', link: '/guide/joining'},
            {label: 'Bedrock Support', link: '/guide/bedrock'},
            {label: 'Public Localhost', link: '/guide/localhost'},
            {label: 'Custom Domains', link: '/guide/domains'},
            {label: 'Forwarding and Topologies', link: '/guide/topologies'},
            {label: 'Compatibility Matrix', link: '/guide/compatibility'},
            {label: 'Login and Auth Plugins', link: '/guide/login-plugins'},
            {label: 'Login Plugin Integration', link: '/guide/login-plugin-integration'},
            {label: 'Offline Mode', link: '/guide/offline-mode'},
            {label: 'About Tunnels', link: '/guide/tunnels'},
            {label: 'DDoS & Bot Protection', link: '/guide/protections'},
            {label: 'Advertising Servers', link: '/guide/advertising'},
            {label: 'Edge Locations', link: '/guide/locations'},
          ],
        },
        {
          label: 'Roadmap',
          collapsed: true,
          items: [
            {label: 'Adoption Plan', link: '/guide/adoption-plan'},
            {label: 'Use cases', link: '/guide/use-cases'},
            {label: 'Development Roadmap', link: '/guide/roadmap'},
            {label: 'Why', link: '/guide/why'},
          ],
        },
        {
          label: 'Minekube',
          items: [
            {label: 'Plans', link: '/plans'},
            {label: 'Changelog', link: '/changelog/'},
            {label: 'Blog', link: '/blog/'},
            {label: 'Team', link: '/team'},
            {label: 'Gate Proxy', link: 'https://gate.minekube.com', attrs: {target: '_blank', rel: 'noopener noreferrer'}},
            {label: 'Dashboard', link: 'https://app.minekube.com', attrs: {target: '_blank', rel: 'noopener noreferrer'}},
            {label: 'Credits', link: '/credits'},
          ],
        },
      ],
    }),
    vue(),
  ],
})
