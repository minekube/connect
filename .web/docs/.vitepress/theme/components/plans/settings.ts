export const plans = {
    plus: {
        price: '$10',
        href: 'https://app.minekube.com/:org/checkout?plan=plus',
        ctaText: 'Upgrade to Plus'
    }
}

const endTime = "2024-03-31T23:59:59"
export const discount = {
  active: new Date() < new Date(endTime),
  endTime: endTime,
  price: '$5',
  note: 'This is our Plus plan launch discount!'
}
