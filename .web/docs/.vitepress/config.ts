import {defineConfig} from 'vitepress'

import {discordLink, editLink, gitHubLink, projects} from '../shared'
import {additionalTitle, commitRef} from "../shared/cloudflare";

const ogUrl = 'https://connect.minekube.com'
const ogImage = `${ogUrl}/og-image.png`
const ogTitle = 'Minekube Connect'
const ogDescription = 'Next Generation Minecraft Networks'

export default defineConfig({
    title: `Connect${additionalTitle}`,
    description: ogDescription,
    appearance: 'dark',

    head: [
        ['link', {rel: 'icon', type: 'image/png', href: '/favicon.png'}],
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
            message: `Plugins are released under the MIT License (version: ${commitRef}) - Not affiliated with Mojang nor Minecraft`,
            copyright: 'Copyright © 2023 Minekube and Contributors'
        },

        nav: [
            {text: 'Guide', link: '/guide/'},
            {text: 'Quick Start', link: '/guide/quick-start'},
            {text: 'Downloads', link: '/guide/downloads'},
            {text: 'API', link: '/guide/api/'},
            {text: 'Join the Team', link: '/team'},
            ...projects,
        ],

        sidebar: {
            '/guide/': [
                {
                    text: 'Getting Started',
                    items: [
                        {text: 'Introduction', link: '/guide/'},
                        {text: 'Quick Start', link: '/guide/quick-start'},
                        {text: 'Why', link: '/guide/why'},
                    ]
                },
                {
                    text: 'Connectors',
                    items: [
                        {text: 'Overview', link: '/guide/connectors/'},
                        {text: 'Gate Proxy', link: '/guide/connectors/gate'},
                        {text: 'Java Plugin', link: '/guide/connectors/plugin'},
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
                            text: 'Public Localhost',
                            link: '/guide/localhost'
                        },
                        {
                            text: 'Endpoint Domains',
                            link: '/guide/domains'
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
