# go-lambda-sls
ビルド
```
make
```

デプロイ
```
sls deploy --aws-profile {AWS Config Profile}
```

シークレットキーの設定
```
aws secretsmanager create-secret --name SlackInvitationChannelToken  --secret-string {Slackチャンネルのトークン} --profile {AWS Config Profile}
```