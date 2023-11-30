import {DefaultTheme} from "vitepress";

export const discordLink = 'https://minekube.com/discord'
export const gitHubLink = 'https://github.com/minekube'

export const projects: DefaultTheme.NavItem = {
    text: 'Gate Proxy',
    link: 'https://gate.minekube.com',
}

export const editLink = (project: string): DefaultTheme.EditLink => {
    return {
        pattern: `${gitHubLink}/docs/edit/main/docs/${project}/:path`,
        text: 'Suggest changes to this page'
    }
}