// Proxies public Bedrock verifier metadata from the Connect watch.
//
// The connect-java plugin hard-pins metadata-origin https://connect.minekube.com
// (BedrockPrincipalConfiguration.validOrigin) and production is served by the
// Cloudflare Pages project (static assets + _redirects). Pages _redirects
// cannot proxy to EXTERNAL hosts (a status-200 rule with an external
// destination silently falls through to the static 404 — verified live
// 2026-08-24, PR #157), so this Pages Function fetches the canonical metadata
// from the watch at the edge instead. The watch (moxy cmd/watch) serves the
// v2 document at /.well-known/minekube-connect/bedrock-principal-v2.json and
// the v1 document at bedrock-identity-keys.json with Cache-Control/ETag.
export async function onRequest(context) {
	return proxyWellKnownMetadata(context.request, context.params);
}

export async function proxyWellKnownMetadata(request, params) {
	const segments = Array.isArray(params.path) ? params.path : [];
	const suffix = segments.length > 0 ? `/${segments.join("/")}` : "";
	const upstream = new URL(
		`https://watch-connect.minekube.net/.well-known/minekube-connect${suffix}`,
	);
	const response = await fetch(upstream, {
		method: request.method,
		redirect: "follow",
	});
	const headers = new Headers(response.headers);
	headers.set("access-control-allow-origin", "*");
	return new Response(response.body, {
		status: response.status,
		statusText: response.statusText,
		headers,
	});
}
