export const connectSpigotDownload =
  'https://github.com/minekube/connect-java/releases/download/latest/connect-spigot.jar'

const legacyHomepage =
  '<!DOCTYPE html><meta http-equiv="Refresh" content="0; url=\'https://minekube.com\'" />'

export default {
  async fetch(request) {
    const { pathname } = new URL(request.url)

    if (pathname === '/connect/download' || pathname === '/connect/download/') {
      return Response.redirect(connectSpigotDownload, 302)
    }

    return new Response(legacyHomepage, {
      headers: {
        'content-type': 'text/html;charset=UTF-8',
      },
    })
  },
}
