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

        algolia: {
            appId: 'HW90LFDYFK',
            apiKey: 'ab80354b5109dc434cd770cb3db6cb2d',
            indexName: 'connect-minekube',
        },

        // carbonAds: {
        //     code: 'CEBIEK3N',
        //     placement: 'vitejsdev'
        // },

        footer: {
            message: `Plugins are released under the MIT License. (web version: ${commitRef})`,
            copyright: 'Copyright © 2022-present Minekube & Contributors'
        },

        nav: [
            {text: 'Guide', link: '/guide/'},
            {text: 'Quick Start', link: '/guide/quick-start'},
            {text: 'Downloads', link: '/guide/downloads'},
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
                            text: 'Joining your Server',
                            link: '/guide/joining'
                        },
                        {
                            text: 'Advertising your Server',
                            link: '/guide/advertising'
                        },
                        {
                            text: 'Server Domains',
                            link: '/guide/domains'
                        },
                        {
                            text: 'DDoS & Bot protection',
                            link: '/guide/protections'
                        },
                        {
                            text: 'Authentication',
                            link: '/guide/authentication'
                        },
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
