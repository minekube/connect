import {DefaultTheme} from "vitepress";

export const discordLink = 'https://minekube.com/discord'
export const gitHubLink = 'https://github.com/minekube'

export const projects: DefaultTheme.NavItem[] = [
    {
        text: 'Gate Proxy',
        link: 'https://gate.minekube.com',
    },
    {
        text: 'Dashboard',
        link: 'https://app.minekube.com',
    },
]

export const editLink = (project: string): DefaultTheme.EditLink => {
    return {
        pattern: `${gitHubLink}/${project}/edit/main/.web/docs/:path`,
        text: 'Suggest changes to this page'
    }
}
