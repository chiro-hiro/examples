# How to migrate to AWS Cognito

User authentication is a critical part of system, It need to be secured and scalable. If you have a experienced team and unlimited resource that's fine to build your own implementation follow an open standard likes [OAuth 2.0](https://oauth.net/2/) otherwise [AWS Cognito](https://aws.amazon.com/cognito/) could be the best choice.

## Why AWS Cognito could be the best choice?

If you said "yes" to following questions then AWS Cognito would be your best choice too.

- Do you want your authorization process to be secured and extendable in the future?
- Do you want to support all kind of cross platform authentication?
- Do you want your user to login with other third party account?
- Do you need to allow other third party resource server to be authenticated?

I spent a lot of time to read OAuth 2.0 document but implementing a new authentication service by myself still a heavy task, it cost resource to develop and research even those I'm not sure, did I following a good practice?.

In this article, I'm only able to introduce a minimal guide that allowed you to migrate your own service to AWS Cognito by using [Hosted UI](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-integration.html) and [OAuth 2.0 Authorization Code Grant](https://oauth.net/2/grant-types/authorization-code/).

## Create user pool on AWS Cognito

This is a Getting Started page of AWS Cognito. You could choose to manage your own user pools or identity pools, in this article we're only take care about User Pools.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20Getting%20started%20-%20Amazon%20Cognito%20-%20Amazon%20Web%20Services.png)

### Step 1: Choose User Pool name

I would say, I'm prefer to use **Step through settings** here are some details you won't able to change after pool was created.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito.png)

### Step 2: Customize attributes

In this steps your could define which attributes of user's profile to be stored in AWS Cognito. To me, I'm prefer common information that will be shared across [App Clients](https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-client-apps.html).

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(1).png)

### Step 3: Configure password strength 

You able to config rules to validate user's password. There are many options to initial your User pool. You could allow users to sign them up or create a new account with a temporary password. AWS Cognito provided tools allow your to complete your tasks in the easiest way. You're even able to import users identity from CVS file.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(2).png)


### Step 4: Multi-factor Authentication

I think, MFA is very important. I choose OTP instead of SMS Text message. I don't think SMS is secure enough or at least in my country. SMS text message is alo need some extra step to create group and enable billing to pay for SMS messages.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(3).png)

### Step 5: Customize email sender and fancy email verification

You could spend time to customize a well look email otherwise your customers would receive a boring email ever. Just take note that, email content must contain placeholder.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(4).png)

### Step 6: Remember me feature

There're not munch to say, you see this once in your life.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(5).png)

### Step 7: Create app client

Each app client need to be configured, You might need to review everything carefully. I suggest to grant app client enough accessible to work. I will disable SRP since I use Hosted UI without customization. If you aren't sure what were you doing then I'm prefer to use Hosted UI without modified. You could customize it later so we don't have to haste.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(6).png)

### Step 8: Customize workflows

If you enable **Enable lambda trigger based custom authentication** from above, you will able to add lambda functions.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(7).png)

### Step 9: Complete creation of new user pool

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(8).png)


### Step 10: Customize domain for hosted UI

I'm using subdomain of AWS Cognito to use your own domain you might need some extra steps to config.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(10).png)

### Step 11: Customize UI

It's only support limited changes, at least you could customize colors. You could plain to work on your own UI by using [AWS Amplify](https://docs.amplify.aws/lib/auth/getting-started/q/platform/js).

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(11).png)


### Steps 12: Configure resource server

Now, your backend is a resource server.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(12).png)


### Steps 13: Setup service endpoint URLs

I'm prefer to use [Authorization Code Grant](https://www.oauth.com/oauth2-servers/server-side-apps/authorization-code/).

Why do I prefer? please check this: [What's going on with the OAuth 2.0 Implicit flow?](https://www.youtube.com/watch?v=CHzERullHe8)

TLDR;
Authorization Code could be securer than other kind of authorization methods since:
- Code only able to use one, It become invalid after redeem for tokens.
- Refresh token stored on server side instead of browser.
- User only hold access token, we able to renew access token after it become invalid by exchange refresh token for new tokens.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/Screenshot_2020-09-13%20User%20Pools%20-%20Amazon%20Cognito(15).png)


## Implementation on server side

Here are simplify workflows that we need to implement.

![image](https://raw.githubusercontent.com/chiro-hiro/examples/master/aws-cognito/images/AWS%20Cognito%20-%20workflows.png)

Here are simplify code of our service, it will take authorization code and redeem for tokens:

```js
const express = require('express')
var axios = require('axios');
const app = express()
const port = 3000
const qs = require('querystring')
const client_id='22ahgdj1ou8esn81vpth59ig7d';
const client_secret='1a3ld3169phtvugodc0u6lk0sb2tv0h3ohmhval9kelia8kmhv37';

app.use(express.json());

axios.interceptors.request.use(request => {
  console.log('Starting Request', request)
  return request
})

axios.interceptors.response.use(response => {
  console.log('Response:', response)
  return response
})

app.get('/sign-in/', async (req, res) => {
  let result = await axios.post('https://chirohirochirohirochirohirochirohirochirohirochirohirochirohiro.auth.ap-southeast-1.amazoncognito.com/oauth2/token',
  qs.stringify({
    grant_type: 'authorization_code',
    code: req.query.code,
    redirect_uri: 'http://localhost:3000/sign-in/',
  }),
  {
    headers: {
      'Content-Type':'application/x-www-form-urlencoded',
      'Authorization': `Basic ${Buffer.from(`${client_id}:${client_secret}`).toString('base64')}`
    }
  });
  console.log(result);
  res.send('Hello World!')
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})
```

There are two options to implement your own service, you could use Amazon Libraries, e.g: [Amazon Cognito Identity SDK for JavaScript](https://www.npmjs.com/package/amazon-cognito-identity-js). Our write your own, It pretty simple since AWS Cognito is following OAuth 2.0. You could write your own code to interactive with AWS Cognito API endpoint e.g: [AUTHORIZATION Endpoint](https://docs.aws.amazon.com/cognito/latest/developerguide/authorization-endpoint.html)