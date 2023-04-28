import {defineConfig} from 'vitepress'

import {
    discordLink,
    gitHubLink,
    additionalTitle,
    commitRef,
    editLink,
   projects
} from '../shared'

const ogUrl = 'https://connect.minekube.com'
const ogImage = `${ogUrl}/og-image.png`
const ogTitle = 'Minekube Connect'
const ogDescription = 'Next Generation Minecraft Networks'

export default defineConfig({
    title: `Connect${additionalTitle}`,
    description: ogDescription,

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
    ],

    vue: {
        reactivityTransform: true
    },

    themeConfig: {
        logo: '/minekube-logo.png',

        editLink: editLink('connect'),

        socialLinks: [
            {icon: 'discord', link: discordLink},
            {icon: 'github', link: `${gitHubLink}/connect-java`}
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
            projects,
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
                    text: 'Guide',
                    items: [
                        {
                            text: 'Downloads',
                            link: '/guide/downloads'
                        },
                        {
                            text: 'Joining Servers',
                            link: '/guide/joining'
                        },
                        {
                            text: 'Public Localhost',
                            link: '/guide/localhost'
                        },
                        {
                            text: 'Advertising Servers',
                            link: '/guide/advertising'
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
                            text: 'Authentication',
                            link: '/guide/api/authentication'
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
