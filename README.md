# Hisab

I built Hisab as a side-project/startup over a few weeks outside of work in my spare time.

Please note, this was 'hacked' together to launch fast and iterate and many flows are incomplete. This was built under the premise that this could be thrown away, and would be rebuilt/refactored if/when it was to scale.

## What Is Hisab?
Allows users to connect their bank accounts & cards to view and manage finances in one place.
Uses [Lean](https://www.leantech.me/) open banking integration.

## Why did I stop?
The project started due to my own struggle to figure out where my money was going in the UAE. Neobanks like Monzo weren't available and Personal Finance Managers like Emma & Cleo didn't cover this region.

Initially I thought I could scale this as a B2C app, but through time figured it'd be very difficult to monetise and generate revenue due to customer base & other startups which had failed to do so in the west [see Mint](https://www.cnbc.com/select/mint-budgeting-app-is-going-away-here-are-some-alternatives/#:~:text=Mint%20will%20go%20offline%20March%2023%2C%202024.&text=The%20Mint%20budgeting%20app%20officially,Karma%2C%20which%20it%20also%20owns).

It was cool to start building something from scratch and figure out the business side of things on the way.

## How far did I get?
Backend:
- User authentication, token refresh flows
- Integrate with Lean (which had a pretty complex flow)
- Allow customers to connect their entities (banks)
- Webhooks to consume user entities, accounts, cards, balances, transactions etc
- Storing this data
- Some APIs to fetch this data
- Setup CI/CD pipelines using Github Actions
- Dockerised
- Deployed to AWS using Elastic Beanstalk & RDS (MySQL)
- Setup VPC and security groups to ensure DB has secure access from backend server only
- Jumpbox to ssh into the server & RDS remotely

Client:
- Authentication flow, token refresh
- 4 page app
- Lean connection flow
- Pull user transactions

## Can I run it?
Probably not on your machine right now. There are some local secrets I haven't pushed and setup/configuration that will take some time to write up.

Happy to take you through it via my machine. It is deployed on AWS but I haven't touched it in some time.