import {defineConfig} from 'vitepress'

import {discordLink, editLink, gitHubLink, projects} from '../shared'
import {additionalTitle} from "../shared/cloudflare";
import {genFeed} from "./theme/components/posts/genFeed";

export const ogUrl = 'https://connect.minekube.com'
const ogImage = `${ogUrl}/og-image.png`
const ogTitle = 'Minekube Connect'
const ogDescription = 'The Ingress Tunnel for Minecraft Servers'
const feedUrl = `${ogUrl}/feed.rss`

export default defineConfig({
    title: `Minekube Connect${additionalTitle}`,
    description: ogDescription,
    appearance: 'dark',

    sitemap: {
        hostname: ogUrl,
    },

    head: [
        ['link', {rel: 'icon', type: 'image/png', href: '/favicon.png'}],
        ['link', {rel: 'alternate', type: 'application/rss+xml', title: 'The Minekube Blog RSS Feed', href: feedUrl}],
        ['meta', {property: 'og:type', content: 'website'}],
        ['meta', {property: 'og:title', content: ogTitle}],
        ['meta', {property: 'og:image', content: ogImage}],
        ['meta', {property: 'og:url', content: ogUrl}],
        ['meta', {property: 'og:description', content: ogDescription}],
        ['meta', {name: 'theme-color', content: '#646cff'}],
        // [
        //     'script',
        //     {
        //         src: 'https://cdn.usefathom.com/script.js',
        //         'data-site': 'CBDFBSLI',
        //         'data-spa': 'auto',
        //         defer: ''
        //     }
        // ]
        [
            'script',
            {},
            `!function(t,e){var o,n,p,r;e.__SV||(window.posthog=e,e._i=[],e.init=function(i,s,a){function g(t,e){var o=e.split(".");2==o.length&&(t=t[o[0]],e=o[1]),t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}}(p=t.createElement("script")).type="text/javascript",p.async=!0,p.src=s.api_host+"/static/array.js",(r=t.getElementsByTagName("script")[0]).parentNode.insertBefore(p,r);var u=e;for(void 0!==a?u=e[a]=[]:a="posthog",u.people=u.people||[],u.toString=function(t){var e="posthog";return"posthog"!==a&&(e+="."+a),t||(e+=" (stub)"),e},u.people.toString=function(){return u.toString(1)+".people (stub)"},o="capture identify alias people.set people.set_once set_config register register_once unregister opt_out_capturing has_opted_out_capturing opt_in_capturing reset isFeatureEnabled onFeatureFlags getFeatureFlag getFeatureFlagPayload reloadFeatureFlags group updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures getActiveMatchingSurveys getSurveys onSessionId".split(" "),n=0;n<o.length;n++)g(u,o[n]);e._i.push([i,s,a])},e.__SV=1)}(document,window.posthog||[]);
            posthog.init('phc_h17apkvCV2aUlSQA4BB7WP7AmZHaU14AKnAe9f3ij5S',{api_host:'https://ph.minekube.com'})`
        ]
    ],

    vue: {
        reactivityTransform: true
    },

    buildEnd: genFeed,

    transformHead({pageData}) {
        const frontmatter = pageData.frontmatter
        if (frontmatter.layout !== 'Post') return []

        const title = String(frontmatter.title)
        const description = String(frontmatter.description || ogDescription)
        const url = pageUrl(pageData.relativePath)
        const imageUrl = frontmatter.imageUrl
            ? new URL(String(frontmatter.imageUrl), ogUrl).href
            : ogImage
        const imageAlt = frontmatter.imageAlt
            ? String(frontmatter.imageAlt)
            : title
        const author = normalizeAuthor(frontmatter.author)
        const published = frontmatter.date
            ? new Date(String(frontmatter.date)).toISOString()
            : undefined
        const imageWidth = String(frontmatter.imageWidth || 1200)
        const imageHeight = String(frontmatter.imageHeight || 630)
        const schema = {
            '@context': 'https://schema.org',
            '@type': 'BlogPosting',
            headline: title,
            description,
            image: [imageUrl],
            datePublished: published,
            dateModified: published,
            author: {
                '@type': 'Person',
                name: author.name,
                url: author.href
            },
            publisher: {
                '@type': 'Organization',
                name: 'Minekube',
                logo: {
                    '@type': 'ImageObject',
                    url: `${ogUrl}/minekube-logo.png`
                }
            },
            mainEntityOfPage: {
                '@type': 'WebPage',
                '@id': url
            },
            isPartOf: {
                '@type': 'Blog',
                name: 'The Minekube Blog',
                url: `${ogUrl}/blog/`
            }
        }

        return [
            ['link', {rel: 'canonical', href: url}],
            ['meta', {name: 'description', content: description}],
            ['meta', {property: 'og:type', content: 'article'}],
            ['meta', {property: 'og:title', content: title}],
            ['meta', {property: 'og:description', content: description}],
            ['meta', {property: 'og:image', content: imageUrl}],
            ['meta', {property: 'og:image:alt', content: imageAlt}],
            ['meta', {property: 'og:image:width', content: imageWidth}],
            ['meta', {property: 'og:image:height', content: imageHeight}],
            ['meta', {property: 'og:url', content: url}],
            ['meta', {property: 'article:published_time', content: published}],
            ['meta', {name: 'author', content: author.name}],
            ['meta', {name: 'twitter:card', content: 'summary_large_image'}],
            ['meta', {name: 'twitter:title', content: title}],
            ['meta', {name: 'twitter:description', content: description}],
            ['meta', {name: 'twitter:image', content: imageUrl}],
            ['meta', {name: 'twitter:image:alt', content: imageAlt}],
            ['script', {type: 'application/ld+json'}, JSON.stringify(schema)]
        ]
    },

    themeConfig: {
        logo: '/minekube-logo.png',

        editLink: editLink('connect'),

        socialLinks: [
            {icon: 'discord', link: discordLink},
            {icon: 'github', link: `${gitHubLink}/connect-java`},
            {icon: 'twitter', link: 'https://x.com/minekube'},
            {
                icon: {
                    svg: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-layout-dashboard"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg>'
                },
                link: 'https://app.minekube.com',
                ariaLabel: 'Dashboard'
            },
        ],

        search: {
            provider: 'algolia',
            options: {
                appId: 'HW90LFDYFK',
                apiKey: 'ab80354b5109dc434cd770cb3db6cb2d',
                indexName: 'connect-minekube',
            },
        },

        // carbonAds: {
        //     code: 'CEBIEK3N',
        //     placement: 'vitejsdev'
        // },

        footer: {
            message: `Not affiliated with Mojang nor Minecraft`,
            copyright: `Copyright © ${new Date().getFullYear()} Minekube and Contributors`
        },

        nav: [
            {text: 'Quick Start', link: '/guide/quick-start'},
            {text: 'Downloads', link: '/guide/connectors/plugin#downloading-the-connect-plugin', activeMatch: '^/guide/connectors/plugin'},
            {text: 'Connectors', link: '/guide/connectors/', activeMatch: '^/guide/connectors/'},
            {text: 'Plans', link: '/plans'},
            {text: 'Blog', link: '/blog/', activeMatch: '^/blog/'},
            ...projects,
        ],

        sidebar: {
            '/guide/': [
                {
                    text: 'Getting Started',
                    items: [
                        {text: 'Introduction', link: '/guide/'},
                        {text: 'Quick Start', link: '/guide/quick-start'},
                    ]
                },
                {
                    text: 'Connectors',
                    items: [
                        {text: 'Overview', link: '/guide/connectors/'},
                        {text: 'Gate Connector', link: '/guide/connectors/gate'},
                        {text: 'Plugin Connector', link: '/guide/connectors/plugin'},
                    ]
                },
                {text: 'AuthSession API', link: '/guide/auth-api'},
                {
                    text: 'Guide',
                    items: [
                        {
                            text: 'Joining Servers',
                            link: '/guide/joining'
                        },
                        {
                            text: 'Bedrock Support',
                            link: '/guide/bedrock'
                        },
                        {
                            text: 'Public Localhost',
                            link: '/guide/localhost'
                        },
                        {
                            text: 'Custom Domains',
                            link: '/guide/domains'
                        },
                        {
                            text: 'Forwarding and Topologies',
                            link: '/guide/topologies'
                        },
                        {
                            text: 'Compatibility Matrix',
                            link: '/guide/compatibility'
                        },
                        {
                            text: 'Offline Mode',
                            link: '/guide/offline-mode'
                        },
                        {
                            text: 'About Tunnels',
                            link: '/guide/tunnels'
                        },
                        {
                            text: 'DDoS & Bot Protection',
                            link: '/guide/protections'
                        },
                        {
                            text: 'Advertising Servers',
                            link: '/guide/advertising'
                        },
                        {
                            text: 'Edge Locations',
                            link: '/guide/locations'
                        }
                    ]
                },
                {
                    text: 'Developers API',
                    items: [
                        {
                            text: 'Overview',
                            link: '/guide/api/'
                        },
                        {
                            text: 'API Clients',
                            link: '/guide/api/clients'
                        },
                        {
                            text: 'Authentication',
                            link: '/guide/api/authentication'
                        },
                        {
                            text: 'Super Endpoints',
                            link: '/guide/api/super-endpoints'
                        },
                        {
                            text: 'Code Examples',
                            link: '/guide/api/examples'
                        }
                    ]
                },
                {
                    text: 'Roadmap',
                    items: [
                        {
                            text: 'Changelog',
                            link: '/guide/changelog'
                        },
                        {
                            text: 'Adoption Plan',
                            link: '/guide/adoption-plan'
                        },
                        {
                            text: 'Use cases',
                            link: '/guide/use-cases'
                        },
                        {
                            text: 'Development Roadmap',
                            link: '/guide/roadmap'
                        },
                        {text: 'Why', link: '/guide/why'},
                    ]
                },
                // {
                //     text: 'APIs',
                //     items: [
                //         {
                //             text: 'Plugin API',
                //             link: '/guide/api-plugin'
                //         },
                //     ]
                // }
            ],
        }
    }
})

function pageUrl(relativePath: string) {
    const path = relativePath
        .replace(/(^|\/)index\.md$/, '$1')
        .replace(/\.md$/, '.html')

    return new URL(path.startsWith('/') ? path : `/${path}`, ogUrl).href
}

function normalizeAuthor(author: unknown) {
    if (author && typeof author === 'object' && 'name' in author) {
        const data = author as { name?: unknown, href?: unknown }
        return {
            name: typeof data.name === 'string' ? data.name : 'Minekube',
            href: typeof data.href === 'string' ? data.href : undefined
        }
    }

    return {
        name: typeof author === 'string' ? author : 'Minekube',
        href: undefined
    }
}
