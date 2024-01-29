# Use a custom domain

When you create a Connect Endpoint, it is automatically given a `play.minekube.net` sub-domain, based on the endpoint’s
name. This is great for testing and private servers, but when you want to go to full production you’ll want your
endpoint to appear on your own domain.
That’s where the Connect custom domains comes in.
Let’s set up the custom domain, first step: directing traffic to your
endpoint.

## Set a CNAME record

The simplest option for directing traffic to your endpoint is to create a CNAME record for your custom domain that
points at your `<endpoint>.play.minekube.net` host. For example, if you have a custom domain called `example.com` and an
endpoint called `mcserver`, then you can create a CNAME record for `example.com`‘s DNS that would look like:

```
CNAME @ mcserver.play.minekube.net
```

You’ll need to configure this with your DNS provider.

Now, accessing `example.com` will tell the DNS system to look up `mcserver.play.minekube.net` and return its results.

## Add your domain

You'll need to go to the [Connect Dashboard](https://app.minekube.com) and add your domain to your endpoint.

1. Go to the endpoint you want to add a domain to
2. Click on `Custom Domains`
3. Click on `Add Domain`
4. Enter your domain

### Domain verification

Your domain will be verified automatically when it detects the CNAME record set up previously.
When you are done, you will be able to join your connected endpoint with your domain, within a few seconds.
