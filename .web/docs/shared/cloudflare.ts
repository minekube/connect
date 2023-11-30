// cloudflare envs
export const commitRef = process.env.CF_PAGES_COMMIT_SHA?.slice(0, 8) || 'dev'

export const deployType = (() => {
    if (commitRef === '') {
        return 'local'
    }
    return 'release'
})()

export const additionalTitle = ((): string => {
    if (deployType === 'release') {
        return ''
    }
    return ' (local)'
})()